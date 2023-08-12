package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	//_, e := util.ParseToken(req.Token)
	//if e != nil {
	//	return &types.FriendListResp{
	//		StatusCode: common.TOKEN_EXPIRE_ERROR,
	//		StatusMsg:  "无效token",
	//		FriendUser: nil,
	//	}, e
	//}
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
	var users []types.FriendUser
	for _, item := range friends.Follows {
		user := types.FriendUser{
			UserId:          item.Id,
			Name:            item.UserName,
			FollowCount:     item.FollowCount,
			FollowerCount:   item.FollowerCount,
			IsFollow:        item.IsFollow,
			Avatar:          item.Avator,
			BackgroundImage: item.BackgroundImage,
			Signature:       item.Signature,
			TotalFavorited:  0, //后三个数据查别人的数据库
			WorkCount:       0,
			FavoriteCount:   0,
		}
		users = append(users, user)
	}
	return &types.FriendListResp{
		StatusCode: common.OK,
		StatusMsg:  "查询好友列表成功",
		FriendUser: users,
	}, nil
}
