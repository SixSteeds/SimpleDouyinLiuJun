package logic

import (
	"context"
	"fmt"

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
	var resp []*pb.Follows
	for _, item := range *follows {
		fmt.Println(item)
		resp = append(resp, &pb.Follows{
			UserId:   item.UserId,
			FollowId: item.FollowId,
		})
	}
	return &pb.GetFollowersByIdResp{
		Follows: resp,
	}, nil
}
