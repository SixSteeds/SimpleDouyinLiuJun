package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
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

func (l *AddFollowsLogic) AddFollows(in *pb.AddFollowsReq) (*pb.AddFollowsResp, error) {
	l.Logger.Info(in)
	userid, _ := strconv.ParseInt(in.UserId, 10, 64)
	followid, _ := strconv.ParseInt(in.FollowId, 10, 64)
	isFollowed, err := l.svcCtx.FollowsModel.CheckIsFollowed(l.ctx, userid, followid) //检查是不是重复操作
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	if isFollowed { //已经关注过了
		return &pb.AddFollowsResp{}, nil
	}
	data, err := util.NewSnowflake(common.UserRpcMachineId)
	if err != nil {
		l.Logger.Error("雪花算法报错", err)
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
