package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserLogic {
	return &CheckUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckUserLogic) CheckUser(in *pb.CheckUserReq) (*pb.CheckUserResp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("CheckUser方法请求参数：", in)
	id, err := l.svcCtx.UserinfoModel.CheckOne(l.ctx, in.Username, in.Password)
	if err != nil {
		return nil, err
	}
	return &pb.CheckUserResp{
		UserId: *id,
	}, nil
}
