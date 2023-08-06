package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserinfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserinfoLogic {
	return &SearchUserinfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserinfoLogic) SearchUserinfo(in *pb.SearchUserinfoReq) (*pb.SearchUserinfoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SearchUserinfoResp{}, nil
}
