package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserinfoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserinfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserinfoByIdLogic {
	return &GetUserinfoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserinfoByIdLogic) GetUserinfoById(in *pb.GetUserinfoByIdReq) (*pb.GetUserinfoByIdResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserinfoByIdResp{}, nil
}
