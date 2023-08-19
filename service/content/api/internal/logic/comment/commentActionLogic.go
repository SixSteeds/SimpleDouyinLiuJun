package comment

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/rpc/pb"
	"log"

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

	//1.根据 token 获取 userid
	parsToken, err0 := util.ParseToken(req.Token)
	if err0 != nil {
		if err := l.svcCtx.KqPusherClient.Push("content_api_comment_commentActionLogic_CommentAction_ParseToken_false"); err != nil {
			log.Fatal(err)
		}
		// 返回token失效错误
		return &types.CommentActionResp{
			StatusCode: common.TOKEN_EXPIRE_ERROR,
			StatusMsg:  common.MapErrMsg(common.TOKEN_EXPIRE_ERROR),
		}, nil
	}
	if action := req.ActionType; action == 1 { // actionType（1评论，2删除评论）
		// 2.新增评论
		_, err1 := l.svcCtx.ContentRpcClient.AddComment(l.ctx, &pb.AddCommentReq{
			VideoId:  req.VideoId,
			UserId:   parsToken.UserID,
			Content:  req.CommentText,
			IsDelete: 0,
		})
		if err1 != nil {
			if err := l.svcCtx.KqPusherClient.Push("content_api_comment_commentActionLogic_AddComment_false"); err != nil {
				log.Fatal(err)
			}
			// 返回数据库查询错误
			return &types.CommentActionResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			}, nil
		}
		//fmt.Println(req.CommentText)
		logx.Error("api-commentAction-用户评论成功")
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
			}, nil
		}
		logx.Error("api-commentAction-用户删除评论成功")
	}
	if err := l.svcCtx.KqPusherClient.Push("content_api_comment_commentActionLogic_CommentAction_success"); err != nil {
		log.Fatal(err)
	}
	return &types.CommentActionResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
	}, nil

}
