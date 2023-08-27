package logic

import (
	"context"
	"database/sql"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"math/rand"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveUserLogic {
	return &SaveUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveUserLogic) SaveUser(in *pb.SaveUserReq) (*pb.SaveUserResp, error) {
	l.Logger.Info("SaveUser方法请求参数：", in)
	snowflake, err := util.NewSnowflake(common.UserRpcMachineId)
	if err != nil {
		return &pb.SaveUserResp{
			Success: false,
		}, nil
	}
	snowId := snowflake.Generate()
	username := sql.NullString{
		String: in.Username,
		Valid:  true,
	}
	password := sql.NullString{
		String: in.Password,
		Valid:  true,
	}
	defaultBackgroundImg := sql.NullString{
		String: l.svcCtx.DefaultBackgroundImg[rand.Intn(len(l.svcCtx.DefaultBackgroundImg))],
		Valid:  true,
	}
	defaultAvatar := sql.NullString{
		String: l.svcCtx.DefaultAvatar[rand.Intn(len(l.svcCtx.DefaultAvatar))],
		Valid:  true,
	}

	_, err = l.svcCtx.UserinfoModel.Insert(l.ctx, &model.Userinfo{
		Id:              snowId,
		Username:        username,
		Password:        password,
		BackgroundImage: defaultBackgroundImg,
		Avatar:          defaultAvatar,
	})
	if err != nil {
		return &pb.SaveUserResp{
			Success: false,
		}, err
	}
	// tzx在创建用户时向数据库中写入自己关注自己，保证刷到自己的视频时不会出现红加号
	// 在取消关注时加以限制，不允许取消关注自己
	data, err := util.NewSnowflake(common.UserRpcMachineId)
	if err != nil {
		l.Logger.Info("雪花算法报错", err)
		return &pb.SaveUserResp{
			Success: false,
		}, nil
	}
	_, err = l.svcCtx.FollowsModel.Insert(l.ctx, &model.Follows{
		Id:       data.Generate(),
		UserId:   strconv.FormatInt(snowId, 10),
		FollowId: strconv.FormatInt(snowId, 10),
		IsDelete: 0,
	})
	if err != nil {
		return &pb.SaveUserResp{
			Success: false,
		}, nil
	}
	return &pb.SaveUserResp{
		Success: true,
		Id:      snowId,
	}, nil
}
