package logic

import (
	"context"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserinfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserinfoLogic {
	return &UpdateUserinfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserinfoLogic) UpdateUserinfo(in *pb.UpdateUserinfoReq) (*pb.UpdateUserinfoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateUserinfoResp{}, nil
}
