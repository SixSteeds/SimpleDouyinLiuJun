package logic

import (
	"context"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"errors"
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
	l.Logger.Info(in)
	data, err := util.NewSnowflake(47)
	if err != nil {
		l.Logger.Info("雪花算法报错", err)
		return nil, errors.New("雪花算法报错")
	}
	_, err = l.svcCtx.FollowsModel.Insert(l.ctx, &model.Follows{
		Id:       data.Generate(),
		UserId:   in.UserId,
		FollowId: in.FollowId,
		//CreateTime: time.Now().Unix(),
		//UpdateTime: time.Now().Unix(),
		IsDelete: 0,
	})
	if err != nil {
		l.Logger.Error("写入数据库报错", err)
		return nil, errors.New("写入数据库报错")
	}
	return &pb.AddFollowsResp{}, nil
}
