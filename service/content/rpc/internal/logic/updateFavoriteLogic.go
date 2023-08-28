package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"errors"
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
	/*
		Author：    刘洋
		Function：  向 favorite 表 根据表项id更新 isDelete参数
		Update：    08.28 对进入逻辑 加log
	*/
	l.Logger.Info("UpdateFavorite方法请求参数：", in)
	//1.根据传入的 isDelete 修改 favorite 表
	err := l.svcCtx.FavoriteModel.Update(l.ctx, &model.Favorite{
		Id:         in.Id,
		UpdateTime: time.Now(),
		IsDelete:   in.IsDelete,
	})
	if err != nil {
		return nil, errors.New("rpc-updateFavorite-修改点赞信息失败")
	}

	if in.IsDelete == 0 {
		logx.Error("rpc-updateFavorite-修改点赞记录为逻辑点赞成功")
	} else {
		logx.Error("rpc-updateFavorite-修改点赞记录为逻辑删除成功")
	}
	return &pb.UpdateFavoriteResp{}, nil
}
