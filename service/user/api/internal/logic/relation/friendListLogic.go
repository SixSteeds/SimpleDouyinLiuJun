package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	contentPB "doushen_by_liujun/service/content/rpc/pb"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"
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
	l.Logger.Info(req)
	_, e := util.ParseToken(req.Token)
	if e != nil {
		l.Logger.Error(e)
		return &types.FriendListResp{
			StatusCode: common.TokenExpireError,
			StatusMsg:  common.MapErrMsg(common.TokenExpireError),
			FriendUser: nil,
		}, nil
	}
	friends, err := l.svcCtx.UserRpcClient.GetFriendsById(l.ctx, &pb.GetFriendsByIdReq{
		Id: req.UserId,
	})
	if err != nil {
		l.Logger.Error(err)
		return &types.FriendListResp{
			StatusCode: common.DbError,
			StatusMsg:  common.MapErrMsg(common.DbError),
			FriendUser: nil,
		}, nil
	}
	var users []types.FriendUser
	redisClient := l.svcCtx.RedisClient
	for _, item := range friends.Follows {
		workCount := 0
		favoriteCount := 0
		totalFavorited := 0
		workCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserWorkPrefix+strconv.Itoa(int(item.Id)))
		if len(workCountRecord) != 0 { //等于0 代表没有记录，查表并存储到redis
			//有记录
			workCount, _ = strconv.Atoi(workCountRecord)
		} else {
			ans, err := l.svcCtx.ContentRpcClient.GetWorkCountByUserId(l.ctx, &contentPB.GetWorkCountByUserIdReq{
				UserId: item.Id,
			})
			if err != nil {
				l.Logger.Error(err)
				return &types.FriendListResp{
					StatusCode: common.DbError,
					StatusMsg:  common.MapErrMsg(common.DbError),
					FriendUser: nil,
				}, nil
			}
			workCount = int(ans.WorkCount)
			err = redisClient.SetCtx(l.ctx, common.CntCacheUserWorkPrefix+strconv.Itoa(int(item.Id)), strconv.Itoa(workCount))
			if err != nil {
				fmt.Printf("redis set err %v\n", err)
			}
		}
		favoriteCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikePrefix+strconv.Itoa(int(item.Id)))
		if len(favoriteCountRecord) != 0 { //等于0 代表没有记录，查表并存储到redis
			//有记录
			favoriteCount, _ = strconv.Atoi(favoriteCountRecord)
		} else {
			ans, err := l.svcCtx.ContentRpcClient.GetFavoriteCountByUserId(l.ctx, &contentPB.GetFavoriteCountByUserIdReq{
				UserId: item.Id,
			})
			if err != nil {
				l.Logger.Error(err)
				return &types.FriendListResp{
					StatusCode: common.DbError,
					StatusMsg:  common.MapErrMsg(common.DbError),
					FriendUser: nil,
				}, nil
			}
			favoriteCount = int(ans.FavoriteCount)
			err = redisClient.SetCtx(l.ctx, common.CntCacheUserLikePrefix+strconv.Itoa(int(item.Id)), strconv.Itoa(favoriteCount))
			if err != nil {
				fmt.Printf("redis set err%v\n", err)
			}
		}
		totalFavoritedRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikedPrefix+strconv.Itoa(int(item.Id)))
		if len(totalFavoritedRecord) != 0 { //等于0 代表没有记录，查表并存储到redis
			//有记录
			totalFavorited, _ = strconv.Atoi(totalFavoritedRecord)
		} else {
			ans, err := l.svcCtx.ContentRpcClient.GetUserPublishAndLikedCntById(l.ctx, &contentPB.GetUserPublishAndLikedCntByIdReq{
				UserId: item.Id,
			})
			if err != nil {
				l.Logger.Error(err)
				return &types.FriendListResp{
					StatusCode: common.DbError,
					StatusMsg:  common.MapErrMsg(common.DbError),
					FriendUser: nil,
				}, nil
			}
			totalFavorited = int(ans.LikedCnt)
			err = redisClient.SetCtx(l.ctx, common.CntCacheUserLikedPrefix+strconv.Itoa(int(item.Id)), strconv.Itoa(totalFavorited))
			if err != nil {
				fmt.Printf("redis set err %v\n", err)
			}
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
	return &types.FriendListResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		FriendUser: users,
	}, nil
}
