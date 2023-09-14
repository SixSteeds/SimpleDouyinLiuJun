package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCommentLogic {
	return &DelCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelCommentLogic) DelComment(in *pb.DelCommentReq) (*pb.DelCommentResp, error) {
	/*
		Author：    刘洋
		Function：  从 comment 表删除评论信息
		Update：    08.28 对进入逻辑 加log
	*/
	l.Logger.Info("DelComment方法请求参数：", in)
	//PS.删除评论不是高频操作，所以不逻辑删除而是直接查库删
	err := l.svcCtx.CommentModel.Delete(l.ctx, in.Id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.New("rpc-DelComment-删除评论数据失败")
	}
	fmt.Println("【rpc-DelComment-删除评论数据成功】")
	return &pb.DelCommentResp{}, nil
}
