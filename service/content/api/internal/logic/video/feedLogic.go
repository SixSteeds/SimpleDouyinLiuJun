package video

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/rpc/pb"
	"fmt"

	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"

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
	// todo: add your logic here and delete this line
	fmt.Println("进入feed流api逻辑")
	fmt.Println(req.LatestTime)
	var userId int64
	token, err := util.ParseToken(req.Token)
	if err != nil {
		// 用户未登录
		userId = 0
	} else {
		userId = token.UserID
	}
	list, err := l.svcCtx.ContentRpcClient.GetFeedList(l.ctx, &pb.FeedListReq{

		UserId:     userId,
		LatestTime: req.LatestTime,
		Size:       5,
	})
	fmt.Println("到这了")
	fmt.Println(list)

	if err != nil {
		return &types.FeedResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
		}, nil
	}
	videoList := list.VideoList
	fmt.Println(videoList)
	if len(videoList) == 0 {
		return &types.FeedResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
		}, nil
	}

	fmt.Println("完成feed流rpc逻辑")
	var FeedVideos []types.Video
	fmt.Println("到这了222222")
	var nextTime int64
	fmt.Println("到这了33333")
	fmt.Println(list)

	for _, video := range videoList {
		fmt.Println("到这了44444")
		user := video.Author
		//在这打印一下author吧，
		fmt.Println("到这了11111111111")
		fmt.Println(user)
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
	fmt.Println("完成对象转换逻辑")

	return &types.FeedResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		VideoList:  FeedVideos,
		NextTime:   nextTime,
	}, nil
}
