package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/chat/rpc/internal/model"
	"doushen_by_liujun/service/chat/rpc/internal/svc"
	"doushen_by_liujun/service/chat/rpc/pb"
	"fmt"

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
	l.Logger.Info("AddChatMessage方法请求参数：", in)

	// generate id
	snowflake, err := util.NewSnowflake(common.ChatRpcMachineId)
	if err != nil {

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

		return nil, fmt.Errorf("fail to add chat message record, error = %s", err)
	}

	return &pb.AddChatMessageResp{}, nil
}
