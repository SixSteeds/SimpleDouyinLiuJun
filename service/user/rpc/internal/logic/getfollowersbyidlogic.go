package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"
	"log"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowersByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowersByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowersByIdLogic {
	return &GetFollowersByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowersByIdLogic) GetFollowersById(in *pb.GetFollowersByIdReq) (*pb.GetFollowersByIdResp, error) {
	// todo: add your logic here and delete this line
	follows, err := l.svcCtx.FollowsModel.FindByFollowId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_getFollowersByIdLogic_GetFollowersById_FindByFollowId_false"); err != nil {
		log.Fatal(err)
	}
	var resp []*pb.Follows
	redisClient := l.svcCtx.RedisClient
	followKey := "followNum_" + strconv.Itoa(int(in.Id))
	followerKey := "followerNum_" + strconv.Itoa(int(in.Id))
	followRecord, _ := redisClient.GetCtx(l.ctx, followKey)
	followNum := 0
	followerNum := 0
	expiration := 3600 //秒
	if len(followRecord) == 0 {
		//没有记录，去查表
		num, err := l.svcCtx.FollowsModel.FindFollowsCount(l.ctx, in.Id)
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
		num, err := l.svcCtx.FollowsModel.FindFollowersCount(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
		_ = redisClient.SetexCtx(l.ctx, followerKey, strconv.Itoa(num), expiration)
		followerNum = num
	} else {
		//有记录
		followerNum, _ = strconv.Atoi(followerRecord)
	}
	for _, item := range *follows {
		fmt.Println(item)
		resp = append(resp, &pb.Follows{
			Id:              item.Id,
			FollowerCount:   int64(followerNum),
			FollowCount:     int64(followNum),
			UserName:        item.UserName,
			Avator:          item.Avator,
			BackgroundImage: item.BackgroundImage,
			Signature:       item.Signature,
			IsFollow:        item.IsFollow,
		})
	}
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_getFollowersByIdLogic_GetFollowersById_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.GetFollowersByIdResp{
		Follows: resp,
	}, nil
}
