package logic

import (
	"context"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/chat/rpc/internal/model"
	"fmt"
	"log"
	"math/rand"
	"time"

	"doushen_by_liujun/service/chat/rpc/internal/svc"
	"doushen_by_liujun/service/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddChatMessageLogic {
	return &AddChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddChatMessageLogic) AddChatMessage(in *pb.AddChatMessageReq) (*pb.AddChatMessageResp, error) {
	// generate id
	rand.Seed(time.Now().UnixNano())
	snowflake, err := util.NewSnowflake(int64(rand.Intn(1023)))
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("chat_rpc_addChatMessageLogic_AddChatMessage_NewSnowflake_false"); err != nil {
			log.Fatal(err)
		}
		return nil, fmt.Errorf("fail to generate id, error = %s", err)
	}

	snowId := snowflake.Generate()

	// add chat message record
	request := &model.ChatMessage{
		Id:       snowId,
		UserId:   in.UserId,
		ToUserId: in.ToUserId,
		Message:  in.Message,
		IsDelete: 0,
	}
	_, err = l.svcCtx.ChatMessageModel.Insert(l.ctx, request)
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("chat_rpc_addChatMessageLogic_AddChatMessage_Insert_false"); err != nil {
			log.Fatal(err)
		}
		return nil, fmt.Errorf("fail to add chat message record, error = %s", err)
	}

	if err := l.svcCtx.KqPusherClient.Push("chat_rpc_addChatMessageLogic_AddChatMessage_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.AddChatMessageResp{}, nil
}
