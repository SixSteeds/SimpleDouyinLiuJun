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

type GetVideoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoByIdLogic {
	return &GetVideoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoByIdLogic) GetVideoById(in *pb.GetVideoByIdReq) (*pb.GetVideoByIdResp, error) {
	/*
		Author：    刘洋
		Function：  从 video 表 根据videoId查找video表项
		Update：    08.28 对进入逻辑 加log
	*/
	l.Logger.Info("GetVideoById方法请求参数：", in)
	videoInfo, err := l.svcCtx.VideoModel.FindOne(l.ctx, in.Id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.New("rpc-GetVideoById-数据查询失败")
	}
	if videoInfo != nil {
		fmt.Println("查到")
		return &pb.GetVideoByIdResp{
			Video: &pb.Video{
				Id:       videoInfo.Id,
				UserId:   videoInfo.UserId,
				PlayUrl:  videoInfo.PlayUrl,
				CoverUrl: videoInfo.CoverUrl.String,
				Title:    videoInfo.Title.String,
				//CreateTime: videoInfo.Title,
				//UpdateTime: videoInfo.UpdateTime,
			},
		}, nil
	}
	fmt.Println("没查到")
	return &pb.GetVideoByIdResp{}, nil
}
