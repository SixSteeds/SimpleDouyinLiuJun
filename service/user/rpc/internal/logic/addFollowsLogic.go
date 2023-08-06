package logic

import (
	"context"
	"database/sql"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"doushen_by_liujun/service/user/rpc/pb"
	"fmt"
	"time"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddFollowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFollowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFollowsLogic {
	return &AddFollowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------鐢ㄦ埛鍩烘湰淇℃伅-----------------------
func (l *AddFollowsLogic) AddFollows(in *pb.AddFollowsReq) (*pb.AddFollowsResp, error) {
	// todo: add your logic here and delete this line
	fmt.Println(in)
	// 将普通的string类型转换为sql.NullString类型
	nSUserId := sql.NullString{
		String: in.UserId,
		Valid:  true,
	}
	nSFollowId := sql.NullString{
		String: in.FollowId,
		Valid:  true,
	}
	res, err := l.svcCtx.FollowsModel.Insert(l.ctx, &model.Follows{
		UserId:     nSUserId,
		FollowId:   nSFollowId,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
		IsDelete:   0,
	})
	if err == nil {
		fmt.Println("res", res)
	}
	fmt.Println("err", err)
	return &pb.AddFollowsResp{}, nil
}
