package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
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
	logger, e := util.ParseToken(req.Token)
	if e != nil {
		return &types.FollowResp{
			StatusCode: common.TOKEN_EXPIRE_ERROR,
			StatusMsg:  "token失效",
		}, err
	}
	//判断是关注还是取消关注
	if req.ActionType == 1 {
		_, err := l.svcCtx.UserRpcClient.AddFollows(l.ctx, &pb.AddFollowsReq{
			UserId:   strconv.Itoa(int(logger.UserID)),
			FollowId: strconv.FormatInt(req.ToUserId, 10),
		})
		if err != nil {
			return &types.FollowResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  "关注失败",
			}, err
		}
		return &types.FollowResp{
			StatusCode: common.OK,
			StatusMsg:  "关注成功",
		}, nil
	} else {
		_, err := l.svcCtx.UserRpcClient.DelFollows(l.ctx, &pb.DelFollowsReq{
			UserId:   strconv.Itoa(int(logger.UserID)),
			FollowId: strconv.FormatInt(req.ToUserId, 10),
		})
		if err != nil {
			return &types.FollowResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  "删除关注失败",
			}, err
		}
		return &types.FollowResp{
			StatusCode: common.OK,
			StatusMsg:  "取关成功",
		}, err
	}
}
