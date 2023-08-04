package logic

import (
	"context"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFollowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFollowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFollowsLogic {
	return &AddFollowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------鐢ㄦ埛鍩烘湰淇℃伅-----------------------
func (l *AddFollowsLogic) AddFollows(in *pb.AddFollowsReq) (*pb.AddFollowsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddFollowsResp{}, nil
}
