package logic

import (
	"context"
	"database/sql"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"log"

	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"

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
	snowflake, err := util.NewSnowflake(2)
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

	logx.Error("save user success")
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_saveUserLogic_SaveUser_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.SaveUserResp{
		Success: true,
		Id:      snowId,
	}, nil
}
