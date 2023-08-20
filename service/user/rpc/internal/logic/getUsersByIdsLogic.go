package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"strconv"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersByIdsLogic {
	return &GetUsersByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsersByIdsLogic) GetUsersByIds(in *pb.GetUsersByIdsReq) (*pb.GetUsersByIdsResp, error) {
	// todo: add your logic here and delete this line
	//userId查id这个人

	infos, err := l.svcCtx.UserinfoModel.FindByIds(l.ctx, in.Ids, in.UserID)
	if err != nil {
		return nil, err
	}
	var users []*pb.Usersinfo
	for _, info := range *infos {
		id := info.Id
		redisClient := l.svcCtx.RedisClient
		followKey := common.FollowNum + strconv.Itoa(int(id))
		followerKey := common.FollowerNum + strconv.Itoa(int(id))
		followRecord, _ := redisClient.GetCtx(l.ctx, followKey)
		followNum := 0
		followerNum := 0
		expiration := 3600 //秒
		if len(followRecord) == 0 {
			//没有记录，去查表
			num, err := l.svcCtx.FollowsModel.FindFollowsCount(l.ctx, id)
			num = num - 1 //剪掉自己
			if err != nil {
				return nil, err
			}
			_ = redisClient.SetexCtx(l.ctx, followKey, strconv.Itoa(num), expiration)
			followNum = num
		} else {
			//有记录
			followNum, _ = strconv.Atoi(followRecord)
		}
		followerRecord, _ := redisClient.GetCtx(l.ctx, followerKey)
		if len(followerRecord) == 0 {
			//没有记录，去查表
			num, err := l.svcCtx.FollowsModel.FindFollowersCount(l.ctx, id)
			num = num - 1 //剪掉自己
			if err != nil {
				return nil, err
			}
			_ = redisClient.SetexCtx(l.ctx, followerKey, strconv.Itoa(num), expiration)
			followerNum = num
		} else {
			//有记录
			followerNum, _ = strconv.Atoi(followerRecord)
		}
		workCount := 0
		favoriteCount := 0
		totalFavorited := 0
		workCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserWorkPrefix+strconv.Itoa(int(id)))
		if len(workCountRecord) != 0 { //等于0 代表没有记录，直接赋值0
			//有记录
			workCount, _ = strconv.Atoi(workCountRecord)
		}
		favoriteCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikePrefix+strconv.Itoa(int(id)))
		if len(favoriteCountRecord) != 0 { //等于0 代表没有记录，直接赋值0
			//有记录
			favoriteCount, _ = strconv.Atoi(favoriteCountRecord)
		}
		totalFavoritedRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikedPrefix+strconv.Itoa(int(id)))
		if len(totalFavoritedRecord) != 0 { //等于0 代表没有记录，直接赋值0
			//有记录
			totalFavorited, _ = strconv.Atoi(totalFavoritedRecord)
		}
		users = append(users, &pb.Usersinfo{
			Id:              info.Id,
			FollowCount:     int64(followNum),
			FollowerCount:   int64(followerNum),
			IsFollow:        info.IsFollow,
			Username:        info.Username.String,
			Avatar:          info.Avatar.String,
			BackgroundImage: info.BackgroundImage.String,
			Signature:       info.Signature.String,
			WorkCount:       int64(workCount),
			FavoriteCount:   int64(favoriteCount),
			TotalFavorited:  int64(totalFavorited),
			Name:            info.Name.String,
		})
	}

	return &pb.GetUsersByIdsResp{
		Users: users,
	}, nil
}
