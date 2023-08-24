package logic

import (
	"context"
	"fmt"

	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkCountByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkCountByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkCountByUserIdLogic {
	return &GetWorkCountByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkCountByUserIdLogic) GetWorkCountByUserId(in *pb.GetWorkCountByUserIdReq) (*pb.GetWorkCountByUserIdResp, error) {
	// todo: add your logic here and delete this line
	fmt.Println("我要去model了，GetWorkCountByUserId")
	count, err := l.svcCtx.VideoModel.GetWorkCountByUserId(l.ctx, in.UserId)
	fmt.Println("我从model回来了，GetWorkCountByUserId")
	fmt.Println(count.Count)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return &pb.GetWorkCountByUserIdResp{
		WorkCount: count.Count,
	}, nil
}
