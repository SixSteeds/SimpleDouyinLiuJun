package userinfo

import (
	"context"
	"doushen_by_liujun/service/user/rpc/pb"

	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.UserRpcClient.AddUserinfo(l.ctx, &pb.AddUserinfoReq{

		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return &types.RegisterResp{
			StatusCode: -1,
			StatusMsg:  "注册失败",
		}, err
	}
	return &types.RegisterResp{
		StatusCode: 200,
		StatusMsg:  "注册成功",
	}, nil
}
