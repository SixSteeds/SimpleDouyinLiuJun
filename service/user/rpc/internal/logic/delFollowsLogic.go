package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"errors"
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
	l.Logger.Info(in)
	if in.UserId == in.FollowId {
		//不能取消关注自己
		return nil, errors.New("不能取消关注自己")
	}
	e := l.svcCtx.FollowsModel.DeleteByUserIdAndFollowId(l.ctx, in.UserId, in.FollowId)
	if e != nil {
		l.Logger.Error("删除关注失败", e)
		return nil, e
	}
	return &pb.DelFollowsResp{}, nil
}
