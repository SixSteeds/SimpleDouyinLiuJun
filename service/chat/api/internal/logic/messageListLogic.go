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

	var resp *types.MessageChatReqResp
	var lastTime int64
	if req.PreMsgTime > 169268692200 {
		lastTime = req.PreMsgTime / 1000
	} else {
		lastTime = req.PreMsgTime
	}
	// parse token
	fmt.Println(lastTime)
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
		//createTime, _ := time.Parse("2006-01-02 15:04:05", item.CreateTime)

		msg := types.Message{
			Id:         item.Id,
			ToUserId:   item.ToUserId,
			FromUserId: item.FromUserId,
			Content:    item.Content,
			CreateTime: item.CreateTime, //strconv.Itoa(int(createTime.Unix() )),
		}
		messages = append(messages, msg)
	}

	resp = &types.MessageChatReqResp{
		StatusCode:  0,
		StatusMsg:   "get chat messages successfully",
		MessageList: messages,
	}

	fmt.Println(resp, resp.StatusCode)
	return resp, nil
}
