package favorite

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/content/rpc/pb"

	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteActionLogic) FavoriteAction(req *types.FavoriteActionReq) (resp *types.FavoriteActionResp, err error) {

	//1.根据 token 获取 userid
	parsToken, err0 := util.ParseToken(req.Token)
	if err0 != nil {
		// 返回token失效错误
		return &types.FavoriteActionResp{
			StatusCode: common.TOKEN_EXPIRE_ERROR,
			StatusMsg:  common.MapErrMsg(common.TOKEN_EXPIRE_ERROR),
		}, nil
	}

	if action := req.ActionType; action == 1 { // actionType（1点赞，2取消）
		// 2.新增点赞
		_, err1 := l.svcCtx.ContentRpcClient.AddFavorite(l.ctx, &pb.AddFavoriteReq{
			UserId:   parsToken.UserID,
			VideoId:  req.VideoId,
			IsDelete: 0,
		})
		if err1 != nil {
			// 返回数据库查询错误
			return &types.FavoriteActionResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			}, nil
		}
		logx.Error("api-favoriteAction-用户点赞成功")
	} else {
		//3.取消点赞
		_, err1 := l.svcCtx.ContentRpcClient.DelFavorite(l.ctx, &pb.DelFavoriteReq{
			UserId:  parsToken.UserID,
			VideoId: req.VideoId,
		})
		if err1 != nil {
			// 返回数据库查询错误
			return &types.FavoriteActionResp{
				StatusCode: common.DB_ERROR,
				StatusMsg:  common.MapErrMsg(common.DB_ERROR),
			}, nil
		}
		logx.Error("api-favoriteAction-用户取消点赞成功")
	}

	return &types.FavoriteActionResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
	}, nil

}
