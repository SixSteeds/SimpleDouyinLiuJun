package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	contentPB "doushen_by_liujun/service/content/rpc/pb"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"
	"math/rand"
	"strconv"

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
	l.Logger.Info(in)
	//userId查id这个人
	infos, err := l.svcCtx.UserinfoModel.FindByIds(l.ctx, in.Ids, in.UserID)
	if err != nil || len(*infos) == 0 {
		return nil, err
	}
	var users []*pb.Usersinfo
	for _, info := range *infos {
		id := info.Id
		redisClient := l.svcCtx.RedisClient
		followKey := common.FollowNum + strconv.Itoa(int(id))
		followerKey := common.FollowerNum + strconv.Itoa(int(id))
		followRecord, _ := redisClient.GetCtx(l.ctx, followKey)
		//followRecord, _ := redisClient.Get(followKey)
		followNum := 0
		followerNum := 0
		//rand.Seed(time.Now().UnixNano())
		expiration := 3000 + rand.Intn(600)
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
		//_, _ = redisClient.Get(followerKey)
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
		if len(workCountRecord) != 0 { //等于0 代表没有记录，查表并存储到redis
			//有记录
			workCount, _ = strconv.Atoi(workCountRecord)
		} else {
			ans, err := l.svcCtx.ContentRpcClient.GetWorkCountByUserId(l.ctx, &contentPB.GetWorkCountByUserIdReq{
				UserId: id,
			})
			if err != nil {
				return nil, err
			}
			workCount = int(ans.WorkCount)
			err = redisClient.SetCtx(l.ctx, common.CntCacheUserWorkPrefix+strconv.Itoa(int(id)), strconv.Itoa(workCount))
			if err != nil {
				fmt.Printf("redis set err %v\n", err)
			}
		}
		favoriteCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikePrefix+strconv.Itoa(int(id)))
		if len(favoriteCountRecord) != 0 { //等于0 代表没有记录，查表并存储到redis
			//有记录
			favoriteCount, _ = strconv.Atoi(favoriteCountRecord)
		} else {
			ans, err := l.svcCtx.ContentRpcClient.GetFavoriteCountByUserId(l.ctx, &contentPB.GetFavoriteCountByUserIdReq{
				UserId: id,
			})
			if err != nil {
				return nil, err
			}
			favoriteCount = int(ans.FavoriteCount)
			err = redisClient.SetCtx(l.ctx, common.CntCacheUserLikePrefix+strconv.Itoa(int(id)), strconv.Itoa(favoriteCount))
			if err != nil {
				fmt.Printf("redis set err%v\n", err)
			}
		}
		totalFavoritedRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikedPrefix+strconv.Itoa(int(id)))
		if len(totalFavoritedRecord) != 0 { //等于0 代表没有记录，查表并存储到redis
			//有记录
			totalFavorited, _ = strconv.Atoi(totalFavoritedRecord)
		} else {
			ans, err := l.svcCtx.ContentRpcClient.GetUserPublishAndLikedCntById(l.ctx, &contentPB.GetUserPublishAndLikedCntByIdReq{
				UserId: id,
			})
			if err != nil {
				return nil, err
			}
			totalFavorited = int(ans.LikedCnt)
			err = redisClient.SetCtx(l.ctx, common.CntCacheUserLikedPrefix+strconv.Itoa(int(id)), strconv.Itoa(totalFavorited))
			if err != nil {
				fmt.Printf("redis set err %v\n", err)
			}
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
