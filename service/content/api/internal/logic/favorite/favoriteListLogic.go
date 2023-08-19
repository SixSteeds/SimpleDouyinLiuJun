package favorite

import (
	"context"
	"doushen_by_liujun/internal/common"
	constants "doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"
	"doushen_by_liujun/service/content/rpc/pb"
	userPB "doushen_by_liujun/service/user/rpc/pb"
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

	// 点赞列表：视频1（作者）、视频2（..）
	redisClient := l.svcCtx.RedisClient

	// 1.根据 user_id 查询 favorite 表，返回点赞的所有 video_id
	favoriteListResp, err0 := l.svcCtx.ContentRpcClient.SearchFavorite(l.ctx, &pb.SearchFavoriteReq{
		UserId: req.UserId,
	})
	if err0 != nil {
		return &types.FavoriteListResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
		}, err0
	}
	fmt.Println(favoriteListResp)
	fmt.Println("查到favoriteList")
	var favoriteList = favoriteListResp.GetFavorite()
	var videoIdList []int64
	for _, item := range favoriteList {
		videoIdList = append(videoIdList, item.VideoId)
	}
	// 2.根据 videoIdList 查询 video 表，得到点赞的所有 video
	if len(videoIdList) == 0 {
		// 如果 videoList 空则直接返回执行完成，不再运行后续逻辑
		return &types.FavoriteListResp{
			StatusCode: common.OK,
			StatusMsg:  common.MapErrMsg(common.OK),
		}, nil
	}
	videoListResp, err1 := l.svcCtx.ContentRpcClient.GetVideoListByIdList(l.ctx, &pb.GetVideoListByIdListReq{
		VideoIdList: videoIdList,
	})
	if err1 != nil {
		return &types.FavoriteListResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
		}, err1
	}
	fmt.Println(videoListResp)
	fmt.Println("查到videoList")

	var authorList []*userPB.Userinfo
	for _, video := range videoListResp.GetVideoList() {
		userResp, err2 := l.svcCtx.UserRpcClient.GetUserById(l.ctx, &userPB.GetUserByIdReq{
			UserID: video.UserId,
		})
		if err2 != nil {
			return &types.FavoriteListResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			}, err2
		}
		authorList = append(authorList, userResp.GetUserinfo())
	}
	fmt.Println(authorList)
	fmt.Println("查到authorList")
	// 4.同时遍历 videoListResp 和 userListResp ，对每个 user 调用服务组合获取计数信息（点赞数、作品数、获赞数）并作为当前遍历到的 video 的成员，
	//   对每个 video 调用服务组合获取计数信息（点赞数、评论数），组装出完整的 video 对象，最终得到 videoList
	var videoList []types.Video
	for index, item := range videoListResp.GetVideoList() {
		// 4.1 将 user 的计数信息组装到 user 对象
		userLikeCntKey := constants.CntCacheUserLikePrefix + strconv.FormatInt(item.Id, 10)
		// 4.1.1 从 redis 中拿到用户点赞数
		userlikeCntRecord, err5 := redisClient.GetCtx(l.ctx, userLikeCntKey)
		if err5 != nil && err5 != redis.Nil {
			// 返回 redis 访问错误
			return &types.FavoriteListResp{
				StatusCode: common.REDIS_ERROR,
				StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			}, err5
		}
		var userLikeCnt int64 = 0
		if len(userlikeCntRecord) != 0 {
			userLikeCnt, _ = strconv.ParseInt(userlikeCntRecord, 10, 64)
		}
		// 4.1.2 根据 userId 调用一个服务组合查询 mysql 多张表，得到用户的发布作品数、获赞数
		cntResp, err6 := l.svcCtx.ContentRpcClient.GetUserPublishAndLikedCntById(l.ctx, &pb.GetUserPublishAndLikedCntByIdReq{
			UserId: authorList[index].Id,
		})
		if err6 != nil {
			return &types.FavoriteListResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			}, err6
		}
		fmt.Println(cntResp)
		fmt.Println("查到publishCnt、likedCnt")
		// 4.1.3 组装出完整的 user 对象
		user := types.User{
			Id:   item.Id,
			Name: authorList[index].Username,
			//FollowCount:     userResp.Userinfo.FollowCount,//关注数
			//FollowerCount:   userResp.Userinfo.FollowerCount,//被关注数
			//IsFollow:        userResp.Userinfo.IsFollow,//当前用户是否关注当前视频作者
			Avatar:          authorList[index].Avatar,
			BackgroundImage: authorList[index].BackgroundImage,
			Signature:       authorList[index].Signature,
			WorkCount:       cntResp.PublishCnt, //作品数量
			FavoriteCount:   userLikeCnt,        //点赞数量
			TotalFavorited:  cntResp.LikedCnt,   //获赞数量
		}
		// 4.2 将 user 和视频计数信息添加到 video 对象，并加入 videoList
		videoLikedCntKey := constants.CntCacheVideoLikedPrefix + strconv.FormatInt(item.Id, 10)
		videoCommentedCntKey := constants.CntCacheVideoCommentedPrefix + strconv.FormatInt(item.Id, 10)
		// 4.2.1 从 redis 获取当前 video 的点赞数
		videoLikedCntRecord, err3 := redisClient.GetCtx(l.ctx, videoLikedCntKey)
		if err3 != nil && err3 != redis.Nil {
			// 返回 redis 访问错误
			return &types.FavoriteListResp{
				StatusCode: common.REDIS_ERROR,
				StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			}, err3
		}
		var videoLikedCnt int64 = 0
		if len(videoLikedCntRecord) != 0 {
			videoLikedCnt, _ = strconv.ParseInt(videoLikedCntRecord, 10, 64)
		}
		// 4.2.2 从 redis 获取当前 video 的评论数
		videoCommentedCntRecord, err4 := redisClient.GetCtx(l.ctx, videoCommentedCntKey)
		if err4 != nil && err4 != redis.Nil {
			// 返回 redis 访问错误
			return &types.FavoriteListResp{
				StatusCode: common.REDIS_ERROR,
				StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			}, err4
		}
		var videoCommentedCnt int64 = 0
		if len(videoLikedCntRecord) != 0 {
			videoCommentedCnt, _ = strconv.ParseInt(videoCommentedCntRecord, 10, 64)
		}
		// 4.2.3 组装出完整的 video 对象，加入 videoList
		videoList = append(videoList, types.Video{
			Id:            item.Id,
			Author:        user,
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: videoLikedCnt,     //视频点赞数
			CommentCount:  videoCommentedCnt, //视频评论数
			IsFavorite:    true,
			Title:         item.Title,
		})
	}

	//// 2.遍历 favoriteList 中信息查询 video 表，返回用户点赞的所有视频
	//fmt.Println(favoriteList)
	//fmt.Println("查到favoriteList")
	//for _, item := range favoriteList { //TODO (videoid_list)->sql
	//	// 2.1 根据 favoriteList 中 video_id 查询 video 信息
	//	videoResp, err1 := l.svcCtx.ContentRpcClient.GetVideoById(l.ctx, &pb.GetVideoByIdReq{
	//		Id: item.VideoId,
	//	})
	//	if err1 != nil {
	//		return &types.FavoriteListResp{
	//			StatusCode: common.DB_ERROR,
	//			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
	//		}, err1
	//	}
	//	fmt.Println(videoResp)
	//	fmt.Println("查到video") //TODO (userid_list)->sql
	//	// 2.2 根据 video 中作者的 user_id 查询 user 信息
	//	userResp, err2 := l.svcCtx.UserRpcClient.GetUserById(l.ctx, &userPB.GetUserByIdReq{
	//		UserID: videoResp.Video.UserId,
	//	})
	//	if err2 != nil {
	//		return &types.FavoriteListResp{
	//			StatusCode: common.DB_ERROR,
	//			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
	//		}, err2
	//	}
	//	fmt.Println(userResp)
	//	fmt.Println("查到user")
	//	// 2.3 将统计信息添加到 user 对象
	//	userLikeCntKey := "CntCache:user_likeCnt:" + strconv.FormatInt(userResp.Userinfo.Id, 10)
	//	// 2.3.1 从 redis 获取该视频作者 user 的点赞数量
	//	userlikeCntRecord, err5 := redisClient.GetCtx(l.ctx, userLikeCntKey)
	//	if err5 != nil && err5 != redis.Nil {
	//		// 返回 redis 访问错误
	//		return &types.FavoriteListResp{
	//			StatusCode: common.DB_ERROR,
	//			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
	//		}, err5
	//	}
	//	var userLikeCnt int64 = 0
	//	if len(userlikeCntRecord) != 0 {
	//		userLikeCnt, _ = strconv.ParseInt(userlikeCntRecord, 10, 64)
	//	}
	//	// 2.3.2 从数据库查找该视频作者 user 的获赞数量
	//	LikedCntResp, err6 := l.svcCtx.ContentRpcClient.GetUserFavoritedCnt(l.ctx, &pb.GetUserFavoritedCntByIdReq{
	//		UserID: videoResp.Video.UserId,
	//	})
	//	if err6 != nil {
	//		return &types.FavoriteListResp{
	//			StatusCode: common.DB_ERROR,
	//			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
	//		}, err6
	//	}
	//
	//	user := types.User{
	//		Id:   userResp.Userinfo.Id,
	//		Name: userResp.Userinfo.Username,
	//		//FollowCount:     userResp.Userinfo.FollowCount,//关注数
	//		//FollowerCount:   userResp.Userinfo.FollowerCount,//被关注数
	//		//IsFollow:        userResp.Userinfo.IsFollow,//当前用户是否关注当前视频作者
	//		Avatar:          userResp.Userinfo.Avatar,
	//		BackgroundImage: userResp.Userinfo.BackgroundImage,
	//		Signature:       userResp.Userinfo.Signature,
	//		//TODO WorkCount:      ,//作品数量
	//		FavoriteCount:  userLikeCnt,      //点赞数量
	//		TotalFavorited: LikedCntResp.Cnt, //获赞数量
	//	}
	//	// 2.4 将统计信息添加到 video 对象，并加入 videoList
	//	videoLikedCntKey := "CntCache:video_likedCnt:" + strconv.FormatInt(videoResp.Video.Id, 10)
	//	videoCommentedCntKey := "CntCache:video_commentedCnt:" + strconv.FormatInt(videoResp.Video.Id, 10)
	//	// 2.4.1 从 redis 获取当前 video 的点赞数
	//	videoLikedCntRecord, err3 := redisClient.GetCtx(l.ctx, videoLikedCntKey)
	//	if err3 != nil && err3 != redis.Nil {
	//		// 返回 redis 访问错误
	//		return &types.FavoriteListResp{
	//			StatusCode: common.DB_ERROR,
	//			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
	//		}, err3
	//	}
	//	var videoLikedCnt int64 = 0
	//	if len(videoLikedCntRecord) != 0 {
	//		videoLikedCnt, _ = strconv.ParseInt(videoLikedCntRecord, 10, 64)
	//	}
	//	// 2.4.2 从 redis 获取当前 video 的评论数
	//	videoCommentedCntRecord, err4 := redisClient.GetCtx(l.ctx, videoCommentedCntKey)
	//	if err4 != nil && err4 != redis.Nil {
	//		// 返回 redis 访问错误
	//		return &types.FavoriteListResp{
	//			StatusCode: common.DB_ERROR,
	//			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
	//		}, err4
	//	}
	//	var videoCommentedCnt int64 = 0
	//	if len(videoLikedCntRecord) != 0 {
	//		videoCommentedCnt, _ = strconv.ParseInt(videoCommentedCntRecord, 10, 64)
	//	}
	//	videoList = append(videoList, types.Video{
	//		Id:            videoResp.Video.Id,
	//		Author:        user,
	//		PlayUrl:       videoResp.Video.PlayUrl,
	//		CoverUrl:      videoResp.Video.CoverUrl,
	//		FavoriteCount: videoLikedCnt,     //视频点赞数
	//		CommentCount:  videoCommentedCnt, //视频评论数
	//		IsFavorite:    true,
	//		Title:         videoResp.Video.Title,
	//	})
	//}

	fmt.Println("【api-查询用户点赞列表成功】")
	fmt.Print(videoList)
	return &types.FavoriteListResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		VideoList:  videoList,
	}, nil
}
