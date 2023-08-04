package logic

import (
	"context"

	"doushen_by_liujun/service/chat/api/internal/svc"
	"doushen_by_liujun/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageListLogic {
	return &MessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageListLogic) MessageList(req *types.MessageChatReq) (resp *types.MessageChatReqResp, err error) {
	// todo: add your logic here and delete this line

	return
}
