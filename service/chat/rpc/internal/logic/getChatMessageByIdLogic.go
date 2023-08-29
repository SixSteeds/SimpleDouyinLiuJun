package logic

import (
	"context"
	"fmt"
	"strconv"

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
	l.Logger.Info("GetChatMessageById方法请求参数：", in)

	var results []*pb.Message
	message, err := l.svcCtx.ChatMessageModel.GetChatMsgByIds(l.ctx, in.UserId, in.ToUserId, in.PreMsgTime)
	if err != nil {
		return nil, fmt.Errorf("fail to getChatMsgByIds, error = %s", err)
	}
	for _, item := range *message {
		createTime := strconv.Itoa(int(item.CreateTime.Unix()))
		results = append(results, &pb.Message{
			Id:         item.Id,
			ToUserId:   item.ToUserId,
			FromUserId: item.UserId,
			Content:    item.Message,
			CreateTime: &createTime,
		})
	}
	return &pb.GetChatMessageByIdResp{MessageList: results}, nil
}
