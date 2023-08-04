package logic

import (
	"context"

	"doushen_by_liujun/service/media/rpc/internal/svc"
	"doushen_by_liujun/service/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddChatMessageLogic {
	return &AddChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------鑱婂ぉ淇℃伅-----------------------
func (l *AddChatMessageLogic) AddChatMessage(in *pb.AddChatMessageReq) (*pb.AddChatMessageResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddChatMessageResp{}, nil
}
