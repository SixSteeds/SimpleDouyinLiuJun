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
	fmt.Println(req.Token)
	var userId int64
	token, err := util.ParseToken(req.Token)
	if err != nil {
		// 用户未登录
		userId = 0
	} else {
		userId = token.UserID
	}
	var lastTime int64
	if req.LatestTime > 169268692200 {
		lastTime = req.LatestTime / 1000
	} else {
		lastTime = req.LatestTime
	}
	fmt.Println("我要去rpc了")
	list, err := l.svcCtx.ContentRpcClient.GetFeedList(l.ctx, &pb.FeedListReq{
		UserId:     userId,
		LatestTime: lastTime,
		Size:       5,
	})
	fmt.Println("我回到api了，我来看看list是什么")
	fmt.Println(list)
	fmt.Println(err)

	if err != nil {
		fmt.Println("陶子勋陶子勋陶子勋，listerr", err)
		return &types.FeedResp{
			StatusCode: common.DB_ERROR,
			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
		}, nil
	}
	if list == nil {
		fmt.Println("陶子勋陶子勋陶子勋，listnil")
		return &types.FeedResp{
			StatusCode: common.DATA_USE_UP,
			StatusMsg:  common.MapErrMsg(common.DATA_USE_UP),
		}, nil
	}

	var videoList []*pb.FeedVideo
	var listLen int
	if list == nil {
		listLen = 0
	} else {
		videoList = list.VideoList
		listLen = len(list.VideoList)
	}

	if listLen < 5 { //tzx新增，使发布后的视频循环播放，不会出现数据库繁忙的报错
		list2, err := l.svcCtx.ContentRpcClient.GetFeedList(l.ctx, &pb.FeedListReq{ //从头查没查完的
			UserId:     userId,
			LatestTime: int64(9999999999),
			Size:       int64(5 - listLen),
		})
		//陶子勋收到的数据！！！！！！ 877806281992183808 1111111111 1
		if err != nil {
			fmt.Println("陶子勋陶子勋陶子勋，list2err", err)
			return &types.FeedResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			}, nil
		}
		if list2 == nil {
			fmt.Println("陶子勋陶子勋陶子勋，list2nil")
			return &types.FeedResp{
				StatusCode: common.DATA_USE_UP,
				StatusMsg:  common.MapErrMsg(common.DATA_USE_UP),
			}, nil
		}
		if listLen == 0 {
			videoList = list2.VideoList
		} else {
			videoList = append(videoList, list2.VideoList...)
		}
	}
	fmt.Println("陶子勋陶子勋陶子勋，视频流长度", len(videoList))
	fmt.Println("完成feed流rpc逻辑")
	var FeedVideos []types.Video
	fmt.Println("到这了222222")
	var nextTime int64
	//fmt.Println("到这了33333")
	fmt.Println(list)

	for _, video := range videoList {
		//fmt.Println("到这了44444")
		user := video.Author
		//在这打印一下author吧，
		//fmt.Println("到这了11111111111")
		//fmt.Println(user)
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
	//fmt.Println("完成对象转换逻辑")

	return &types.FeedResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		VideoList:  FeedVideos,
		NextTime:   nextTime,
	}, nil
}
