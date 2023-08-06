package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchFollowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchFollowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchFollowsLogic {
	return &SearchFollowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchFollowsLogic) SearchFollows(in *pb.SearchFollowsReq) (*pb.SearchFollowsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SearchFollowsResp{}, nil
}
