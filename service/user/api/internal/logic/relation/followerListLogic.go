package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"strconv"
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
			StatusMsg:    common.MapErrMsg(common.TOKEN_EXPIRE_ERROR),
			FollowerList: nil,
		}, e
	}
	followers, e := l.svcCtx.UserRpcClient.GetFollowersById(l.ctx, &pb.GetFollowersByIdReq{
		Id: req.UserId,
	})
	if e != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_api_relation_followerListLogic_FollowerList_GetFollowersById_false"); err != nil {
			log.Fatal(err)
		}
		return &types.FollowerListResp{
			StatusCode:   common.DB_ERROR,
			StatusMsg:    common.MapErrMsg(common.DB_ERROR),
			FollowerList: nil,
		}, e
	}
	var users []types.User
	redisClient := l.svcCtx.RedisClient
	for _, item := range followers.Follows {
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
		user := types.User{
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
	if err := l.svcCtx.KqPusherClient.Push("user_api_relation_followerListLogic_FollowerList_success"); err != nil {
		log.Fatal(err)
	}
	return &types.FollowerListResp{
		StatusCode:   common.OK,
		StatusMsg:    common.MapErrMsg(common.OK),
		FollowerList: users,
	}, nil
}
