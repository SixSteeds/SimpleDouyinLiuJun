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

type GetUserFavoritedCntLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFavoritedCntLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFavoritedCntLogic {
	return &GetUserFavoritedCntLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFavoritedCntLogic) GetUserFavoritedCnt(in *pb.GetUserFavoritedCntByIdReq) (*pb.GetUserFavoritedCntByIdResp, error) {
	var count int64 = 0
	//1. 根据 userId 查找用户发布的所有 videoId
	videoList, err := l.svcCtx.VideoModel.FindVideoListByUserId(l.ctx, in.UserID)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.New("数据查询失败")
	}
	for _, video := range *videoList { //todo (video_list)->
		fmt.Println(video)
		//2. 根据 videoId 查找所有的点赞记录，累加点赞数
		favoriteList, err2 := l.svcCtx.FavoriteModel.FindFavoriteListByVideoId(l.ctx, video.Id)
		if err2 != nil && err2 != model.ErrNotFound {
			return nil, errors.New("数据查询失败")
		}
		for _, like := range *favoriteList {
			if like.IsDelete == 0 { //逻辑删除的不返回给api
				count += 1
			}
		}
		//TODO 可以改为 count 查询
	}
	return &pb.GetUserFavoritedCntByIdResp{
		Cnt: count,
	}, nil
}
