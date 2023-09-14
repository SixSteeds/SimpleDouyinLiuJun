package userinfo

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type UserinfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserinfoLogic {
	return &UserinfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserinfoLogic) Userinfo(req *types.UserinfoReq) (resp *types.UserinfoResp, err error) {
	//fmt.Println(req)//for test
	//a := pb.GetUsersByIdsReq{
	//	UserID: 879260009437863936,
	//	Ids:    []int64{879651486764638208, 879004869187346432},
	//}
	//info, err := l.svcCtx.UserRpcClient.GetUsersByIds(l.ctx, &a)
	//fmt.Println("在api里", err)
	//fmt.Println(info)
	//return nil, nil
	l.Logger.Info(req)
	logger, e := util.ParseToken(req.Token)
	if e != nil {
		l.Logger.Error(e)
		return &types.UserinfoResp{
			StatusCode: common.TokenExpireError,
			StatusMsg:  common.MapErrMsg(common.TokenExpireError),
			User:       types.User{},
		}, nil
	}
	IntUserId, _ := strconv.Atoi(logger.ID)

	info, e := l.svcCtx.UserRpcClient.GetUserinfoById(l.ctx, &pb.GetUserinfoByIdReq{
		Id:     req.UserId,
		UserID: int64(IntUserId),
	})
	var user types.User
	if e != nil {
		l.Logger.Error(e)
		return &types.UserinfoResp{
			StatusCode: common.DbError,
			StatusMsg:  common.MapErrMsg(common.DbError),
			User:       user,
		}, nil
	}
	user = types.User{
		UserId:          info.Userinfo.Id,
		Name:            info.Userinfo.Username,
		FollowCount:     info.Userinfo.FollowCount,
		FollowerCount:   info.Userinfo.FollowerCount,
		IsFollow:        info.Userinfo.IsFollow, //我对这个的理解就是当前用户对这条数据的用户是否关注
		Avatar:          info.Userinfo.Avatar,
		BackgroundImage: info.Userinfo.BackgroundImage,
		Signature:       info.Userinfo.Signature,
		WorkCount:       info.Userinfo.WorkCount,
		FavoriteCount:   info.Userinfo.FavoriteCount,
		TotalFavorited:  info.Userinfo.TotalFavorited, //查表
	}
	return &types.UserinfoResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		User:       user,
	}, nil
}
