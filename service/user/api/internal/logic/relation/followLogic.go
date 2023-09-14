package relation

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/juju/ratelimit"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"strconv"
	"time"
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	bucket *ratelimit.Bucket
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		bucket: ratelimit.NewBucketWithRate(10, 10), //每秒钟生成 10 个令牌——令牌桶限流1s内最多处理10个请求
	}
}

func (l *FollowLogic) Follow(req *types.FollowReq) (resp *types.FollowResp, err error) {
	l.Logger.Info(req)
	logger, e := util.ParseToken(req.Token)
	if e != nil {
		l.Logger.Error(e)
		return &types.FollowResp{
			StatusCode: common.TokenExpireError,
			StatusMsg:  common.MapErrMsg(common.TokenExpireError),
		}, nil
	}
	//修正redis中数据，判断是关注还是取消关注
	redisClient := l.svcCtx.RedisClient
	followKey := common.FollowNum + strconv.Itoa(int(logger.UserID))
	followerKey := common.FollowerNum + strconv.Itoa(int(req.ToUserId))
	followRecord, _ := redisClient.GetCtx(l.ctx, followKey)
	followerRecord, _ := redisClient.GetCtx(l.ctx, followerKey)
	//rand.Seed(time.Now().UnixNano())
	expiration := 3000 + rand.Intn(600)
	if req.ActionType == 1 {
		if len(followRecord) != 0 { //有记录更新一下计数，没有就算了,还是等查用户信息时再查表
			followNum, _ := strconv.Atoi(followRecord)
			followNum += 1
			_ = redisClient.SetexCtx(l.ctx, followKey, strconv.Itoa(followNum), expiration)
		}
		if len(followerRecord) != 0 {
			followerNum, _ := strconv.Atoi(followerRecord)
			followerNum += 1
			_ = redisClient.SetexCtx(l.ctx, followerKey, strconv.Itoa(followerNum), expiration)
		}
	} else {
		if len(followRecord) != 0 { //有记录更新一下计数，没有就算了,还是等查用户信息时再查表
			followNum, _ := strconv.Atoi(followRecord)
			followNum -= 1
			_ = redisClient.SetexCtx(l.ctx, followKey, strconv.Itoa(followNum), expiration)
		}
		if len(followerRecord) != 0 {
			followerNum, _ := strconv.Atoi(followerRecord)
			followerNum -= 1
			_ = redisClient.SetexCtx(l.ctx, followerKey, strconv.Itoa(followerNum), expiration)
		}
	}
	//写入数据库
	if l.bucket.TakeAvailable(1) == 0 {
		// 令牌不足，限流处理
		randomInterval := time.Duration(rand.Int63n(int64(10 * time.Minute))) //10分钟内的随机时间
		//判断是关注还是取消关注
		if req.ActionType == 1 {
			// 使用time.AfterFunc延迟执行写入数据库的操作
			time.AfterFunc(randomInterval, func() {
				userId := strconv.Itoa(int(logger.UserID))
				followId := strconv.FormatInt(req.ToUserId, 10)
				_, err := l.svcCtx.UserRpcClient.AddFollows(l.ctx, &pb.AddFollowsReq{
					UserId:   userId,
					FollowId: followId,
				})
				if err != nil { //没写进去则隔一段时间重试
					randomInterval := time.Duration(rand.Int63n(int64(20 * time.Second))) //20s内的随机时间
					ticker := time.NewTicker(randomInterval)                              //延迟执行写入数据库
				OuterLoop:
					for {
						select {
						case _ = <-ticker.C:
							_, err := l.svcCtx.UserRpcClient.AddFollows(l.ctx, &pb.AddFollowsReq{
								UserId:   userId,
								FollowId: followId,
							})
							if err == nil { //确实写进去了再结束
								ticker.Stop()
								break OuterLoop
							}
						}
					}
				}
			})
			return &types.FollowResp{
				StatusCode: common.OK,
				StatusMsg:  common.MapErrMsg(common.OK),
			}, nil
		} else {
			// 使用time.AfterFunc延迟执行写入数据库的操作
			time.AfterFunc(randomInterval, func() {
				userId := strconv.Itoa(int(logger.UserID))
				followId := strconv.FormatInt(req.ToUserId, 10)
				_, err := l.svcCtx.UserRpcClient.DelFollows(l.ctx, &pb.DelFollowsReq{
					UserId:   userId,
					FollowId: followId,
				})
				if err != nil { //没写进去则隔一段时间重试
					randomInterval := time.Duration(rand.Int63n(int64(20 * time.Second))) //20s内的随机时间
					ticker := time.NewTicker(randomInterval)                              //延迟执行写入数据库
				OuterLoop:
					for {
						select {
						case _ = <-ticker.C:
							_, err := l.svcCtx.UserRpcClient.DelFollows(l.ctx, &pb.DelFollowsReq{
								UserId:   userId,
								FollowId: followId,
							})
							if err == nil { //确实写进去了再结束
								ticker.Stop()
								break OuterLoop
							}
						}
					}
				}
			})
			return &types.FollowResp{
				StatusCode: common.OK,
				StatusMsg:  common.MapErrMsg(common.OK),
			}, nil
		}
	} else {
		//判断是关注还是取消关注
		if req.ActionType == 1 {
			_, err := l.svcCtx.UserRpcClient.AddFollows(l.ctx, &pb.AddFollowsReq{
				UserId:   strconv.Itoa(int(logger.UserID)),
				FollowId: strconv.FormatInt(req.ToUserId, 10),
			})
			if err != nil {
				l.Logger.Error(err)
				return &types.FollowResp{
					StatusCode: common.DbError,
					StatusMsg:  common.MapErrMsg(common.DbError),
				}, nil
			}
			return &types.FollowResp{
				StatusCode: common.OK,
				StatusMsg:  common.MapErrMsg(common.OK),
			}, nil
		} else {
			_, err := l.svcCtx.UserRpcClient.DelFollows(l.ctx, &pb.DelFollowsReq{
				UserId:   strconv.Itoa(int(logger.UserID)),
				FollowId: strconv.FormatInt(req.ToUserId, 10),
			})
			if err != nil {
				l.Logger.Error(err)
				return &types.FollowResp{
					StatusCode: common.DbError,
					StatusMsg:  common.MapErrMsg(common.DbError),
				}, nil
			}
			return &types.FollowResp{
				StatusCode: common.OK,
				StatusMsg:  common.MapErrMsg(common.OK),
			}, nil
		}
	}
}
