package logic

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("fail to getChatMsgByIds, error = ?", err)
	}
	for _, item := range *message {
		results = append(results, &pb.ChatMessage{
			Id:         item.Id,
			UserId:     item.UserId,
			ToUserId:   item.ToUserId,
			Message:    item.Message,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}
	return &pb.GetChatMessageByIdResp{ChatMessage: results}, nil
}
