package logic

import (
	"context"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/chat/rpc/pb"
	"fmt"
	"log"

	"doushen_by_liujun/service/chat/api/internal/svc"
	"doushen_by_liujun/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageActionLogic {
	return &MessageActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageActionLogic) MessageAction(req *types.MessageActionReq) (resp *types.MessageActionReqResp, err error) {
	// get params
	token := req.Token
	toUserID := req.ToUserId
	actionType := req.Action_type
	content := req.Content

	// perform corresponding actions based on actionType
	switch actionType {
	case 1:
		// send message
		if err = l.SendMessage(token, content, toUserID); err != nil {
			if err := l.svcCtx.KqPusherClient.Push("chat_api_messageActionLogic_MessageAction_SendMessage_false"); err != nil {
				log.Fatal(err)
			}
			resp = &types.MessageActionReqResp{
				StatusCode: 1,
				StatusMsg:  "fail to send message",
			}
			return resp, fmt.Errorf("fail to send message, error = %s", err)
		}
	default:
		// unknown operation type
		resp = &types.MessageActionReqResp{
			StatusCode: 1,
			StatusMsg:  "fail to send message",
		}
		return resp, fmt.Errorf("unknown operation type")
	}

	// send successfully
	resp = &types.MessageActionReqResp{
		StatusCode: 0,
		StatusMsg:  "send message successfully",
	}
	if err := l.svcCtx.KqPusherClient.Push("chat_api_messageActionLogic_MessageAction_success"); err != nil {
		log.Fatal(err)
	}

	return resp, nil
}

func (l *MessageActionLogic) SendMessage(token, content string, toUserId int64) error {
	// TODO：get permission
	res, err := util.ParseToken(token)
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("chat_api_messageActionLogic_SendMessage_ParseToken_false"); err != nil {
			log.Fatal(err)
		}
		return fmt.Errorf("fail to parse token, error = %s", err)
	}

	// get userId
	userId := res.UserID

	// TODO：checkUserExists

	// add message
	request := &pb.AddChatMessageReq{
		UserId:   userId,
		ToUserId: toUserId,
		Message:  content,
		IsDelete: 0,
	}
	_, err = l.svcCtx.ChatRpcClient.AddChatMessage(l.ctx, request)
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("chat_api_messageActionLogic_SendMessage_AddChatMessage_false"); err != nil {
			log.Fatal(err)
		}
		return fmt.Errorf("fail to send message, error = %s", err)
	}

	return nil
}
