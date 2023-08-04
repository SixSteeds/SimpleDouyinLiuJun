package logic

import (
	"context"

	"doushen_by_liujun/service/media/rpc/internal/svc"
	"doushen_by_liujun/service/media/rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.GetChatMessageByIdResp{}, nil
}
