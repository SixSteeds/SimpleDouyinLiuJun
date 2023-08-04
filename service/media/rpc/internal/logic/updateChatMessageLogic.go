package logic

import (
	"context"

	"doushen_by_liujun/service/media/rpc/internal/svc"
	"doushen_by_liujun/service/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChatMessageLogic {
	return &UpdateChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateChatMessageLogic) UpdateChatMessage(in *pb.UpdateChatMessageReq) (*pb.UpdateChatMessageResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateChatMessageResp{}, nil
}
