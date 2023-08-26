package logic

import (
	"context"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowersCountByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowersCountByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowersCountByIdLogic {
	return &GetFollowersCountByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowersCountByIdLogic) GetFollowersCountById(in *pb.GetFollowersCountByIdReq) (*pb.GetFollowersCountByIdResp, error) {
	l.Logger.Info(in)
	count, err := l.svcCtx.FollowsModel.FindFollowersCount(l.ctx, in.Id)
	return &pb.GetFollowersCountByIdResp{
		Count: int64(count),
	}, err
}
