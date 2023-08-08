package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowsCountByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowsCountByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowsCountByIdLogic {
	return &GetFollowsCountByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowsCountByIdLogic) GetFollowsCountById(in *pb.GetFollowsCountByIdReq) (*pb.GetFollowsCountByIdResp, error) {
	// todo: add your logic here and delete this line
	count, err := l.svcCtx.FollowsModel.FindFollowsCount(l.ctx, in.Id)
	return &pb.GetFollowsCountByIdResp{
		Count: int64(count),
	}, err
}
