package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
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
	//_, e := util.ParseToken(req.Token)
	//if e != nil {
	//	return &types.FollowerListResp{
	//		StatusCode:   common.TOKEN_EXPIRE_ERROR,
	//		StatusMsg:    "无效token",
	//		FollowerList: nil,
	//	}, e
	//}
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
	var users []types.User
	for _, item := range followers.Follows {
		user := types.User{
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
	return &types.FollowerListResp{
		StatusCode:   common.OK,
		StatusMsg:    "查询粉丝列表成功",
		FollowerList: users,
	}, nil
}
