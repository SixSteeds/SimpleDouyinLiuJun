package logic

import (
	"context"
	"database/sql"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"log"
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
	// todo: add your logic here and delete this line
	snowflake, err := util.NewSnowflake(common.UserRpcMachineId)
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_rpc_saveUserLogic_SaveUser_NewSnowflake_false"); err != nil {
			log.Fatal(err)
		}
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

	_, err = l.svcCtx.UserinfoModel.Insert(l.ctx, &model.Userinfo{
		Id:       snowId,
		Username: username,
		Password: password,
	})
	if err != nil {
		if err := l.svcCtx.KqPusherClient.Push("user_rpc_saveUserLogic_SaveUser_insert_false"); err != nil {
			log.Fatal(err)
		}
		return &pb.SaveUserResp{
			Success: false,
		}, err
	}
	//tzx在创建用户时向数据库中写入自己关注自己，保证刷到自己的视频时不会出现红加号
	//在取消关注时加以限制，不允许取消关注自己
	data, err := util.NewSnowflake(47)
	if err != nil {
		l.Logger.Info("雪花算法报错", err)
		if err := l.svcCtx.KqPusherClient.Push("user_rpc_addFollowsLogic_AddFollows_NewSnowflake_false"); err != nil {
			log.Fatal(err)
		}
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

	if err := l.svcCtx.KqPusherClient.Push("user_rpc_saveUserLogic_SaveUser_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.SaveUserResp{
		Success: true,
		Id:      snowId,
	}, nil
}
