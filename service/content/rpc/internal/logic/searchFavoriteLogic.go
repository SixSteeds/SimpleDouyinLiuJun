package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchFavoriteLogic {
	return &SearchFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchFavoriteLogic) SearchFavorite(in *pb.SearchFavoriteReq) (*pb.SearchFavoriteResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SearchFavoriteResp{}, nil
}
