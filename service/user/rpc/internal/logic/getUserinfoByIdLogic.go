package logic

import (
	"context"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserinfoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserinfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserinfoByIdLogic {
	return &GetUserinfoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserinfoByIdLogic) GetUserinfoById(in *pb.GetUserinfoByIdReq) (*pb.GetUserinfoByIdResp, error) {
	// todo: add your logic here and delete this line
	info, err := l.svcCtx.UserinfoModel.FindOne(l.ctx, in.Id, in.UserID)
	if err != nil {
		return nil, err
	}
	userInfo := pb.Userinfo{
		Id:              info.Id,
		FollowCount:     info.FollowCount,
		FollowerCount:   info.FollowerCount,
		IsFollow:        info.IsFollow,
		Username:        info.Username.String,
		Avatar:          info.Avatar.String,
		BackgroundImage: info.BackgroundImage.String,
		Signature:       info.Signature.String,
	}
	return &pb.GetUserinfoByIdResp{
		Userinfo: &userInfo,
	}, nil
}
