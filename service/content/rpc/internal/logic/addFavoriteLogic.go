package logic

import (
	"context"
	constants "doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"errors"
	"fmt"
	"time"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFavoriteLogic {
	return &AddFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddFavoriteLogic) AddFavorite(in *pb.AddFavoriteReq) (*pb.AddFavoriteResp, error) {
	/*
		Author：    刘洋
		Function：  向 favorite 表添加点赞信息
		Update：    08.28 对进入逻辑 加log
	*/
	l.Logger.Info("AddFavorite方法请求参数：", in)
	//1.根据（userId、videoId）查找 favorite 表
	favorite, err0 := l.svcCtx.FavoriteModel.FindFavoriteByUserIdVideoId(l.ctx, in.UserId, in.VideoId)
	fmt.Println(favorite)
	if err0 != nil && !errors.Is(err0, model.ErrNotFound) {
		return nil, errors.New("rpc-AddFavorite-数据查询失败")
	}
	//2.favorite记录存在，则置 isDelete=0 选项到 favorite 表项
	if favorite != nil {
		fmt.Println("查到")
		err := l.svcCtx.FavoriteModel.Update(l.ctx, &model.Favorite{
			Id:         favorite.Id,
			UserId:     in.UserId,
			VideoId:    in.VideoId,
			UpdateTime: time.Now(),
			IsDelete:   0,
		})
		if err != nil {
			return nil, errors.New("rpc-AddFavorite-新增点赞数据失败")
		}
	} else {
		//3.favorite记录不存在，则新增点赞信息到 favorite 表项
		fmt.Println("没查到")
		//雪花算法生成id
		snowflake, err1 := util.NewSnowflake(constants.ContentRpcMachineId)
		if err1 != nil {
			return nil, errors.New("rpc-AddFavorite-新增评论，snowflake生成id失败")
		}
		snowId := snowflake.Generate()
		_, err := l.svcCtx.FavoriteModel.Insert(l.ctx, &model.Favorite{
			Id:         snowId,
			UserId:     in.UserId,
			VideoId:    in.VideoId,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			IsDelete:   0,
		})
		if err != nil {
			return nil, errors.New("rpc-AddFavorite-新增点赞数据失败")
		}
	}

	fmt.Println("【rpc-AddFavorite-新增点赞数据成功】")
	return &pb.AddFavoriteResp{}, nil
}
