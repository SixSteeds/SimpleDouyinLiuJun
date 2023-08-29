package logic

import (
	"context"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/chat/api/internal/svc"
	"doushen_by_liujun/service/chat/api/internal/types"
	"doushen_by_liujun/service/chat/rpc/pb"
	"fmt"
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

func (l *MessageListLogic) MessageList(req *types.MessageChatReq) (*types.MessageChatReqResp, error) {
	l.Logger.Info(req)

	var resp *types.MessageChatReqResp
	var lastTime int64
	if req.PreMsgTime > 169268692200 {
		// 获取第三位数字
		thirdDigit := (req.PreMsgTime / 100) % 10
		// 进行四舍五入
		if thirdDigit >= 5 {
			lastTime = req.PreMsgTime/1000 + 1
		} else {
			lastTime = req.PreMsgTime / 1000
		}
	} else {
		lastTime = req.PreMsgTime
	}

	// parse token
	res, err := util.ParseToken(req.Token)
	if err != nil {
		resp = &types.MessageChatReqResp{
			StatusCode:  1,
			StatusMsg:   "fail to parse token",
			MessageList: nil,
		}
		return resp, fmt.Errorf("fail to parse token, error = %s", err)
	}

	// get params
	userId := res.UserID
	toUserId := req.ToUserId

	request := pb.GetChatMessageByIdReq{
		UserId:     userId,
		ToUserId:   toUserId,
		PreMsgTime: lastTime,
	}
	// get chat messages
	message, err := l.svcCtx.ChatRpcClient.GetChatMessageById(l.ctx, &request)
	if err != nil {
		resp = &types.MessageChatReqResp{
			StatusCode:  1,
			StatusMsg:   "fail to get chat message",
			MessageList: nil,
		}
		return resp, fmt.Errorf("fail to get chat message, error = %s", err)
	}

	var messages []types.Message
	for _, item := range message.MessageList {
		fmt.Println(item.CreateTime)

		msg := types.Message{
			Id:         item.Id,
			ToUserId:   item.ToUserId,
			FromUserId: item.FromUserId,
			Content:    item.Content,
			CreateTime: item.CreateTime,
		}
		messages = append(messages, msg)
	}

	resp = &types.MessageChatReqResp{
		StatusCode:  0,
		StatusMsg:   "get chat messages successfully",
		MessageList: messages,
	}

	return resp, nil
}
