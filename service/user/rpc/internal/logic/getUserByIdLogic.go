package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"errors"
	"fmt"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *pb.GetUserByIdReq) (*pb.GetUserByIdResp, error) {
	userInfo, err := l.svcCtx.UserinfoModel.FindUserById(l.ctx, in.UserID)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.New("数据查询失败")
	}
	if userInfo != nil {
		fmt.Println("查到")
		return &pb.GetUserByIdResp{
			Userinfo: &pb.Userinfo{
				Id: userInfo.Id,
				//FollowCount: ,//用户关注数
				//FollowerCount: ,//用户被关注数
				//IsFollow: ,//当前登录用户对这个用户（视频作者）是否关注
				Username:        userInfo.Username.String,
				Avatar:          userInfo.Avatar.String,
				BackgroundImage: userInfo.BackgroundImage.String,
				Signature:       userInfo.Signature.String,
			},
		}, nil
	}
	fmt.Println("没查到")
	return &pb.GetUserByIdResp{}, nil
}
