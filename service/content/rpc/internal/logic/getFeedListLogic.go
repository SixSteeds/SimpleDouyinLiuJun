package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"
	userPb "doushen_by_liujun/service/user/rpc/pb"
	"fmt"

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
		fmt.Println("===========================================")
		fmt.Println(feed.UserId)
		userIds = append(userIds, feed.UserId)
	}

	// 通过userIds获取到所有的user信息
	usersByIds, err := l.svcCtx.UserRpcClient.GetUsersByIds(l.ctx, &userPb.GetUsersByIdsReq{
		Ids:    userIds,
		UserID: in.UserId,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(usersByIds)

	fmt.Println("完成feed流rpc逻辑11111111111")
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
	for count, feed := range *feedList {
		FeedVideos = append(FeedVideos, &pb.FeedVideo{
			Id:            feed.Id,
			Author:        feedUserList[count],
			PlayUrl:       feed.PlayUrl,
			CoverUrl:      feed.CoverUrl,
			Title:         feed.Title,
			FavoriteCount: feed.FavoriteCount,
			CommentCount:  feed.CommentCount,
			IsFavorite:    feed.IsFavorite,
			NextTime:      feed.UpdateTime.Unix(),
		})
	}
	if err != nil {
		return nil, err
	}

	fmt.Println("tttttttt获取到的feed流信息为：", FeedVideos)
	return &pb.FeedListResp{
		VideoList: FeedVideos,
	}, nil
}
