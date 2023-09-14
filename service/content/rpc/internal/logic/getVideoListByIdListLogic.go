package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoListByIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoListByIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListByIdListLogic {
	return &GetVideoListByIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoListByIdListLogic) GetVideoListByIdList(in *pb.GetVideoListByIdListReq) (*pb.GetVideoListByIdListResp, error) {
	/*
		Author：    刘洋
		Function：  从 video 表拿出多个videoId对应的video数据
		Update：    08.28 对进入逻辑 加log
	*/
	l.Logger.Info("GetVideoListByIdList方法请求参数：", in)
	videoList, err := l.svcCtx.VideoModel.FindVideoListByIdList(l.ctx, &in.VideoIdList)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.New("rpc-GetVideoListByIdList-数据查询失败")
	}
	var resp []*pb.Video
	for _, item := range *videoList {
		resp = append(resp, &pb.Video{
			Id:         item.Id,
			UserId:     item.UserId,
			PlayUrl:    item.PlayUrl,
			CoverUrl:   item.CoverUrl.String,
			Title:      item.Title.String,
			CreateTime: item.CreateTime.Unix(),
			UpdateTime: item.UpdateTime.Unix(),
		})
	}
	fmt.Println("【rpc-GetVideoListByIdList-根据视频id列表查询视频列表成功】")
	return &pb.GetVideoListByIdListResp{
		VideoList: resp,
	}, nil
}
