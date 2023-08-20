package favorite

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"

	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"

	constants "doushen_by_liujun/internal/common"
	"github.com/zeromicro/go-zero/core/logx"
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

//func executeCntRedis(l *FavoriteActionLogic, redisKey string, incFlag bool) (err error) {
//	// incFlag=true：  redis计数自增，没有记录则新建 redis 中计数并设置初始值为1
//	// incFlag=false： redis计数自减，没有则返回 redis 查询错误
//	redisClient := l.svcCtx.RedisClient
//	if incFlag == true { //自增计数
//		info, e := redisClient.GetCtx(l.ctx, redisKey)
//		if e != nil && e != redis.Nil { //查询redis报错
//			return e
//		}
//		if len(info) == 0 {
//			// 没有记录，新增记录并令 cnt=1
//			redisClient.SetCtx(l.ctx, redisKey, "1")
//		} else {
//			// 有记录，cnt 自增1
//			redisClient.IncrCtx(l.ctx, redisKey)
//		}
//	} else { //自减计数
//		info, e := redisClient.GetCtx(l.ctx, redisKey)
//		if e != nil && e != redis.Nil { //查询redis报错
//			return e
//		}
//		if len(info) == 0 { // 没有记录无法再减少，返回错误
//			return redis.Nil
//		} else { // 有记录，cnt 自减1
//			redisClient.DecrCtx(l.ctx, redisKey)
//		}
//	}
//	fmt.Println("executeCntRedis-执行成功")
//	return nil
//}

func executeCntRedis(l *FavoriteActionLogic, redisKey string, pipeline redis.Pipeliner, incFlag bool) (err error) {
	// incFlag=true：  redis计数自增，没有记录则新建 redis 中计数并设置初始值为1
	// incFlag=false： redis计数自减，没有则返回 redis 查询错误
	redisClient := l.svcCtx.RedisClient
	if incFlag == true { //自增计数
		info, e := redisClient.GetCtx(l.ctx, redisKey)
		if e != nil && e != redis.Nil { //查询redis报错
			return e
		}
		if len(info) == 0 {
			// 没有记录，新增记录并令 cnt=1
			pipeline.Set(l.ctx, redisKey, "1", time.Duration(-1))
		} else {
			// 有记录，cnt 自增1
			pipeline.Incr(l.ctx, redisKey)
		}
	} else { //自减计数
		info, e := redisClient.GetCtx(l.ctx, redisKey)
		if e != nil && e != redis.Nil { //查询redis报错
			return e
		}
		if len(info) == 0 { // 没有记录无法再减少，返回错误
			return redis.Nil
		} else { // 有记录，cnt 自减1
			pipeline.Decr(l.ctx, redisKey)
		}
	}
	fmt.Println("executeCntRedis-执行成功")
	return nil
}

func (l *FavoriteActionLogic) FavoriteAction(req *types.FavoriteActionReq) (resp *types.FavoriteActionResp, err error) {

	//1.根据 token 获取 userid
	parsToken, err0 := util.ParseToken(req.Token)
	if err0 != nil {
		// 返回token失效错误
		return &types.FavoriteActionResp{
			StatusCode: common.TOKEN_EXPIRE_ERROR,
			StatusMsg:  common.MapErrMsg(common.TOKEN_EXPIRE_ERROR),
		}, nil
	}
	//var test_useid int64 = 8

	// TODO 2.加入redis缓存
	redisClient := l.svcCtx.RedisClient
	videoLikedKey := constants.LikeCacheVideoLikedPrefix + strconv.FormatInt(req.VideoId, 10)
	videoLikedCntKey := constants.CntCacheVideoLikedPrefix + strconv.FormatInt(req.VideoId, 10)
	userLikeCntKey := constants.CntCacheUserLikePrefix + strconv.FormatInt(parsToken.UserID, 10)

	if action := req.ActionType; action == 1 { // actionType（1点赞，2取消）
		// 2.新增点赞
		// 2.1 查询 redis 点赞记录
		likeRecord, err1 := redisClient.HgetCtx(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10))
		if err1 != nil && err1 != redis.Nil {
			// 返回 redis 访问错误
			return &types.FavoriteActionResp{
				StatusCode: common.REDIS_ERROR,
				StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			}, err1
		}
		if len(likeRecord) != 0 && likeRecord == "0" {
			logx.Error("api-favoriteAction-已点赞，重复操作无效")
		} else {
			// 一起执行 pipeline 操作
			e0 := redisClient.PipelinedCtx(l.ctx, func(pipeline redis.Pipeliner) error {
				// 2.2 新增 redis video 被点赞记录
				pipeline.HSet(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10), "0")
				// 2.3 redis 中 video 被点赞计数自增
				pipeline.Incr(l.ctx, videoLikedCntKey)
				// 2.4 redis 中 user 点赞计数自增
				pipeline.Incr(l.ctx, userLikeCntKey)
				// 2.5 pipeline 执行
				pipeline.Exec(l.ctx)
				return nil
			})
			if e0 != nil && e0 != redis.Nil {
				// pipeline 操作失败
				return &types.FavoriteActionResp{
					StatusCode: common.REDIS_ERROR,
					StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
				}, e0
			}
			fmt.Println("执行pipeline成功")
			//// 2.2 新增 redis video 被点赞记录
			//redisClient.HsetCtx(l.ctx, videoLikedKey, strconv.FormatInt(test_useid, 10), "1")
			//// 2.3 redis 中 video 被点赞计数自增
			//err2 := executeCntRedis(l, videoLikedCntKey, true)
			//if err2 != nil {
			//	// 返回 redis 访问错误
			//	return &types.FavoriteActionResp{
			//		StatusCode: common.DB_ERROR,
			//		StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			//	}, err2
			//}
			//// 2.4 redis 中 user 点赞计数自增
			//err3 := executeCntRedis(l, userLikeCntKey, true)
			//if err3 != nil {
			//	// 返回 redis 访问错误
			//	return &types.FavoriteActionResp{
			//		StatusCode: common.DB_ERROR,
			//		StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			//	}, err3
			//}

			// TODO 3. 新增点赞数据库修改->新增redis写数据库定时任务
			//_, err4 := l.svcCtx.ContentRpcClient.AddFavorite(l.ctx, &pb.AddFavoriteReq{
			//	UserId:   test_useid, //parsToken.UserID,
			//	VideoId:  req.VideoId,
			//	IsDelete: 0,
			//})
			//if err4 != nil {
			//	// 返回数据库新增错误
			//	return &types.FavoriteActionResp{
			//		StatusCode: common.DB_ERROR,
			//		StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			//	}, err4
			//}
			fmt.Println("【api-favoriteAction-用户点赞成功】")
		}
	} else {
		// 4.取消点赞
		// 4.1 查询 redis 点赞记录
		likeRecord, err1 := redisClient.HgetCtx(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10))
		if err1 != nil && err1 != redis.Nil {
			// 返回 redis 访问错误
			return &types.FavoriteActionResp{
				StatusCode: common.REDIS_ERROR,
				StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			}, err1
		}
		if len(likeRecord) != 0 && likeRecord == "1" {
			logx.Error("api-favoriteAction-已取消点赞，重复操作无效")
		} else {
			// 一起执行 pipeline 操作
			e0 := redisClient.PipelinedCtx(l.ctx, func(pipeline redis.Pipeliner) error {
				// 4.2 取消 redis 视频点赞用户记录
				pipeline.HSet(l.ctx, videoLikedKey, strconv.FormatInt(parsToken.UserID, 10), "1")
				// 4.3 redis 中 video 被点赞计数自减
				pipeline.Decr(l.ctx, videoLikedCntKey)
				// 4.4 redis 中 user 点赞计数自减
				pipeline.Decr(l.ctx, userLikeCntKey)
				// 2.5 pipeline 执行
				pipeline.Exec(l.ctx)
				return nil
			})
			if e0 != nil && e0 != redis.Nil {
				// 事务执行失败，已回滚
				return &types.FavoriteActionResp{
					StatusCode: common.REDIS_ERROR,
					StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
				}, e0
			}
			fmt.Println("执行pipeline成功")
			//// 4.2 取消 redis 视频点赞用户记录
			//redisClient.HsetCtx(l.ctx, videoLikedKey, strconv.FormatInt(test_useid, 10), "0")
			//// 4.3 redis 中 video 被点赞计数自减
			//err2 := executeCntRedis(l, videoLikedCntKey, false)
			//if err2 != nil {
			//	// 返回 redis 访问错误
			//	return &types.FavoriteActionResp{
			//		StatusCode: common.DB_ERROR,
			//		StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			//	}, err2
			//}
			//// 4.4 redis 中 user 点赞计数自减
			//err3 := executeCntRedis(l, userLikeCntKey, false)
			//if err3 != nil {
			//	// 返回 redis 访问错误
			//	return &types.FavoriteActionResp{
			//		StatusCode: common.DB_ERROR,
			//		StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			//	}, err3
			//}

			// TODO 5. 取消点赞数据库修改->新增redis写数据库定时任务
			//_, err4 := l.svcCtx.ContentRpcClient.DelFavorite(l.ctx, &pb.DelFavoriteReq{
			//	UserId:  test_useid, //parsToken.UserID,
			//	VideoId: req.VideoId,
			//})
			//if err4 != nil {
			//	// 返回数据库删除错误
			//	return &types.FavoriteActionResp{
			//		StatusCode: common.DB_ERROR,
			//		StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			//	}, err4
			//}
			fmt.Println("【api-favoriteAction-用户取消点赞成功】")
		}
	}
	return &types.FavoriteActionResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
	}, nil

}
