package relation

import (
	"context"
	"doushen_by_liujun/service/user/rpc/pb"

	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowLogic) Follow(req *types.FollowReq) (resp *types.FollowResp, err error) {
	// todo: add your logic here and delete this line
	// redis测试
	err = l.svcCtx.RedisClient.Set("test", "test")
	if err != nil {
		return nil, err
	}
	// rpcClientc测试
	one, err := l.svcCtx.UserRpcClient.GetFollowsById(l.ctx, &pb.GetFollowsByIdReq{
		Id: 1,
	})
	if err != nil {
		println("rpcClientc测试失败")
		l.Logger.Info("rpcClientc测试失败")
	}
	l.Logger.Info(one)
	return &types.FollowResp{
		StatusCode: 200,
		StatusMsg:  "关注成功",
	}, nil
}
