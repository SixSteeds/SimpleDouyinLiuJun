package relation

import (
	"context"

	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FridendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFridendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FridendListLogic {
	return &FridendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FridendListLogic) FridendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
