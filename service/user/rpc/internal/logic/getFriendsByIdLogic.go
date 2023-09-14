package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"math/rand"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsByIdLogic {
	return &GetFriendsByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendsByIdLogic) GetFriendsById(in *pb.GetFriendsByIdReq) (*pb.GetFriendsByIdResp, error) {
	l.Logger.Info(in)
	follows, err := l.svcCtx.FollowsModel.FindFriendsByUserId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	var resp []*pb.Follows
	redisClient := l.svcCtx.RedisClient
	for _, item := range *follows {
		followKey := common.FollowNum + strconv.Itoa(int(item.Id))
		followerKey := common.FollowerNum + strconv.Itoa(int(item.Id))
		followRecord, _ := redisClient.GetCtx(l.ctx, followKey)
		followNum := 0
		followerNum := 0
		//rand.Seed(time.Now().UnixNano())
		expiration := 3000 + rand.Intn(600)
		if len(followRecord) == 0 {
			//没有记录，去查表
			num, err := l.svcCtx.FollowsModel.FindFollowsCount(l.ctx, item.Id)
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
			num, err := l.svcCtx.FollowsModel.FindFollowersCount(l.ctx, item.Id)
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
		resp = append(resp, &pb.Follows{
			Id:              item.Id,
			FollowerCount:   int64(followerNum),
			FollowCount:     int64(followNum),
			UserName:        item.UserName.String,
			Avator:          item.Avator.String,
			BackgroundImage: item.BackgroundImage.String,
			Signature:       item.Signature.String,
			IsFollow:        item.IsFollow,
		})
	}
	return &pb.GetFriendsByIdResp{
		Follows: resp,
	}, nil
}
