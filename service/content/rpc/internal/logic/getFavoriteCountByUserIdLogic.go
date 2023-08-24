package logic

import (
	"context"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByUserIdLogic {
	return &GetFavoriteCountByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteCountByUserIdLogic) GetFavoriteCountByUserId(in *pb.GetFavoriteCountByUserIdReq) (*pb.GetFavoriteCountByUserIdResp, error) {
	count, err := l.svcCtx.FavoriteModel.GetFavoriteCountByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetFavoriteCountByUserIdResp{
		FavoriteCount: *count,
	}, nil
}
