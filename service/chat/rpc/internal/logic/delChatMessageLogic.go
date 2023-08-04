package logic

import (
	"context"

	"doushen_by_liujun/service/chat/rpc/internal/svc"
	"doushen_by_liujun/service/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelChatMessageLogic {
	return &DelChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelChatMessageLogic) DelChatMessage(in *pb.DelChatMessageReq) (*pb.DelChatMessageResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DelChatMessageResp{}, nil
}
