package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFavoritedCntLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFavoritedCntLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFavoritedCntLogic {
	return &GetUserFavoritedCntLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFavoritedCntLogic) GetUserFavoritedCnt(in *pb.GetUserFavoritedCntByIdReq) (*pb.GetUserFavoritedCntByIdResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserFavoritedCntByIdResp{}, nil
}
