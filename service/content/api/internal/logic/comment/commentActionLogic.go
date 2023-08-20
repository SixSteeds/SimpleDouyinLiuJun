package comment

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/content/rpc/pb"

	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"log"
	"strconv"

	constants "doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func executeCntRedis(l *CommentActionLogic, redisKey string, incFlag bool) (err error) {
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
			redisClient.SetCtx(l.ctx, redisKey, "1")
		} else {
			// 有记录，cnt 自增1
			redisClient.IncrCtx(l.ctx, redisKey)
		}
	} else { //自减计数
		info, e := redisClient.GetCtx(l.ctx, redisKey)
		if e != nil && e != redis.Nil { //查询redis报错
			return e
		}
		if len(info) == 0 { // 没有记录无法再减少，返回错误
			return redis.Nil
		} else { // 有记录，cnt 自减1
			redisClient.DecrCtx(l.ctx, redisKey)
		}
	}
	fmt.Println("【executeCntRedis-执行成功】")
	return nil
}

func (l *CommentActionLogic) CommentAction(req *types.CommentActionReq) (resp *types.CommentActionResp, err error) {

	//1.根据 token 获取 userid
	//parsToken, err0 := util.ParseToken(req.Token)
	//if err0 != nil {
	//	// 返回token失效错误
	//	return &types.CommentActionResp{
	//		StatusCode: common.TOKEN_EXPIRE_ERROR,
	//		StatusMsg:  common.MapErrMsg(common.TOKEN_EXPIRE_ERROR),
	//	}, nil
	//}
	var test_useid int64 = 7
	videoCommentedCntKey := constants.CntCacheVideoCommentedPrefix + strconv.FormatInt(req.VideoId, 10)

	if action := req.ActionType; action == 1 { // actionType（1评论，2删除评论）
		// 2.新增评论
		_, err1 := l.svcCtx.ContentRpcClient.AddComment(l.ctx, &pb.AddCommentReq{
			VideoId:  req.VideoId,
			UserId:   test_useid, //parsToken.UserID,
			Content:  req.CommentText,
			IsDelete: 0,
		})
		if err1 != nil {
			if err := l.svcCtx.KqPusherClient.Push("content_api_comment_commentActionLogic_AddComment_false"); err != nil {
				log.Fatal(err)
			}
			// 返回数据库查询错误
			return &types.CommentActionResp{
				StatusCode: common.REDIS_ERROR,
				StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			}, err1
		}
		// 2.1 redis 中 video 被评论计数自增
		err2 := executeCntRedis(l, videoCommentedCntKey, true)
		if err2 != nil {
			// 返回 redis 访问错误
			return &types.CommentActionResp{
				StatusCode: common.REDIS_ERROR,
				StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			}, err2
		}
		fmt.Println("【api-commentAction-用户评论成功】")
	} else {
		//3.删除评论
		_, err1 := l.svcCtx.ContentRpcClient.DelComment(l.ctx, &pb.DelCommentReq{
			Id: req.CommentId,
		})
		if err1 != nil {
			if err := l.svcCtx.KqPusherClient.Push("content_api_comment_commentActionLogic_DelComment_false"); err != nil {
				log.Fatal(err)
			}
			// 返回数据库查询错误
			return &types.CommentActionResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			}, err1
		}
		// 3.1 redis 中 video 被评论计数自减
		err2 := executeCntRedis(l, videoCommentedCntKey, false)
		if err2 != nil {
			// 返回 redis 访问错误
			return &types.CommentActionResp{
				StatusCode: common.REDIS_ERROR,
				StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			}, err2
		}
		fmt.Println("【api-commentAction-用户删除评论成功】")
	}
	if err := l.svcCtx.KqPusherClient.Push("content_api_comment_commentActionLogic_CommentAction_success"); err != nil {
		log.Fatal(err)
	}
	return &types.CommentActionResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
	}, nil

}
