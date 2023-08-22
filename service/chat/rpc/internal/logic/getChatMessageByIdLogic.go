package logic

import (
	"context"
	"fmt"
	"log"

	"doushen_by_liujun/service/chat/rpc/internal/svc"
	"doushen_by_liujun/service/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatMessageByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatMessageByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatMessageByIdLogic {
	return &GetChatMessageByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetChatMessageByIdLogic) GetChatMessageById(in *pb.GetChatMessageByIdReq) (*pb.GetChatMessageByIdResp, error) {
	var results []*pb.ChatMessage
	message, err := l.svcCtx.ChatMessageModel.GetChatMsgByIds(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		if err = l.svcCtx.KqPusherClient.Push("chat_rpc_getChatMessageByIdLogic_GetChatMsgByIds_false"); err != nil {
			log.Fatal(err)
		}
		return nil, fmt.Errorf("fail to getChatMsgByIds, error = ?", err)
	}
	for _, item := range *message {
		results = append(results, &pb.ChatMessage{
			Id:         item.Id,
			UserId:     item.UserId,
			ToUserId:   item.ToUserId,
			Message:    item.Message,
			CreateTime: item.CreateTime.String(),
			UpdateTime: item.UpdateTime.String(),
		})
	}
	if err = l.svcCtx.KqPusherClient.Push("chat_rpc_getChatMessageByIdLogic_ GetChatMessageById_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.GetChatMessageByIdResp{ChatMessage: results}, nil
}
