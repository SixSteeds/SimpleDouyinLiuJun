package logic

import (
	"context"
	"fmt"
	"log"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

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
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_getFollowersByIdLogic_GetFollowersById_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.GetFollowersByIdResp{
		Follows: resp,
	}, nil
}
