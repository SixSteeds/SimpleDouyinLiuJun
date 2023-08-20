package comment

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"
	"doushen_by_liujun/service/content/rpc/pb"
	userPB "doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"strconv"
	"time"
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
	logger, e := util.ParseToken(req.Token)
	if e != nil {
		if err := l.svcCtx.KqPusherClient.Push("content_api_comment_CommentListLogic_CommentList_ParseToken_false"); err != nil {
			log.Fatal(err)
		}
		return &types.CommentListResp{
			StatusCode:  common.TOKEN_EXPIRE_ERROR,
			StatusMsg:   common.MapErrMsg(common.TOKEN_EXPIRE_ERROR),
			CommentList: []types.Comment{},
		}, e
	}
	follows, e := l.svcCtx.ContentRpcClient.GetCommentById(l.ctx, &pb.GetCommentByIdReq{
		Id: req.VideoId,
	})
	var comments []types.Comment
	if e != nil {
		if err := l.svcCtx.KqPusherClient.Push("content_api_comment_CommentListLogic_CommentList_GetCommentById_false"); err != nil {
			log.Fatal(err)
		}
		return &types.CommentListResp{
			StatusCode:  common.DB_ERROR,
			StatusMsg:   common.MapErrMsg(common.DB_ERROR),
			CommentList: []types.Comment{},
		}, e
	}
	IntUserId, _ := strconv.Atoi(logger.ID)
	//IntUserId := 223
	var ids []int64
	for _, item := range follows.Comment {
		ids = append(ids, item.UserId)
	}
	info, e := l.svcCtx.UserRpcClient.GetUsersByIds(l.ctx, &userPB.GetUsersByIdsReq{
		Ids:    ids,
		UserID: int64(IntUserId),
	})
	if e != nil {
		return &types.CommentListResp{
			StatusCode:  common.DB_ERROR,
			StatusMsg:   common.MapErrMsg(common.DB_ERROR),
			CommentList: []types.Comment{},
		}, e
	}
	users := info.Users
	for index, item := range follows.Comment {
		var user types.User
		user = types.User{
			Id:              users[index].Id,
			Name:            users[index].Username,
			FollowCount:     users[index].FollowCount,
			FollowerCount:   users[index].FollowerCount,
			IsFollow:        users[index].IsFollow, //我对这个的理解就是当前用户对这条数据的用户是否关注
			Avatar:          users[index].Avatar,
			BackgroundImage: users[index].BackgroundImage,
			Signature:       users[index].Signature,
			WorkCount:       users[index].WorkCount,
			FavoriteCount:   users[index].FavoriteCount,
			TotalFavorited:  users[index].TotalFavorited,
		}
		comments = append(comments, types.Comment{
			Id:         item.Id,
			User:       user,
			Content:    item.Content,
			CreateDate: time.Unix(item.CreateTime, 0).Format("01-02"),
		})
	}
	return &types.CommentListResp{
		StatusCode:  common.OK,
		StatusMsg:   common.MapErrMsg(common.OK),
		CommentList: comments,
	}, nil
}
