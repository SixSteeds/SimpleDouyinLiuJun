package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"
	userPb "doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedListLogic {
	return &GetFeedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFeedListLogic) GetFeedList(in *pb.FeedListReq) (*pb.FeedListResp, error) {
	// todo: add your logic here and delete this line
	feedList, err := l.svcCtx.VideoModel.GetFeedList(l.ctx, in.UserId, in.LatestTime, in.Size)
	if err != nil {
		return nil, err
	}
	// 将feedlist中的userId全部拿出来转换为一个数组
	var userIds []int64
	for _, feed := range feedList {
		userIds = append(userIds, feed.UserId)
	}
	// 通过userIds获取到所有的user信息
	usersByIds, err := l.svcCtx.UserRpcClient.GetUsersByIds(l.ctx, &userPb.GetUsersByIdsReq{
		Ids:    userIds,
		UserID: in.UserId,
	})
	var FeedVideos []*pb.FeedVideo
	var feedUserList []pb.FeedUser
	for _, user := range usersByIds.Users {
		feedUserList = append(feedUserList, pb.FeedUser{
			Id:              user.Id,
			Name:            user.Username,
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

	for count, feed := range feedList {
		FeedVideos = append(FeedVideos, &pb.FeedVideo{
			Id:            feed.Id,
			Author:        &feedUserList[count],
			PlayUrl:       feed.PlayUrl,
			CoverUrl:      feed.CoverUrl,
			Title:         feed.Title,
			FavoriteCount: feed.FavoriteCount,
			CommentCount:  feed.CommentCount,
			IsFavorite:    feed.IsFavorite == 1,
		})
	}

	if err != nil {
		return nil, err
	}

	return &pb.FeedListResp{}, nil
}
