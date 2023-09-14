package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"errors"

	"fmt"

	"time"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelFavoriteLogic {
	return &DelFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelFavoriteLogic) DelFavorite(in *pb.DelFavoriteReq) (*pb.DelFavoriteResp, error) {
	/*
		Author：    刘洋
		Function：  从 favorite 删除点赞信息
		Update：    08.28 对进入逻辑 加log
	*/
	l.Logger.Info("DelFavorite方法请求参数：", in)
	//1.根据（userId、videoId）查找 favorite 表
	favorite, err0 := l.svcCtx.FavoriteModel.FindFavoriteByUserIdVideoId(l.ctx, in.UserId, in.VideoId)
	if err0 != nil && !errors.Is(err0, model.ErrNotFound) {
		return nil, errors.New("rpc-delFavorite-数据查询失败")
	}
	if favorite == nil || favorite.IsDelete == 1 {
		return nil, errors.New("rpc-delFavorite-没有找到该条点赞数据")
	}
	//2.逻辑删除，置 isDelete=1 选项到 favorite 表项
	err1 := l.svcCtx.FavoriteModel.Update(l.ctx, &model.Favorite{
		Id:         favorite.Id,
		UserId:     in.UserId,
		VideoId:    in.VideoId,
		UpdateTime: time.Now(),
		IsDelete:   1,
	})
	if err1 != nil {
		return nil, errors.New("rpc-delFavorite-删除点赞数据失败")
	}
	fmt.Println("【rpc-delFavorite-删除点赞数据成功】")
	return &pb.DelFavoriteResp{}, nil
}
