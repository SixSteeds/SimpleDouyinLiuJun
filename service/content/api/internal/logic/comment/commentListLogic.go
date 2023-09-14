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
	l.Logger.Info(req)
	logger, e := util.ParseToken(req.Token)
	if e != nil {
		l.Logger.Error(e)
		return &types.CommentListResp{
			StatusCode:  common.TokenExpireError,
			StatusMsg:   common.MapErrMsg(common.TokenExpireError),
			CommentList: []types.Comment{},
		}, nil
	}
	comment, e := l.svcCtx.ContentRpcClient.GetCommentById(l.ctx, &pb.GetCommentByIdReq{
		Id: req.VideoId,
	})
	var comments []types.Comment
	if e != nil {
		l.Logger.Error(e)
		return &types.CommentListResp{
			StatusCode:  common.DbError,
			StatusMsg:   common.MapErrMsg(common.DbError),
			CommentList: []types.Comment{},
		}, nil
	}
	if len(comment.Comment) == 0 {
		return &types.CommentListResp{
			StatusCode:  common.OK,
			StatusMsg:   common.MapErrMsg(common.OK),
			CommentList: comments,
		}, nil
	}
	IntUserId, _ := strconv.Atoi(logger.ID)
	//IntUserId := 223
	var ids []int64
	for _, item := range comment.Comment {
		ids = append(ids, item.UserId)
	}
	info, e := l.svcCtx.UserRpcClient.GetUsersByIds(l.ctx, &userPB.GetUsersByIdsReq{
		Ids:    ids,
		UserID: int64(IntUserId),
	})
	if e != nil {
		l.Logger.Error(e)
		return &types.CommentListResp{
			StatusCode:  common.DbError,
			StatusMsg:   common.MapErrMsg(common.DbError),
			CommentList: []types.Comment{},
		}, nil
	}
	users := info.Users
	for index, item := range comment.Comment {
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
