package favorite

import (
	"context"
	"doushen_by_liujun/internal/common"
	constants "doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"
	"fmt"
	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
)

type FavoriteActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteActionLogic) FavoriteAction(req *types.FavoriteActionReq) (resp *types.FavoriteActionResp, err error) {
	/*
		Author：    刘洋
		Function：  点赞、取消点赞（ ActionType=1 点赞 ，ActionType=2 取消 ）
		Update：    08.28 对进入逻辑、异常 加log
	*/
	l.Logger.Info(req)
	// 1.根据 token 获取 userid
	parsToken, err0 := util.ParseToken(req.Token)
	if err0 != nil {
		l.Logger.Error(err0)
		return &types.FavoriteActionResp{
			StatusCode: common.TokenExpireError,
			StatusMsg:  common.MapErrMsg(common.TokenExpireError),
		}, nil
	}

	// 2.使用redis缓存
	redisClient := l.svcCtx.RedisClient
	videoLikedKey := constants.LikeCacheVideoLikedPrefix + strconv.FormatInt(req.VideoId, 10)
	videoLikedCntKey := constants.CntCacheVideoLikedPrefix + strconv.FormatInt(req.VideoId, 10)
	userLikeCntKey := constants.CntCacheUserLikePrefix + strconv.FormatInt(parsToken.UserID, 10)

	if action := req.ActionType; action == 1 { // actionType（1点赞，2取消）
		// 3.新增点赞
		// 3.1 查询 redis 点赞记录
		likeRecord, err1 := redisClient.HgetCtx(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10))
		if err1 != nil && err1 != redis.Nil {
			l.Logger.Error(err1)
			return &types.FavoriteActionResp{
				StatusCode: common.RedisError,
				StatusMsg:  common.MapErrMsg(common.RedisError),
			}, nil
		}
		if len(likeRecord) != 0 && likeRecord == "0" {
			logx.Error("api-favoriteAction-已点赞，重复操作无效")
		} else {
			// 新建 redis 连接
			c := red.NewClient(&red.Options{
				Addr:     "127.0.0.1:8094",
				Password: common.DefaultPass,
			})
			//一起执行 pipeline 事务操作
			pipeline := c.TxPipeline()
			// 3.2 新增 redis video 被点赞记录
			pipeline.HSet(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10), "0")
			// 3.3 redis 中 video 被点赞计数自增
			pipeline.Incr(l.ctx, videoLikedCntKey)
			// 3.4 redis 中 user 点赞计数自增
			pipeline.Incr(l.ctx, userLikeCntKey)
			// 3.5 pipeline 执行
			_, e := pipeline.Exec(l.ctx)

			//一起执行 pipeline 操作
			//e := redisClient.PipelinedCtx(l.ctx, func(pipeline redis.Pipeliner) error {
			//	// 3.2 新增 redis video 被点赞记录
			//	pipeline.HSet(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10), "0")
			//	// 3.3 redis 中 video 被点赞计数自增
			//	pipeline.Incr(l.ctx, videoLikedCntKey)
			//	// 3.4 redis 中 user 点赞计数自增
			//	pipeline.Incr(l.ctx, userLikeCntKey)
			//	// 3.5 pipeline 执行
			//	pipeline.Exec(l.ctx)
			//	return nil
			//})
			if e != nil && e != redis.Nil {
				// pipeline 操作失败
				l.Logger.Error(e)
				return &types.FavoriteActionResp{
					StatusCode: common.RedisError,
					StatusMsg:  common.MapErrMsg(common.RedisError),
				}, nil
			}
			fmt.Println("执行pipeline成功")
			// TODO 4. redis写数据库定时任务开启中
			fmt.Println("【api-favoriteAction-用户点赞成功】")
		}
	} else {
		// 5.取消点赞
		// 5.1 查询 redis 点赞记录
		likeRecord, err1 := redisClient.HgetCtx(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10))
		if err1 != nil && err1 != redis.Nil {
			l.Logger.Error(err1)
			return &types.FavoriteActionResp{
				StatusCode: common.RedisError,
				StatusMsg:  common.MapErrMsg(common.RedisError),
			}, nil
		}
		if len(likeRecord) != 0 && likeRecord == "1" {
			logx.Error("api-favoriteAction-已取消点赞，重复操作无效")
		} else {
			// 新建 redis 连接
			c := red.NewClient(&red.Options{
				Addr:     "127.0.0.1:8094",
				Password: common.DefaultPass,
			})
			//一起执行 pipeline 事务操作
			pipeline := c.TxPipeline()
			// 5.2 取消 redis 视频点赞用户记录
			pipeline.HSet(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10), "1")
			// 5.3 redis 中 video 被点赞计数自减
			pipeline.Decr(l.ctx, videoLikedCntKey)
			// 5.4 redis 中 user 点赞计数自减
			pipeline.Decr(l.ctx, userLikeCntKey)
			// 5.5 pipeline 执行
			_, e := pipeline.Exec(l.ctx)

			//// 一起执行 pipeline 操作
			//e := redisClient.PipelinedCtx(l.ctx, func(pipeline redis.Pipeliner) error {
			//	// 5.2 取消 redis 视频点赞用户记录
			//	pipeline.HSet(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10), "1")
			//	// 5.3 redis 中 video 被点赞计数自减
			//	pipeline.Decr(l.ctx, videoLikedCntKey)
			//	// 5.4 redis 中 user 点赞计数自减
			//	pipeline.Decr(l.ctx, userLikeCntKey)
			//	// 5.5 pipeline 执行
			//	pipeline.Exec(l.ctx)
			//	return nil
			//})
			if e != nil && e != redis.Nil {
				// pipeline 操作失败
				l.Logger.Error(e)
				return &types.FavoriteActionResp{
					StatusCode: common.RedisError,
					StatusMsg:  common.MapErrMsg(common.RedisError),
				}, nil
			}
			fmt.Println("执行pipeline成功")
			// TODO 6. redis写数据库定时任务执行中
			fmt.Println("【api-favoriteAction-用户取消点赞成功】")
		}
	}
	return &types.FavoriteActionResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
	}, nil

}
