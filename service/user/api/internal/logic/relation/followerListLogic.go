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

type FollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowerListLogic {
	return &FollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowerListLogic) FollowerList(req *types.FollowerListReq) (resp *types.FollowerListResp, err error) {
	_, e := util.ParseToken(req.Token)
	if e != nil {
		return &types.FollowerListResp{
			StatusCode:   common.TOKEN_EXPIRE_ERROR,
			StatusMsg:    "无效token",
			FollowerList: nil,
		}, e
	}
	followers, e := l.svcCtx.UserRpcClient.GetFollowersById(l.ctx, &pb.GetFollowersByIdReq{
		Id: req.UserId,
	})
	if e != nil {
		return &types.FollowerListResp{
			StatusCode:   common.DB_ERROR,
			StatusMsg:    "查询粉丝列表失败",
			FollowerList: nil,
		}, e
	}
	userInfo := userinfo.NewUserinfoLogic(l.ctx, l.svcCtx)
	var users []types.User
	for _, item := range followers.Follows {
		fmt.Println(item.FollowId)
		id, _ := strconv.Atoi(item.UserId)
		resp, err := userInfo.Userinfo(&types.UserinfoReq{UserId: int64(id)})
		if err != nil {
			return &types.FollowerListResp{
				StatusCode:   common.DB_ERROR,
				StatusMsg:    "查询粉丝列表失败",
				FollowerList: nil,
			}, err
		}
		isFollowed, e := l.svcCtx.UserRpcClient.CheckIsFollow(l.ctx, &pb.CheckIsFollowReq{
			Userid:   item.FollowId,
			Followid: item.UserId,
		})
		if e != nil {
			return &types.FollowerListResp{
				StatusCode:   common.DB_ERROR,
				StatusMsg:    "查询粉丝列表失败",
				FollowerList: nil,
			}, err
		}
		resp.User.IsFollow = isFollowed.IsFollowed
		users = append(users, resp.User)
	}
	return &types.FollowerListResp{
		StatusCode:   common.OK,
		StatusMsg:    "查询粉丝列表成功",
		FollowerList: users,
	}, nil
}
