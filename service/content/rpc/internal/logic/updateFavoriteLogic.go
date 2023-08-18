package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"errors"
	"log"
	"time"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteLogic {
	return &UpdateFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFavoriteLogic) UpdateFavorite(in *pb.UpdateFavoriteReq) (*pb.UpdateFavoriteResp, error) {

	//1.根据传入的 isDelete 修改 favorite 表
	err := l.svcCtx.FavoriteModel.Update(l.ctx, &model.Favorite{
		Id:         in.Id,
		UpdateTime: time.Now(),
		IsDelete:   in.IsDelete,
	})
	if err != nil {
		return nil, errors.New("rpc-updateFavorite-修改点赞信息失败")
	}
	if err := l.svcCtx.KqPusherClient.Push("content_rpc_updateFavoriteLogic_UpdateFavorite_Update_false"); err != nil {
		log.Fatal(err)
	}

	if in.IsDelete == 0 {
		logx.Error("rpc-updateFavorite-修改点赞记录为逻辑点赞成功")
	} else {
		logx.Error("rpc-updateFavorite-修改点赞记录为逻辑删除成功")
	}
	if err := l.svcCtx.KqPusherClient.Push("content_rpc_updateFavoriteLogic_UpdateFavorite_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.UpdateFavoriteResp{}, nil
}
