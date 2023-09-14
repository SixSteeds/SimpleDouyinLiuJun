package comment

import (
	"context"
	"database/sql"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/rpc/pb"
	"errors"

	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
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

func (l *CommentActionLogic) CommentAction(req *types.CommentActionReq) (resp *types.CommentActionResp, err error) {
	/*
		Author：    刘洋
		Function：  评论、删除评论（ ActionType=1 评论，ActionType=2 删除 ）
		Update：    08.28 对进入逻辑、异常 加log
	*/
	l.Logger.Info(req)
	redisClient := l.svcCtx.RedisClient
	videoCommentedCntKey := constants.CntCacheVideoCommentedPrefix + strconv.FormatInt(req.VideoId, 10)

	// 1.根据 token 获取 userid
	parsToken, err0 := util.ParseToken(req.Token)
	if err0 != nil {
		l.Logger.Error(err0)
		return &types.CommentActionResp{
			StatusCode: common.TokenExpireError,
			StatusMsg:  common.MapErrMsg(common.TokenExpireError),
		}, nil
	}
	if action := req.ActionType; action == 1 { // actionType（1新增评论，2删除评论）
		// 2.新增评论
		_, err1 := l.svcCtx.ContentRpcClient.AddComment(l.ctx, &pb.AddCommentReq{
			VideoId:  req.VideoId,
			UserId:   parsToken.UserID,
			Content:  req.CommentText,
			IsDelete: 0,
		})
		if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
			l.Logger.Error(err1)
			return &types.CommentActionResp{
				StatusCode: common.DbError,
				StatusMsg:  common.MapErrMsg(common.DbError),
			}, nil
		}
		// 2.1 redis 中 video 被评论计数自增
		_, err2 := redisClient.IncrCtx(l.ctx, videoCommentedCntKey)
		if err2 != nil && err2 != redis.Nil {
			l.Logger.Error(err2)
			return &types.CommentActionResp{
				StatusCode: common.RedisError,
				StatusMsg:  common.MapErrMsg(common.RedisError),
			}, nil
		}
		fmt.Println("【api-commentAction-用户评论成功】")
	} else {
		// 3.删除评论
		_, err1 := l.svcCtx.ContentRpcClient.DelComment(l.ctx, &pb.DelCommentReq{
			Id: req.CommentId,
		})
		if err1 != nil {
			l.Logger.Error(err1)
			return &types.CommentActionResp{
				StatusCode: common.DbError,
				StatusMsg:  common.MapErrMsg(common.DbError),
			}, nil
		}
		// 3.1 redis 中 video 被评论计数自减
		_, err2 := redisClient.DecrCtx(l.ctx, videoCommentedCntKey)
		if err2 != nil && err2 != redis.Nil {
			l.Logger.Error(err2)
			return &types.CommentActionResp{
				StatusCode: common.RedisError,
				StatusMsg:  common.MapErrMsg(common.RedisError),
			}, nil
		}
		fmt.Println("【api-commentAction-用户删除评论成功】")
	}
	return &types.CommentActionResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
	}, nil

}
