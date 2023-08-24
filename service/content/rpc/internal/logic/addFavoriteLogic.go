package logic

import (
	"context"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"errors"
	"fmt"
	"log"
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

	//1.根据（userId、videoId）查找 favorite 表
	favorite, err0 := l.svcCtx.FavoriteModel.FindFavoriteByUserIdVideoId(l.ctx, in.UserId, in.VideoId)
	fmt.Println(favorite)
	if err0 != nil && err0 != model.ErrNotFound {
		return nil, errors.New("rpc-AddFavorite-数据查询失败")
	}
	if err := l.svcCtx.KqPusherClient.Push("content_rpc_addFavoriteLogic_AddFavorite_FindFavoriteByUserIdVideoId_false"); err != nil {
		log.Fatal(err)
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
		if err := l.svcCtx.KqPusherClient.Push("content_rpc_addFavoriteLogic_AddFavorite_Update_false"); err != nil {
			log.Fatal(err)
		}
	} else {
		//3.favorite记录不存在，则新增点赞信息到 favorite 表项
		fmt.Println("没查到")
		//雪花算法生成id
		snowflake, err1 := util.NewSnowflake(3)
		if err1 != nil {
			return nil, errors.New("rpc-AddFavorite-新增评论，snowflake生成id失败")
		}
		if err := l.svcCtx.KqPusherClient.Push("content_rpc_addFavoriteLogic_AddFavorite_NewSnowflake_false"); err != nil {
			log.Fatal(err)
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
		if err := l.svcCtx.KqPusherClient.Push("content_rpc_addFavoriteLogic_AddFavorite_Insert_false"); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("【rpc-AddFavorite-新增点赞数据成功】")
	logx.Error("rpc-AddFavorite-新增点赞数据成功")
	if err := l.svcCtx.KqPusherClient.Push("content_rpc_addFavoriteLogic_AddFavorite_success"); err != nil {
		log.Fatal(err)
	}

	return &pb.AddFavoriteResp{}, nil
}
