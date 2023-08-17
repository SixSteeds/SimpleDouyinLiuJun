package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
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
	//_, e := util.ParseToken(req.Token)
	//if e != nil {
	//	return &types.FollowListResp{
	//		StatusCode: common.TOKEN_EXPIRE_ERROR,
	//		StatusMsg:  "无效token",
	//		FollowList: nil,
	//	}, e
	//}
	follows, e := l.svcCtx.UserRpcClient.GetFollowsById(l.ctx, &pb.GetFollowsByIdReq{
		Id: req.UserId,
	})
	if e != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_api_relation_followListLogic_FollowList_GetFollowsById_false"); err != nil {
			log.Fatal(err)
		}
		return &types.FollowListResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  "查询关注列表失败",
			FollowList: nil,
		}, e
	}
	var users []types.User
	for _, item := range follows.Follows {
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
	if err := l.svcCtx.KqPusherClient.Push("user_api_relation_followListLogic_FollowList_success"); err != nil {
		log.Fatal(err)
	}
	return &types.FollowListResp{
		StatusCode: common.OK,
		StatusMsg:  "查询关注列表成功",
		FollowList: users,
	}, nil
}
