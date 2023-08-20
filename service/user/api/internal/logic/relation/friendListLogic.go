package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"log"
	"strconv"

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
	_, e := util.ParseToken(req.Token)
	if e != nil {
		return &types.FriendListResp{
			StatusCode: common.TOKEN_EXPIRE_ERROR,
			StatusMsg:  common.MapErrMsg(common.TOKEN_EXPIRE_ERROR),
			FriendUser: nil,
		}, e
	}
	friends, err := l.svcCtx.UserRpcClient.GetFriendsById(l.ctx, &pb.GetFriendsByIdReq{
		Id: req.UserId,
	})
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_api_relation_friendListLogic_FriendList_GetFriendsById_false"); err != nil {
			log.Fatal(err)
		}
		return &types.FriendListResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			FriendUser: nil,
		}, err
	}
	var users []types.FriendUser
	redisClient := l.svcCtx.RedisClient
	for _, item := range friends.Follows {
		workCount := 0
		favoriteCount := 0
		totalFavorited := 0
		workCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserWorkPrefix+strconv.Itoa(int(item.Id)))
		if len(workCountRecord) != 0 { //等于0 代表没有记录，直接赋值0
			//有记录
			workCount, _ = strconv.Atoi(workCountRecord)
		}
		favoriteCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikePrefix+strconv.Itoa(int(item.Id)))
		if len(favoriteCountRecord) != 0 { //等于0 代表没有记录，直接赋值0
			//有记录
			favoriteCount, _ = strconv.Atoi(favoriteCountRecord)
		}
		totalFavoritedRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikedPrefix+strconv.Itoa(int(item.Id)))
		if len(totalFavoritedRecord) != 0 { //等于0 代表没有记录，直接赋值0
			//有记录
			totalFavorited, _ = strconv.Atoi(totalFavoritedRecord)
		}
		user := types.FriendUser{
			UserId:          item.Id,
			Name:            item.UserName,
			FollowCount:     item.FollowCount,
			FollowerCount:   item.FollowerCount,
			IsFollow:        item.IsFollow,
			Avatar:          item.Avator,
			BackgroundImage: item.BackgroundImage,
			Signature:       item.Signature,
			WorkCount:       int64(workCount),
			FavoriteCount:   int64(favoriteCount),
			TotalFavorited:  int64(totalFavorited),
		}
		users = append(users, user)
	}
	if err := l.svcCtx.KqPusherClient.Push("user_api_relation_friendListLogic_FriendList_success"); err != nil {
		log.Fatal(err)
	}
	return &types.FriendListResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		FriendUser: users,
	}, nil
}
