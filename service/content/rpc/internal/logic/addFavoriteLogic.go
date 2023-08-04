package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFavoriteLogic {
	return &AddFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------鐐硅禐淇℃伅-----------------------
func (l *AddFavoriteLogic) AddFavorite(in *pb.AddFavoriteReq) (*pb.AddFavoriteResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddFavoriteResp{}, nil
}
