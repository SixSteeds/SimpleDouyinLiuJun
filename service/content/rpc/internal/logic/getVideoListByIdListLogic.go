package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoListByIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoListByIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListByIdListLogic {
	return &GetVideoListByIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoListByIdListLogic) GetVideoListByIdList(in *pb.GetVideoListByIdListReq) (*pb.GetVideoListByIdListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetVideoListByIdListResp{}, nil
}
