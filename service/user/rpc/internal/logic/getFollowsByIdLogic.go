package logic

import (
	"context"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowsByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowsByIdLogic {
	return &GetFollowsByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowsByIdLogic) GetFollowsById(in *pb.GetFollowsByIdReq) (*pb.GetFollowsByIdResp, error) {
	// todo: add your logic here and delete this line
	// 数据库测试

	one, err := l.svcCtx.FollowsModel.FindOne(l.ctx, 1)
	println("数据库测试")
	println(one)
	//println(err.Error())
	println("数据库测试1111111111")
	println(one.Id)
	if err != nil {
		println("数据库测试失败")
		return nil, err
	}
	println("数据库测试")
	l.Logger.Info(one.Id)
	var pbFollows pb.Follows
	//pbFollows.Id = one.Id
	//println(pbFollows.Id)
	return &pb.GetFollowsByIdResp{
		Follows: &pbFollows,
	}, nil
}
