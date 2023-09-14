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

type GetUserListByIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListByIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListByIdListLogic {
	return &GetUserListByIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListByIdListLogic) GetUserListByIdList(in *pb.GetUserListByIdListReq) (*pb.GetUserListByIdListResp, error) {
	videoList, err := l.svcCtx.UserinfoModel.FindUserListByIdList(l.ctx, &in.UserIdList)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.New("rpc-GetUserListByIdList-数据查询失败")
	}
	var resp []*pb.Userinfo
	for _, item := range *videoList {
		resp = append(resp, &pb.Userinfo{
			Id: item.Id,
			//FollowCount: ,
			//FollowerCount: ,
			//IsFollow: ,
			Username:        item.Username.String,
			Avatar:          item.Avatar.String,
			BackgroundImage: item.BackgroundImage.String,
			Signature:       item.Signature.String,
		})
	}
	fmt.Println("【rpc-GetUserListByIdList-根据用户id列表查询列表用户成功】")
	return &pb.GetUserListByIdListResp{
		UserList: resp,
	}, nil

}
