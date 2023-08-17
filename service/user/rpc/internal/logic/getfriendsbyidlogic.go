package logic

import (
	"context"
	"fmt"
	"log"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

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
	follows, err := l.svcCtx.FollowsModel.FindFriendsByUserId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_getFriendsByIdLogic_GetFriendById_FindFriendsByUserId_false"); err != nil {
		log.Fatal(err)
	}
	var resp []*pb.Follows
	for _, item := range *follows {
		fmt.Println(item)
		resp = append(resp, &pb.Follows{
			Id:              item.Id,
			FollowerCount:   item.FollowerCount,
			FollowCount:     item.FollowCount,
			UserName:        item.UserName,
			Avator:          item.Avator,
			BackgroundImage: item.BackgroundImage,
			Signature:       item.Signature,
			IsFollow:        item.IsFollow,
		})
	}
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_getFriendsByIdLogic_GetFriendById_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.GetFriendsByIdResp{
		Follows: resp,
	}, nil
}
