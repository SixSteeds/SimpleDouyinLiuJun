package favorite

import (
	"context"
	"database/sql"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"errors"

	constants "doushen_by_liujun/internal/common"

	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"
	"doushen_by_liujun/service/content/rpc/pb"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
)

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListReq) (resp *types.FavoriteListResp, err error) {
	/*
		Author：    刘洋
		Function：  在用户界面返回点赞列表：视频1 [作者（用户信息）、视频信息]、视频2[...]、...
		Update：    08.29 修改点击别人视频作者出现当前用户喜欢列表bug。将使用 token 中 userId 查询改为通过方法传入的 userId 查询
	*/
	l.Logger.Info(req)
	redisClient := l.svcCtx.RedisClient

	// 1.根据 token 获取 userid
	_, err0 := util.ParseToken(req.Token)
	if err0 != nil {
		l.Logger.Error(err0)
		return &types.FavoriteListResp{
			StatusCode: common.TokenExpireError,
			StatusMsg:  common.MapErrMsg(common.TokenExpireError),
		}, nil
	}

	// 2.根据 user_id 查询 favorite 表，返回点赞的所有 video_id
	favoriteListResp, err1 := l.svcCtx.ContentRpcClient.SearchFavorite(l.ctx, &pb.SearchFavoriteReq{
		UserId: req.UserId,
	})
	if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
		l.Logger.Error(err1)
		return &types.FavoriteListResp{
			StatusCode: common.DbError,
			StatusMsg:  common.MapErrMsg(common.DbError),
		}, nil
	}
	fmt.Println("查到favoriteList")
	var favoriteList = favoriteListResp.GetFavorite()

	// 3.依次根据 videoIdList 查询 video 表，得到点赞的所有 video
	var videoIdList []int64
	for _, item := range favoriteList {
		videoIdList = append(videoIdList, item.VideoId)
	}
	if len(videoIdList) == 0 {
		// 如果 videoList 空则直接返回执行完成，不再运行后续逻辑
		return &types.FavoriteListResp{
			StatusCode: common.OK,
			StatusMsg:  common.MapErrMsg(common.OK),
		}, nil
	}
	videoListResp, err2 := l.svcCtx.ContentRpcClient.GetVideoListByIdList(l.ctx, &pb.GetVideoListByIdListReq{
		VideoIdList: videoIdList,
	})
	if err2 != nil && !errors.Is(err2, sql.ErrNoRows) {
		l.Logger.Error(err2)
		return &types.FavoriteListResp{
			StatusCode: common.DbError,
			StatusMsg:  common.MapErrMsg(common.DbError),
		}, nil
	}
	fmt.Println("查到videoList")

	// 4.依次根据 videoList 中的作者id 查询 user 表，得到所有视频作者的 userInfo
	//var authorList []*userPB.Userinfo
	//for _, video := range videoListResp.GetVideoList() {
	//	userResp, err3 := l.svcCtx.UserRpcClient.GetUserById(l.ctx, &userPB.GetUserByIdReq{
	//		UserID: video.UserId,
	//	})
	//	if err3 != nil && err3 != sql.ErrNoRows {
	//		return &types.FavoriteListResp{
	//			StatusCode: common.DB_ERROR,
	//			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
	//		}, nil
	//	}
	//	authorList = append(authorList, userResp.GetUserinfo())
	//}
	//fmt.Println("查到authorList")

	// 4.同时遍历 videoListResp 和 userListResp ，对每个 user 调用服务组合获取计数信息（点赞数、作品数、获赞数）并作为当前遍历到的 video 的成员，
	//   对每个 video 调用服务组合获取计数信息（点赞数、评论数），组装出完整的 video 对象，最终得到 videoList
	var videoList []types.Video
	for index, item := range videoListResp.GetVideoList() {
		fmt.Println(index)
		//// 4.1 将 user 的计数信息组装到 user 对象
		//userLikeCntKey := constants.CntCacheUserLikePrefix + strconv.FormatInt(item.Id, 10)
		//// 4.1.1 从 redis 中拿到用户点赞数
		//userlikeCntRecord, err4 := redisClient.GetCtx(l.ctx, userLikeCntKey)
		//if err4 != nil && err4 != redis.Nil {
		//	return &types.FavoriteListResp{
		//		StatusCode: common.REDIS_ERROR,
		//		StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
		//	}, nil
		//}
		//var userLikeCnt int64 = 0
		//if len(userlikeCntRecord) != 0 {
		//	userLikeCnt, _ = strconv.ParseInt(userlikeCntRecord, 10, 64)
		//}
		//// 4.1.2 根据 userId 调用一个服务组合查询 mysql 多张表，得到用户的发布作品数、获赞数
		//cntResp, err5 := l.svcCtx.ContentRpcClient.GetUserPublishAndLikedCntById(l.ctx, &pb.GetUserPublishAndLikedCntByIdReq{
		//	UserId: authorList[index].Id,
		//})
		//if err5 != nil && err5 != sql.ErrNoRows {
		//	return &types.FavoriteListResp{
		//		StatusCode: common.DB_ERROR,
		//		StatusMsg:  common.MapErrMsg(common.DB_ERROR),
		//	}, nil
		//}
		//fmt.Println("查到publishCnt、likedCnt")
		//// 4.1.3 组装出完整的 user 对象
		//user := types.User{
		//	Id:   item.Id,
		//	Name: authorList[index].Username,
		//	//FollowCount:     userResp.Userinfo.FollowCount,//关注数
		//	//FollowerCount:   userResp.Userinfo.FollowerCount,//被关注数
		//	//IsFollow:        userResp.Userinfo.IsFollow,//当前用户是否关注当前视频作者
		//	Avatar:          authorList[index].Avatar,
		//	BackgroundImage: authorList[index].BackgroundImage,
		//	Signature:       authorList[index].Signature,
		//	WorkCount:       cntResp.PublishCnt, //作品数量
		//	FavoriteCount:   userLikeCnt,        //点赞数量
		//	TotalFavorited:  cntResp.LikedCnt,   //获赞数量
		//}
		// 4.2 将 user 和视频计数信息添加到 video 对象，并加入 videoList

		videoLikedCntKey := constants.CntCacheVideoLikedPrefix + strconv.FormatInt(item.Id, 10)
		videoCommentedCntKey := constants.CntCacheVideoCommentedPrefix + strconv.FormatInt(item.Id, 10)
		// 4.2.1 从 redis 获取当前 video 的点赞数
		videoLikedCntRecord, err3 := redisClient.GetCtx(l.ctx, videoLikedCntKey)
		if err3 != nil && err3 != redis.Nil {
			l.Logger.Error(err3)
			return &types.FavoriteListResp{
				StatusCode: common.RedisError,
				StatusMsg:  common.MapErrMsg(common.RedisError),
			}, nil
		}
		var videoLikedCnt int64 = 0
		if len(videoLikedCntRecord) != 0 {
			videoLikedCnt, _ = strconv.ParseInt(videoLikedCntRecord, 10, 64)
		}
		// 4.2.2 从 redis 获取当前 video 的评论数
		videoCommentedCntRecord, err6 := redisClient.GetCtx(l.ctx, videoCommentedCntKey)
		if err6 != nil && err6 != redis.Nil {
			l.Logger.Error(err6)
			return &types.FavoriteListResp{
				StatusCode: common.RedisError,
				StatusMsg:  common.MapErrMsg(common.RedisError),
			}, nil
		}
		var videoCommentedCnt int64 = 0
		if len(videoLikedCntRecord) != 0 {
			videoCommentedCnt, _ = strconv.ParseInt(videoCommentedCntRecord, 10, 64)
		}
		// 4.2.3 组装出完整的 video 对象，加入 videoList
		videoList = append(videoList, types.Video{
			Id: item.Id,
			//Author:        user,
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: videoLikedCnt,     //视频点赞数
			CommentCount:  videoCommentedCnt, //视频评论数
			IsFavorite:    true,
			Title:         item.Title,
		})
	}
	fmt.Print(videoList)
	fmt.Println("【api-查询用户点赞列表成功】")
	return &types.FavoriteListResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		VideoList:  videoList,
	}, nil
}
