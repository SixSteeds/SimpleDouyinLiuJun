package logic

import (
	"context"
	"log"

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
	id, err := l.svcCtx.UserinfoModel.CheckOne(l.ctx, in.Username, in.Password)

	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_rpc_checkUSerLogic_CheckUser_CheckOne_false"); err != nil {
			log.Fatal(err)
		}
		return nil, err
	}
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_checkUserLogic_CheckUser_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.CheckUserResp{
		UserId: *id,
	}, nil
}
