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

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListReq) (resp *types.PublishListResp, err error) {
	l.Logger.Info("PublishList方法请求参数：", req)
	var userId int64
	token, err := util.ParseToken(req.Token)
	if err != nil {
		// 用户未登录（此处避免抖声app在退出登录时乱发请求所做限流）
		return &types.PublishListResp{
			StatusCode: common.OK,
			StatusMsg:  common.MapErrMsg(common.OK),
			VideoList:  nil,
		}, nil
	} else {
		userId = token.UserID
	}
	data, err := l.svcCtx.ContentRpcClient.GetPublishList(l.ctx, &pb.PublishListReq{
		CheckUserId: req.UserId,
		UserId:      userId,
	})
	if err != nil {
		l.Logger.Error(err)
		return &types.PublishListResp{
			StatusCode: common.DbError,
			StatusMsg:  common.MapErrMsg(common.DbError),
		}, nil
	}
	if data == nil {
		l.Logger.Error(err)
		return &types.PublishListResp{
			StatusCode: common.OK,
			StatusMsg:  common.MapErrMsg(common.OK),
		}, nil
	}

	// 通过userIds获取到所有的user信息
	usersByIds, err := l.svcCtx.UserRpcClient.GetUsersByIds(l.ctx, &userPb.GetUsersByIdsReq{
		Ids:    data.UserIds,
		UserID: userId,
	})
	if err != nil {
		l.Logger.Error(err)
		return &types.PublishListResp{
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

	videoList := data.VideoList

	var FeedVideos []types.Video

	for index, video := range videoList {
		user := feedUserList[index]
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
	}
	return &types.PublishListResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		VideoList:  FeedVideos,
	}, nil
}
