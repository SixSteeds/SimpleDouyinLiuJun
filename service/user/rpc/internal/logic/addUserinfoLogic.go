package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserinfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserinfoLogic {
	return &AddUserinfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------鐢ㄦ埛鍩烘湰淇℃伅-----------------------
func (l *AddUserinfoLogic) AddUserinfo(in *pb.AddUserinfoReq) (*pb.AddUserinfoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddUserinfoResp{}, nil
}
