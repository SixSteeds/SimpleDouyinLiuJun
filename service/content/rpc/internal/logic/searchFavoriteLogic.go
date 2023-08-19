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

type SearchFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchFavoriteLogic {
	return &SearchFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchFavoriteLogic) SearchFavorite(in *pb.SearchFavoriteReq) (*pb.SearchFavoriteResp, error) {

	// 1.根据 user_id 查询 favorite 表，返回所有点赞信息
	favoriteList, err := l.svcCtx.FavoriteModel.FindFavoriteListByUserId(l.ctx, in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.New("数据查询失败")
	}
	var resp []*pb.Favorite
	for _, item := range *favoriteList {
		if item.IsDelete == 0 { //逻辑删除的不返回给api
			resp = append(resp, &pb.Favorite{
				Id:         item.Id,
				VideoId:    item.VideoId,
				UserId:     item.UserId,
				CreateTime: item.CreateTime.Unix(),
				UpdateTime: item.UpdateTime.Unix(),
			})
		}
	}
	fmt.Println("【rpc-SearchFavorite-查询用户点赞列表成功】")
	return &pb.SearchFavoriteResp{
		Favorite: resp,
	}, nil
}
