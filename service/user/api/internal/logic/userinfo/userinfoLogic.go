package userinfo

import (
	"context"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line
	fmt.Println("用户信息") //校验token
	info, e := l.svcCtx.UserRpcClient.GetUserinfoById(l.ctx, &pb.GetUserinfoByIdReq{
		Id: req.UserId,
	})
	var user types.User
	if e != nil {
		return &types.UserinfoResp{
			StatusCode: -1,
			StatusMsg:  e.Error(),
			User:       user,
		}, err
	}
	followCount, e := l.svcCtx.UserRpcClient.GetFollowsCountById(l.ctx, &pb.GetFollowsCountByIdReq{
		Id: req.UserId,
	})
	if e != nil {
		return &types.UserinfoResp{
			StatusCode: -1,
			StatusMsg:  "查询关注数量失败",
			User:       user,
		}, err
	}
	followerCount, e := l.svcCtx.UserRpcClient.GetFollowersCountById(l.ctx, &pb.GetFollowersCountByIdReq{
		Id: req.UserId,
	})
	if e != nil {
		return &types.UserinfoResp{
			StatusCode: -1,
			StatusMsg:  "查询粉丝数量失败",
			User:       user,
		}, err
	}
	user = types.User{
		UserId:          info.Userinfo.Id,
		Name:            info.Userinfo.Name,
		FollowCount:     followCount.Count,
		FollowerCount:   followerCount.Count,
		IsFollow:        false, //查表
		Avatar:          info.Userinfo.Avatar,
		BackgroundImage: info.Userinfo.BackgroundImage,
		Signature:       info.Userinfo.Signature,
		WorkCount:       0, //查表
		FavoriteCount:   0, //查表
	}
	return &types.UserinfoResp{
		StatusCode: 0,
		StatusMsg:  "查询成功",
		User:       user,
	}, nil
}
