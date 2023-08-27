package logic

import (
	"context"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPasswordByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPasswordByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPasswordByUsernameLogic {
	return &GetPasswordByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPasswordByUsernameLogic) GetPasswordByUsername(in *pb.GetPasswordByUsernameReq) (*pb.GetPasswordByUsernameResp, error) {
	l.Logger.Info("GetPasswordByUsername方法请求参数：", in)
	data, err := l.svcCtx.UserinfoModel.GetPasswordByUsername(l.ctx, in.Username)
	if err != nil {
		return nil, err
	}
	return &pb.GetPasswordByUsernameResp{
		Password: data.Password,
		Id:       data.Id,
	}, nil
}
