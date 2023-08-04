package logic

import (
	"context"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelFollowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelFollowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelFollowsLogic {
	return &DelFollowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelFollowsLogic) DelFollows(in *pb.DelFollowsReq) (*pb.DelFollowsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DelFollowsResp{}, nil
}
