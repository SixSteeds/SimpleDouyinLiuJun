package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"errors"
	"fmt"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPublishAndLikedCntByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserPublishAndLikedCntByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPublishAndLikedCntByIdLogic {
	return &GetUserPublishAndLikedCntByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserPublishAndLikedCntByIdLogic) GetUserPublishAndLikedCntById(in *pb.GetUserPublishAndLikedCntByIdReq) (*pb.GetUserPublishAndLikedCntByIdResp, error) {

	/*
		Author：    刘洋
		Function：  从 video表和favorite表 查找用户发布作品数、用户被点赞总数
		Update：    08.28 对进入逻辑 加log
	*/
	l.Logger.Info("GetUserPublishAndLikedCntById方法请求参数：", in)
	//1. 根据 userId 查找 video 表，得到用户发布的所有 videoId, 并计数 publishCnt
	var publishCnt int64 = 0
	videoList, err := l.svcCtx.VideoModel.FindVideoListByUserId(l.ctx, in.UserId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.New("数据查询失败")
	}
	var videoIdList []int64
	for _, item := range *videoList {
		videoIdList = append(videoIdList, item.Id)
		publishCnt++
	}
	fmt.Println(videoIdList)
	fmt.Println("查到所有videoId")
	if len(videoIdList) == 0 {
		return &pb.GetUserPublishAndLikedCntByIdResp{
			PublishCnt: publishCnt,
			LikedCnt:   0,
		}, nil
	}
	//2. 根据 videoIdList 查找 favorite 表，count 得到这些所有作品的总获赞数
	var likedCnt int64 = 0
	likedCntResp, err2 := l.svcCtx.FavoriteModel.FindFavoritedCntByVideoIdList(l.ctx, &videoIdList)
	if err2 != nil && !errors.Is(err2, model.ErrNotFound) {
		return nil, errors.New("数据查询失败")
	}
	likedCnt = likedCntResp
	fmt.Println("【rpc-GetUserPublishAndLikedCntById-查询计数数据成功】")
	return &pb.GetUserPublishAndLikedCntByIdResp{
		PublishCnt: publishCnt,
		LikedCnt:   likedCnt,
	}, nil
}
