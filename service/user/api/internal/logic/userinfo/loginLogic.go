package userinfo

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/gloabalType"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"time"

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
	l.Logger.Info("login方法请求参数：", req)

	if len(req.Username) > 32 || len(req.Username) < 2 || len(req.Password) < 5 || len(req.Password) > 32 {
		l.Logger.Error("login方法参数错误")
		return &types.LoginResp{
			StatusCode: common.RequestParamError,
			StatusMsg:  common.MapErrMsg(common.RequestParamError),
		}, nil
	}

	data, err := l.svcCtx.UserRpcClient.GetPasswordByUsername(l.ctx, &pb.GetPasswordByUsernameReq{
		Username: req.Username,
	})

	if err != nil {
		l.Logger.Error(err)
		return &types.LoginResp{
			StatusCode: common.AuthorizationError,
			StatusMsg:  common.MapErrMsg(common.AuthorizationError),
		}, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		l.Logger.Error(err)
		return &types.LoginResp{
			StatusCode: common.AuthorizationError,
			StatusMsg:  common.MapErrMsg(common.AuthorizationError),
		}, nil
	}

	token, err := util.GenToken(data.Id, req.Username)
	if err != nil {
		l.Logger.Error(err)
		return &types.LoginResp{
			StatusCode: common.RequestParamError,
			StatusMsg:  common.MapErrMsg(common.RequestParamError),
		}, nil
	}
	ip := l.ctx.Value("ip")
	ipString, ok := ip.(string)
	message := gloabalType.LoginSuccessMessage{}
	if ok {
		message.IP = ipString
		message.Logintime = time.Now()
		message.UserId = data.Id
		messageBytes, err := json.Marshal(message)
		if err != nil {
			l.Logger.Error("无法序列化 message 结构体为 JSON：", err)
		}
		if err := l.svcCtx.LoginLogKqPusherClient.Push(string(messageBytes)); err != nil {
			l.Logger.Error("login方法kafka日志处理失败")
		}
	} else {
		l.Logger.Error("nginx出问题啦")
	}

	return &types.LoginResp{
		UserId:     data.Id,
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		Token:      token,
	}, nil
}
