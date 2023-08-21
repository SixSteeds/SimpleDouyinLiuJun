package userinfo

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/gloabalType"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"
	"log"
	"time"

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
	if len(req.Username) > 32 || len(req.Username) < 2 || len(req.Password) < 5 || len(req.Password) > 32 {
		return &types.LoginResp{
			StatusCode: common.REUQEST_PARAM_ERROR,
			StatusMsg:  common.MapErrMsg(common.REUQEST_PARAM_ERROR),
		}, err
	}

	data, err := l.svcCtx.UserRpcClient.CheckUser(l.ctx, &pb.CheckUserReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_loginLogic_Login_CheckUser_false"); err != nil {
			log.Fatal(err)
		}
		return &types.LoginResp{
			StatusCode: common.AUTHORIZATION_ERROR,
			StatusMsg:  common.MapErrMsg(common.AUTHORIZATION_ERROR),
		}, err
	}

	token, err := util.GenToken(data.UserId, req.Username)
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_loginLogic_Login_genToken_false"); err != nil {
			log.Fatal(err)
		}
		return &types.LoginResp{
			StatusCode: common.REUQEST_PARAM_ERROR,
			StatusMsg:  common.MapErrMsg(common.REUQEST_PARAM_ERROR),
		}, err
	}
	if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_loginLogic_Login_success"); err != nil {
		log.Fatal(err)
	}
	ip := l.ctx.Value("ip")
	ipString, ok := ip.(string)
	message := gloabalType.LoginSuccessMessage{}
	if ok {
		message.IP = ipString
		message.Logintime = time.Now()
		message.UserId = data.UserId
		messageString := fmt.Sprintf("%v", message)
		if err := l.svcCtx.KqPusherClient.Push(messageString); err != nil {
			l.Logger.Error("login方法kafka日志处理失败")
		}
	} else {
		l.Logger.Error("nginx出问题啦")
	}

	return &types.LoginResp{
		UserId:     data.UserId,
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		Token:      token,
	}, nil
}
