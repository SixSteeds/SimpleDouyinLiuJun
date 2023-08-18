package logic

import (
	"context"
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
		followKey := "followNum_" + strconv.Itoa(int(id))
		followerKey := "followerNum_" + strconv.Itoa(int(id))
		followRecord, _ := redisClient.GetCtx(l.ctx, followKey)
		followNum := 0
		followerNum := 0
		expiration := 3600 //秒
		if len(followRecord) == 0 {
			//没有记录，去查表
			num, err := l.svcCtx.FollowsModel.FindFollowsCount(l.ctx, id)
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
			if err != nil {
				return nil, err
			}
			_ = redisClient.SetexCtx(l.ctx, followerKey, strconv.Itoa(num), expiration)
			followerNum = num
		} else {
			//有记录
			followerNum, _ = strconv.Atoi(followerRecord)
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
			WorkCount:       0,
			FavoriteCount:   0,
			TotalFavorited:  0,
			Name:            info.Name.String,
		})
	}

	return &pb.GetUsersByIdsResp{
		Users: users,
	}, nil
}
