package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteByIdLogic {
	return &GetFavoriteByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteByIdLogic) GetFavoriteById(in *pb.GetFavoriteByIdReq) (*pb.GetFavoriteByIdResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetFavoriteByIdResp{}, nil
}
