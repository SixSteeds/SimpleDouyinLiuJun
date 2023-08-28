package logic

import (
	"context"

	"doushen_by_liujun/service/chat/rpc/internal/svc"
	"doushen_by_liujun/service/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchChatMessageLogic {
	return &SearchChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchChatMessageLogic) SearchChatMessage(in *pb.SearchChatMessageReq) (*pb.SearchChatMessageResp, error) {
	l.Logger.Info(in)

	message, _ := l.svcCtx.ChatMessageModel.FindOne(l.ctx, in.Id)

	var results []*pb.ChatMessage
	results = append(results, &pb.ChatMessage{
		Id:         message.Id,
		UserId:     message.UserId,
		ToUserId:   message.ToUserId,
		Message:    message.Message,
		CreateTime: message.CreateTime.String(),
		UpdateTime: message.UpdateTime.String(),
	})

	return &pb.SearchChatMessageResp{ChatMessage: results}, nil
}
