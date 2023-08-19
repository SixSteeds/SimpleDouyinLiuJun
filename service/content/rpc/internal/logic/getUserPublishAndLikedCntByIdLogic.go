package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPublishAndLikedCntByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserPublishAndLikedCntByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPublishAndLikedCntByIdLogic {
	return &GetUserPublishAndLikedCntByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserPublishAndLikedCntByIdLogic) GetUserPublishAndLikedCntById(in *pb.GetUserPublishAndLikedCntByIdReq) (*pb.GetUserPublishAndLikedCntByIdResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserPublishAndLikedCntByIdResp{}, nil
}
