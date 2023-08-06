package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFollowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFollowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFollowsLogic {
	return &UpdateFollowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFollowsLogic) UpdateFollows(in *pb.UpdateFollowsReq) (*pb.UpdateFollowsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateFollowsResp{}, nil
}
