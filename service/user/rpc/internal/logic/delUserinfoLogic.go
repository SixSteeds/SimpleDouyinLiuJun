package logic

import (
	"context"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserinfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserinfoLogic {
	return &DelUserinfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelUserinfoLogic) DelUserinfo(in *pb.DelUserinfoReq) (*pb.DelUserinfoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DelUserinfoResp{}, nil
}
