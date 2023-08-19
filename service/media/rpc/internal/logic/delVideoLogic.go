package logic

import (
	"context"

	"doushen_by_liujun/service/media/rpc/internal/svc"
	"doushen_by_liujun/service/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelVideoLogic {
	return &DelVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelVideoLogic) DelVideo(in *pb.DelVideoReq) (*pb.DelVideoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DelVideoResp{}, nil
}
