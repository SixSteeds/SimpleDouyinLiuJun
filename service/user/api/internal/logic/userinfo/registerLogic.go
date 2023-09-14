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
	l.Logger.Info("register方法请求参数：", req)
	if len(req.Username) > 32 || len(req.Username) < 2 || len(req.Password) < 5 || len(req.Password) > 32 {
		return &types.RegisterResp{
			StatusCode: common.RequestParamError,
			StatusMsg:  common.MapErrMsg(common.RequestParamError),
		}, nil
	}

	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &types.RegisterResp{
			StatusCode: common.ServerCommonError,
			StatusMsg:  common.MapErrMsg(common.ServerCommonError),
		}, nil
	}

	data, err := l.svcCtx.UserRpcClient.SaveUser(l.ctx, &pb.SaveUserReq{
		Username: req.Username,
		Password: string(hashedPassword),
	})

	if err != nil || !data.Success {
		return &types.RegisterResp{
			StatusCode: common.UsernameRepetition,
			StatusMsg:  common.MapErrMsg(common.UsernameRepetition),
		}, nil
	}
	token, err := util.GenToken(data.Id, req.Username)
	if err != nil {
		return &types.RegisterResp{
			StatusCode: common.TokenGenerateError,
			StatusMsg:  common.MapErrMsg(common.TokenGenerateError),
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

	return &types.RegisterResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		UserId:     data.Id,
		Token:      token,
	}, nil
}
