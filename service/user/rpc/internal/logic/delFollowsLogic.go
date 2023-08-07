package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelFollowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelFollowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelFollowsLogic {
	return &DelFollowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelFollowsLogic) DelFollows(in *pb.DelFollowsReq) (*pb.DelFollowsResp, error) {
	// todo: add your logic here and delete this line
	//e := l.svcCtx.FollowsModel.DeleteByUserIdAndFollowId(l.ctx, in.UserId, in.FollowId)
	//if e != nil {
	//	l.Logger.Info("删除关注失败", e)
	//	return nil, errors.New("删除关注失败")
	//}
	return &pb.DelFollowsResp{}, nil
}
