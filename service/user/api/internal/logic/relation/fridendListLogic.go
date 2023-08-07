package relation

import (
	"context"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"

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
	fmt.Println(req.UserId, req.Token) //校验token
	follows, e := l.svcCtx.UserRpcClient.GetFriendsById(l.ctx, &pb.GetFriendsByIdReq{
		Id: req.UserId,
	})
	fmt.Println("查好友列表啦！！！！！！")
	fmt.Println(follows, e)
	//拿到了两条数据，还要查别人的表，之后写
	return
}
