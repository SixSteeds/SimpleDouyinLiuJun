package logic

import (
	"context"
	"doushen_by_liujun/service/media/rpc/internal/model"
	"doushen_by_liujun/service/media/rpc/internal/svc"
	"doushen_by_liujun/service/media/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type SaveVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveVideoLogic {
	return &SaveVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveVideoLogic) SaveVideo(in *pb.SaveVideoReq) (*pb.SaveVideoResp, error) {
	l.Logger.Info(in)
	_, err := l.svcCtx.MediaModel.Save(l.ctx, &model.Video{
		Id:       in.Id,
		UserId:   in.UserId,
		CoverUrl: in.CoverUrl,
		PlayUrl:  in.PlayUrl,
		Title:    in.Title,
	})
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return &pb.SaveVideoResp{}, nil
}
