package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteLogic {
	return &UpdateFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFavoriteLogic) UpdateFavorite(in *pb.UpdateFavoriteReq) (*pb.UpdateFavoriteResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateFavoriteResp{}, nil
}
