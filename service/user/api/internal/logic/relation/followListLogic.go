package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/api/internal/logic/userinfo"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"
	"strconv"

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
	_, e := util.ParseToken(req.Token)
	if e != nil {
		return &types.FollowListResp{
			StatusCode: common.TOKEN_EXPIRE_ERROR,
			StatusMsg:  "无效token",
			FollowList: nil,
		}, e
	}
	fmt.Println(req.UserId, req.Token) //校验token
	follows, e := l.svcCtx.UserRpcClient.GetFollowsById(l.ctx, &pb.GetFollowsByIdReq{
		Id: req.UserId,
	})
	if e != nil {
		return &types.FollowListResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  "查询关注列表失败",
			FollowList: nil,
		}, e
	}
	userInfo := userinfo.NewUserinfoLogic(l.ctx, l.svcCtx)
	var users []types.User
	for _, item := range follows.Follows {
		fmt.Println(item.FollowId)
		id, _ := strconv.Atoi(item.FollowId)
		resp, err := userInfo.Userinfo(&types.UserinfoReq{UserId: int64(id)})
		if err != nil {
			return &types.FollowListResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  "查询关注列表失败",
				FollowList: nil,
			}, err
		}
		resp.User.IsFollow = true
		users = append(users, resp.User)
	}
	fmt.Println(follows, e)
	return &types.FollowListResp{
		StatusCode: common.OK,
		StatusMsg:  "查询关注列表成功",
		FollowList: users,
	}, nil
}
