package logic

import (
	"context"
	"fmt"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsByIdLogic {
	return &GetFriendsByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendsByIdLogic) GetFriendsById(in *pb.GetFriendsByIdReq) (*pb.GetFriendsByIdResp, error) {
	// todo: add your logic here and delete this line
	follows, err := l.svcCtx.FollowsModel.FindFriendsByUserId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	var resp []*pb.Follows
	for _, item := range *follows {
		fmt.Println(item)
		resp = append(resp, &pb.Follows{
			UserId:   item.UserId,
			FollowId: item.FollowId,
		})
	}
	return &pb.GetFriendsByIdResp{
		Follows: resp,
	}, nil
}
