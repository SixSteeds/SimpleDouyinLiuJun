package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
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
	userid, _ := strconv.ParseInt(in.UserId, 10, 64)
	followid, _ := strconv.ParseInt(in.FollowId, 10, 64)
	isFollowed, err := l.svcCtx.FollowsModel.CheckIsFollowed(l.ctx, userid, followid) //检查是不是重复操作
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	if !isFollowed { //已经取消关注了
		return &pb.DelFollowsResp{}, nil
	}
	e := l.svcCtx.FollowsModel.DeleteByUserIdAndFollowId(l.ctx, in.UserId, in.FollowId)
	if e != nil {
		l.Logger.Error("删除关注失败", e)
		return nil, e
	}
	return &pb.DelFollowsResp{}, nil
}
