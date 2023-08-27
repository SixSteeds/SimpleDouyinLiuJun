package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"fmt"
	"strconv"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublishListLogic {
	return &GetPublishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPublishListLogic) GetPublishList(in *pb.PublishListReq) (*pb.PublishListResp, error) {
	l.Logger.Info("GetPublishList方法请求参数：", in)

	feedList, err := l.svcCtx.VideoModel.GetPublishList(l.ctx, in.UserId, in.CheckUserId)
	//fmt.Println(feedList)
	if err != nil {
		return nil, err
	}
	if len(*feedList) == 0 {
		return nil, nil
	}

	// 将feedlist中的userId全部拿出来转换为一个数组
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

	var FeedVideos []*pb.FeedVideo

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

	//fmt.Println("获取到的feed流信息为：", FeedVideos)
	fmt.Println("退出GetPublishList rpc逻辑")
	return &pb.PublishListResp{
		VideoList: FeedVideos,
		UserIds:   userIds,
	}, nil
}
