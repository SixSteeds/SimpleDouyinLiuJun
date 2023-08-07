package relation

import (
	"context"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"

	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowListLogic) FollowList(req *types.FollowListReq) (resp *types.FollowListResp, err error) {
	// todo: add your logic here and delete this line
	fmt.Println(req.UserId, req.Token) //校验token
	follows, e := l.svcCtx.UserRpcClient.GetFollowsByFollowId(l.ctx, &pb.GetFollowsByIdReq{
		Id: req.UserId,
	})
	fmt.Println(follows, e)
	//拿到了两条数据，还要查别人的表，之后写
	return
}
