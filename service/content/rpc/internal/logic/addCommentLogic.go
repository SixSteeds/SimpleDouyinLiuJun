package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"errors"

	"time"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCommentLogic) AddComment(in *pb.AddCommentReq) (*pb.AddCommentResp, error) {
	/*
		Author：    刘洋
		Function：  向 comment 表添加评论
		Update：    08.28 对进入逻辑 加log
	*/
	l.Logger.Info("AddComment方法请求参数：", in)
	//1. 雪花算法生成 id
	snowflake, err1 := util.NewSnowflake(common.ContentRpcMachineId)
	if err1 != nil {
		return nil, errors.New("rpc-AddComment-新增评论，snowflake生成id失败")
	}
	snowId := snowflake.Generate()
	//2. 新增评论信息到 comment 表项
	_, err := l.svcCtx.CommentModel.Insert(l.ctx, &model.Comment{
		Id:         snowId,
		UserId:     in.UserId,
		VideoId:    in.VideoId,
		Content:    in.Content,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		IsDelete:   0,
	})
	if err != nil {
		return nil, errors.New("rpc-AddComment-新增评论数据失败")
	}

	return &pb.AddCommentResp{}, nil
}
