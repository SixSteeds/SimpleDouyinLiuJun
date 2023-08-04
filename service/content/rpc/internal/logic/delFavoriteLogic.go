package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelFavoriteLogic {
	return &DelFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelFavoriteLogic) DelFavorite(in *pb.DelFavoriteReq) (*pb.DelFavoriteResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DelFavoriteResp{}, nil
}
