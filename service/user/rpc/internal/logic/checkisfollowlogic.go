package logic

import (
	"context"
	"strconv"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckIsFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIsFollowLogic {
	return &CheckIsFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckIsFollowLogic) CheckIsFollow(in *pb.CheckIsFollowReq) (*pb.CheckIsFollowResp, error) {
	// todo: add your logic here and delete this line
	userid, _ := strconv.ParseInt(in.Userid, 10, 64)
	followid, _ := strconv.ParseInt(in.Followid, 10, 64)
	isFollowed, err := l.svcCtx.FollowsModel.CheckIsFollowed(l.ctx, userid, followid)
	return &pb.CheckIsFollowResp{
		IsFollowed: isFollowed,
	}, err
}
