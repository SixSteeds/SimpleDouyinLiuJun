package comment

import (
	"context"
	"doushen_by_liujun/service/content/rpc/pb"
	"fmt"

	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	// todo: add your logic here and delete this line
	fmt.Println(req.VideoId, req.Token) //校验token
	follows, e := l.svcCtx.ContentRpcClient.GetCommentById(l.ctx, &pb.GetCommentByIdReq{
		Id: req.VideoId,
	})
	fmt.Println("查评论列表啦！！！！！！")
	fmt.Println(follows, e)
	//拿到了两条数据，还要查别人的表，之后写
	return
}
