package logic

import (
	"context"
	"fmt"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentByIdLogic {
	return &GetCommentByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentByIdLogic) GetCommentById(in *pb.GetCommentByIdReq) (*pb.GetCommentByIdResp, error) {
	comments, err := l.svcCtx.CommentModel.FindConmentsByVideoId(l.ctx, in.Id)
	//a, err := l.svcCtx.CommentForUserModel.FindOne(l.ctx, 3)
	//fmt.Println("dr数据库", a)
	if err != nil {
		return nil, err
	}
	var resp []*pb.Comment
	for _, item := range *comments {
		fmt.Println(item)
		resp = append(resp, &pb.Comment{
			Id:         item.Id,
			VideoId:    item.VideoId,
			Content:    item.Content,
			CreateTime: item.CreateTime.Unix(),
			UpdateTime: item.UpdateTime.Unix(),
		})
	}
	return &pb.GetCommentByIdResp{
		Comment: resp,
	}, nil
}
