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
	"fmt"
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
	if len(req.Username) > 32 || len(req.Username) < 2 || len(req.Password) < 5 || len(req.Password) > 32 {
		return &types.RegisterResp{
			StatusCode: common.REUQEST_PARAM_ERROR,
			StatusMsg:  common.MapErrMsg(common.REUQEST_PARAM_ERROR),
		}, err
	}

	data, err := l.svcCtx.UserRpcClient.SaveUser(l.ctx, &pb.SaveUserReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil || !data.Success {
		return &types.RegisterResp{
			StatusCode: common.USERNAME_REPETITION,
			StatusMsg:  common.MapErrMsg(common.USERNAME_REPETITION),
		}, err
	}
	token, err := util.GenToken(data.Id, req.Username)
	if err != nil {
		return &types.RegisterResp{
			StatusCode: common.TOKEN_GENERATE_ERROR,
			StatusMsg:  common.MapErrMsg(common.TOKEN_GENERATE_ERROR),
		}, err
	}

	ip := l.ctx.Value("ip")
	ipString, ok := ip.(string)
	message := gloabalType.LoginSuccessMessage{}
	if ok {
		fmt.Println("sdasdasdad")
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
