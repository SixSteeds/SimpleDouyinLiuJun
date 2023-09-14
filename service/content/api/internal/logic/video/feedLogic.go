package video

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"
	"doushen_by_liujun/service/content/rpc/pb"
	userPb "doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedReq) (resp *types.FeedResp, err error) {

	l.Logger.Info("Feed方法请求参数：", req)

	var userId int64
	token, err := util.ParseToken(req.Token)
	if err != nil {
		l.Logger.Error(err)
		// 用户未登录
		userId = 0
	} else {
		userId = token.UserID
	}
	var lastTime int64
	if req.LatestTime > 169268692200 {
		// 获取第三位数字
		thirdDigit := (req.LatestTime / 100) % 10

		// 进行四舍五入
		if thirdDigit >= 5 {
			lastTime = req.LatestTime/1000 + 1
		} else {
			lastTime = req.LatestTime / 1000
		}
	} else {
		lastTime = req.LatestTime
	}
	data, err := l.svcCtx.ContentRpcClient.GetFeedList(l.ctx, &pb.FeedListReq{
		UserId:     userId,
		LatestTime: lastTime,
		Size:       5,
	})
	if err != nil {
		l.Logger.Error(err)
		return &types.FeedResp{
			StatusCode: common.DbError,
			StatusMsg:  common.MapErrMsg(common.DbError),
		}, nil
	}
	if data == nil {

		l.Logger.Error(err)

		return &types.FeedResp{
			StatusCode: common.DataUseUp,
			StatusMsg:  common.MapErrMsg(common.DataUseUp),
		}, nil
	}

	var videoList []*pb.FeedVideo
	var userIds []int64
	var listLen int
	if data == nil {
		listLen = 0
	} else {
		videoList = data.VideoList
		listLen = len(data.VideoList)
		userIds = data.UserIds
	}
	if listLen < 5 { //tzx新增，使发布后的视频循环播放，不会出现数据库繁忙的报错
		data2, err := l.svcCtx.ContentRpcClient.GetFeedList(l.ctx, &pb.FeedListReq{ //从头查没查完的
			UserId:     userId,
			LatestTime: int64(9999999999),
			Size:       int64(5 - listLen),
		})
		//陶子勋收到的数据！！！！！！
		if err != nil {
			l.Logger.Error(err)
			return &types.FeedResp{
				StatusCode: common.DbError,
				StatusMsg:  common.MapErrMsg(common.DbError),
			}, nil
		}
		if data2 == nil {

			l.Logger.Error(err)

			return &types.FeedResp{
				StatusCode: common.DataUseUp,
				StatusMsg:  common.MapErrMsg(common.DataUseUp),
			}, nil
		}
		if listLen == 0 {
			videoList = data.VideoList
			userIds = data2.UserIds
		} else {
			videoList = append(videoList, data2.VideoList...)
			userIds = append(userIds, data2.UserIds...)
		}
	}
	// 通过userIds获取到所有的user信息
	usersByIds, err := l.svcCtx.UserRpcClient.GetUsersByIds(l.ctx, &userPb.GetUsersByIdsReq{
		Ids:    userIds,
		UserID: userId,
	})
	if err != nil {
		l.Logger.Error(err)
		return &types.FeedResp{
			StatusCode: common.DbError,
			StatusMsg:  common.MapErrMsg(common.DbError),
		}, nil
	}

	var feedUserList []*pb.FeedUser
	for _, user := range usersByIds.Users {
		feedUserList = append(feedUserList, &pb.FeedUser{
			Id:              user.Id,
			Name:            user.Name,
			FollowCount:     &user.FollowCount,
			FollowerCount:   &user.FollowerCount,
			IsFollow:        user.IsFollow,
			Avatar:          &user.Avatar,
			BackgroundImage: &user.BackgroundImage,
			Signature:       &user.Signature,
			TotalFavorited:  &user.TotalFavorited,
			WorkCount:       &user.WorkCount,
			FavoriteCount:   &user.FavoriteCount,
		})
	}

	var FeedVideos []types.Video
	var nextTime int64

	for count, video := range videoList {
		user := feedUserList[count]
		var author = &types.User{
			Id:              user.Id,
			Name:            user.Name,
			FollowCount:     *user.FollowCount,
			FollowerCount:   *user.FollowerCount,
			IsFollow:        user.IsFollow,
			Avatar:          *user.Avatar,
			BackgroundImage: *user.BackgroundImage,
			Signature:       *user.Signature,
			TotalFavorited:  *user.TotalFavorited,
			WorkCount:       *user.WorkCount,
			FavoriteCount:   *user.FavoriteCount,
		}
		FeedVideos = append(FeedVideos, types.Video{
			Id:            video.Id,
			Author:        *author,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		})
		nextTime = video.NextTime
	}

	return &types.FeedResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		VideoList:  FeedVideos,
		NextTime:   nextTime,
	}, nil
}
