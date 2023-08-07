package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowsByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowsByIdLogic {
	return &GetFollowsByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowsByIdLogic) GetFollowsById(in *pb.GetFollowsByIdReq) (*pb.GetFollowsByIdResp, error) {
	// todo: add your logic here and delete this line
	follows, err := l.svcCtx.FollowsModel.FindByUserId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	var resp []pb.Follows
	for _, item := range *follows {
		fmt.Println(item)
		resp = append(resp, pb.Follows{
			Id:         item.Id,
			UserId:     item.UserId,
			FollowId:   item.FollowId,
			UpdateTime: item.UpdateTime.Unix(),
		})
	}
	return &pb.GetFollowsByIdResp{
		Follows: &resp,
	}, nil
}
