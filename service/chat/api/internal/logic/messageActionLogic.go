package logic

import (
	"context"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/chat/api/internal/svc"
	"doushen_by_liujun/service/chat/api/internal/types"
	"doushen_by_liujun/service/chat/rpc/pb"
	"doushen_by_liujun/service/user/rpc/user"
	"fmt"

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

func (l *MessageActionLogic) MessageAction(req *types.MessageActionReq) (*types.MessageActionReqResp, error) {
	var resp *types.MessageActionReqResp
	// get params
	token := req.Token
	toUserID := req.ToUserId
	actionType := req.Action_type
	content := req.Content

	// perform corresponding actions based on actionType
	switch actionType {
	case 1:
		// send message
		if err := l.SendMessage(token, content, toUserID); err != nil {
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
		return resp, fmt.Errorf("fail to send message, error = unknown operation type")
	}

	// send successfully
	resp = &types.MessageActionReqResp{
		StatusCode: 0,
		StatusMsg:  "send message successfully",
	}

	return resp, nil
}

func (l *MessageActionLogic) SendMessage(token, content string, toUserId int64) error {
	// get permission
	res, err := util.ParseToken(token)
	if err != nil {
		return fmt.Errorf("fail to parse token, error = %s", err)
	}

	// get userId
	userId := res.UserID

	// checkUserExists
	userReq := user.GetUserinfoByIdReq{
		Id:     toUserId,
		UserID: userId,
	}
	response, userInfoErr := l.svcCtx.UserRpcClient.GetUserinfoById(l.ctx, &userReq)
	if userInfoErr != nil {
		return fmt.Errorf("fail to getUserInfo by id, error = %s", userInfoErr)
	}
	if response == nil {
		return fmt.Errorf("No user with id %s", toUserId)
	}

	// add message
	request := &pb.AddChatMessageReq{
		UserId:   userId,
		ToUserId: toUserId,
		Message:  content,
		IsDelete: 0,
	}
	_, err = l.svcCtx.ChatRpcClient.AddChatMessage(l.ctx, request)
	if err != nil {
		logx.Error(err)
		return fmt.Errorf("fail to send message, error = %s", err)
	}

	return nil
}
