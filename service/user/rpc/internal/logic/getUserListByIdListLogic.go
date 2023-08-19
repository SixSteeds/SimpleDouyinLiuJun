package logic

import (
	"context"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListByIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListByIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListByIdListLogic {
	return &GetUserListByIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListByIdListLogic) GetUserListByIdList(in *pb.GetUserListByIdListReq) (*pb.GetUserListByIdListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserListByIdListResp{}, nil
}
