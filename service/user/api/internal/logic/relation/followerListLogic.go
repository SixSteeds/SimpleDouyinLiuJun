package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	contentPB "doushen_by_liujun/service/content/rpc/pb"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
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
		}, nil
	}
	followers, e := l.svcCtx.UserRpcClient.GetFollowersById(l.ctx, &pb.GetFollowersByIdReq{
		Id: req.UserId,
	})
	if e != nil {
		return &types.FollowerListResp{
			StatusCode:   common.DB_ERROR,
			StatusMsg:    common.MapErrMsg(common.DB_ERROR),
			FollowerList: nil,
		}, nil
	}
	var users []types.User
	redisClient := l.svcCtx.RedisClient
	for _, item := range followers.Follows {
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
				return nil, err
			}
			workCount = int(ans.WorkCount)
			redisClient.SetCtx(l.ctx, common.CntCacheUserWorkPrefix+strconv.Itoa(int(item.Id)), strconv.Itoa(workCount))
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
				return nil, err
			}
			favoriteCount = int(ans.FavoriteCount)
			redisClient.SetCtx(l.ctx, common.CntCacheUserLikePrefix+strconv.Itoa(int(item.Id)), strconv.Itoa(favoriteCount))
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
				return nil, err
			}
			totalFavorited = int(ans.LikedCnt)
			redisClient.SetCtx(l.ctx, common.CntCacheUserLikedPrefix+strconv.Itoa(int(item.Id)), strconv.Itoa(totalFavorited))
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
	return &types.FollowerListResp{
		StatusCode:   common.OK,
		StatusMsg:    common.MapErrMsg(common.OK),
		FollowerList: users,
	}, nil
}
