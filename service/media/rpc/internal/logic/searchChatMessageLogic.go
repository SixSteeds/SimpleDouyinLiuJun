package logic

import (
	"context"

	"doushen_by_liujun/service/media/rpc/internal/svc"
	"doushen_by_liujun/service/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchChatMessageLogic {
	return &SearchChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchChatMessageLogic) SearchChatMessage(in *pb.SearchChatMessageReq) (*pb.SearchChatMessageResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SearchChatMessageResp{}, nil
}
