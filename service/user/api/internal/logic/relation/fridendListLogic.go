package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/api/internal/logic/userinfo"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"
	"strconv"

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
	_, e := util.ParseToken(req.Token)
	if e != nil {
		return &types.FriendListResp{
			StatusCode: common.TOKEN_EXPIRE_ERROR,
			StatusMsg:  "无效token",
			FriendUser: nil,
		}, e
	}
	friends, e := l.svcCtx.UserRpcClient.GetFriendsById(l.ctx, &pb.GetFriendsByIdReq{
		Id: req.UserId,
	})
	if e != nil {
		return &types.FriendListResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  "查询好友列表失败",
			FriendUser: nil,
		}, e
	}
	userInfo := userinfo.NewUserinfoLogic(l.ctx, l.svcCtx)
	var users []types.FriendUser
	for _, item := range friends.Follows {
		fmt.Println(item.FollowId)
		id, _ := strconv.Atoi(item.FollowId)
		resp, err := userInfo.Userinfo(&types.UserinfoReq{UserId: int64(id)})
		if err != nil {
			return &types.FriendListResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  "查询好友列表失败",
				FriendUser: nil,
			}, err
		}
		users = append(users, types.FriendUser{
			UserId:          resp.User.UserId,
			Name:            resp.User.Name,
			FollowCount:     resp.User.FollowCount,
			FollowerCount:   resp.User.FollowerCount,
			IsFollow:        true,
			Avatar:          resp.User.Avatar,
			BackgroundImage: resp.User.BackgroundImage,
			Signature:       resp.User.Signature,
			TotalFavorited:  resp.User.TotalFavorited,
			WorkCount:       resp.User.WorkCount,
			FavoriteCount:   resp.User.FavoriteCount,
			Message:         "Message和MsgType后续调用别人接口查询",
			MsgType:         0,
		})
	}
	return &types.FriendListResp{
		StatusCode: common.OK,
		StatusMsg:  "查询好友列表成功",
		FriendUser: users,
	}, nil
}
