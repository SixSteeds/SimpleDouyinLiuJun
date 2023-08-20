package logic

import (
	"context"
	userPb "doushen_by_liujun/service/user/rpc/pb"
	"fmt"

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
	// todo: add your logic here and delete this line

	fmt.Println("进入GetPublishList rpc逻辑")
	feedList, err := l.svcCtx.VideoModel.GetPublishList(l.ctx, in.UserId, in.CheckUserId)
	fmt.Println(feedList)
	if err != nil {
		return nil, err
	}
	// 将feedlist中的userId全部拿出来转换为一个数组
	var userIds []int64
	for _, feed := range *feedList {
		userIds = append(userIds, feed.UserId)
	}

	// 通过userIds获取到所有的user信息
	usersByIds, err := l.svcCtx.UserRpcClient.GetUsersByIds(l.ctx, &userPb.GetUsersByIdsReq{
		Ids:    userIds,
		UserID: in.UserId,
	})
	fmt.Println(usersByIds)

	fmt.Println("完成feed流rpc逻辑11111111111")
	var FeedVideos []*pb.FeedVideo
	var feedUserList []pb.FeedUser
	for _, user := range usersByIds.Users {
		feedUserList = append(feedUserList, pb.FeedUser{
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
			Author:        &feedUserList[count],
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

	fmt.Println("获取到的feed流信息为：", FeedVideos)
	return &pb.PublishListResp{
		VideoList: FeedVideos,
	}, nil
}
