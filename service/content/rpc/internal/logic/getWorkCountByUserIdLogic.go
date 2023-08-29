package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkCountByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkCountByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkCountByUserIdLogic {
	return &GetWorkCountByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkCountByUserIdLogic) GetWorkCountByUserId(in *pb.GetWorkCountByUserIdReq) (*pb.GetWorkCountByUserIdResp, error) {
	// todo: add your logic here and delete this line
	count, err := l.svcCtx.VideoModel.GetWorkCountByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetWorkCountByUserIdResp{
		WorkCount: count.Count,
	}, nil
}
