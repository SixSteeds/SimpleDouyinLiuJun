package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"
	"strconv"

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
	l.Logger.Info("GetFeedList方法请求参数：", in)
	feedList, err := l.svcCtx.VideoModel.GetFeedList(l.ctx, in.UserId, &in.LatestTime, in.Size)
	if err != nil {
		return nil, err
	}
	var FeedVideos []*pb.FeedVideo
	if len(*feedList) == 0 {
		return &pb.FeedListResp{
			VideoList: FeedVideos,
		}, nil
	}

	// 将feedlist中的userId全部拿出来转换为一个数组 1692684
	var userIds []int64
	for _, feed := range *feedList {
		IsFavorite, _ := l.svcCtx.RedisClient.HgetCtx(l.ctx, common.LikeCacheVideoLikedPrefix+strconv.FormatInt(feed.Id, 10), strconv.FormatInt(feed.UserId, 10))
		if len(IsFavorite) != 0 {
			if IsFavorite == "0" {
				feed.IsFavorite = true
			}
		}
		userIds = append(userIds, feed.UserId)
	}
	for _, feed := range *feedList {
		FeedVideos = append(FeedVideos, &pb.FeedVideo{
			Id:            feed.Id,
			Author:        nil,
			PlayUrl:       feed.PlayUrl,
			CoverUrl:      feed.CoverUrl,
			Title:         feed.Title,
			FavoriteCount: feed.FavoriteCount,
			CommentCount:  feed.CommentCount,
			IsFavorite:    feed.IsFavorite,
			NextTime:      feed.CreateTime.Unix(),
		})
	}
	if err != nil {
		return nil, err
	}

	return &pb.FeedListResp{
		VideoList: FeedVideos,
		UserIds:   userIds,
	}, nil
}
