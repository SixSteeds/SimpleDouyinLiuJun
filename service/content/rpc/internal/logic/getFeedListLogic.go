package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"
	"fmt"
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
	// todo: add your logic here and delete this line
	fmt.Println("陶子勋收到的数据！！！！！！GetFeedList", in.UserId, in.LatestTime, in.Size)
	fmt.Println("进入feed流rpc逻辑")
	feedList, err := l.svcCtx.VideoModel.GetFeedList(l.ctx, in.UserId, &in.LatestTime, in.Size)
	fmt.Println("在rpc里看model返回值")
	fmt.Println(feedList)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	var FeedVideos []*pb.FeedVideo
	if len(*feedList) == 0 {
		fmt.Println("model查询为空")
		return &pb.FeedListResp{
			VideoList: FeedVideos,
		}, nil
	}

	// 将feedlist中的userId全部拿出来转换为一个数组 1692684
	var userIds []int64
	for _, feed := range *feedList {
		IsFavorite, _ := l.svcCtx.RedisClient.HgetCtx(l.ctx, common.LikeCacheVideoLikedPrefix+strconv.FormatInt(feed.Id, 10), strconv.FormatInt(feed.UserId, 10))
		fmt.Println("666666666666666666666666666666")
		if len(IsFavorite) != 0 {
			if IsFavorite == "0" {
				feed.IsFavorite = true
			}
		}
		fmt.Println("===========================================")
		fmt.Println(feed.UserId)
		userIds = append(userIds, feed.UserId)
	}

	fmt.Println("完成feed流rpc逻辑11111111111")

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

	fmt.Println("tttttttt获取到的feed流信息为：", FeedVideos)
	return &pb.FeedListResp{
		VideoList: FeedVideos,
		UserIds:   userIds,
	}, nil
}
