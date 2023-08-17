package userinfo

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/rpc/pb"
	"log"

	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	data, err := l.svcCtx.UserRpcClient.CheckUser(l.ctx, &pb.CheckUserReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_loginLogic_Login_CheckUser_false"); err != nil {
			log.Fatal(err)
		}
		return nil, err
	}

	token, err := util.GenToken(data.UserId, req.Username)
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_loginLogic_Login_genToken_false"); err != nil {
			log.Fatal(err)
		}
		return nil, err
	}
	if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_loginLogic_Login_success"); err != nil {
		log.Fatal(err)
	}
	return &types.LoginResp{
		UserId:     data.UserId,
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		Token:      token,
	}, nil
}
