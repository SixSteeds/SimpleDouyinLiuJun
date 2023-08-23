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
	if len(req.Username) > 32 || len(req.Username) < 2 || len(req.Password) < 5 || len(req.Password) < 32 {
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
		if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_registerLogic_Register_SaveUser_false"); err != nil {
			log.Fatal(err)
		}
		return &types.RegisterResp{
			StatusCode: common.USERNAME_REPETITION,
			StatusMsg:  common.MapErrMsg(common.USERNAME_REPETITION),
		}, err
	}
	token, err := util.GenToken(data.Id, req.Username)
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_registerLogic_Register_genToken_false"); err != nil {
			log.Fatal(err)
		}
		return &types.RegisterResp{
			StatusCode: common.TOKEN_GENERATE_ERROR,
			StatusMsg:  common.MapErrMsg(common.TOKEN_GENERATE_ERROR),
		}, err
	}
	if err := l.svcCtx.KqPusherClient.Push("user_api_userinfo_registerLogic_Register_success"); err != nil {
		log.Fatal(err)
	}
	return &types.RegisterResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		UserId:     data.Id,
		Token:      token,
	}, nil
}
