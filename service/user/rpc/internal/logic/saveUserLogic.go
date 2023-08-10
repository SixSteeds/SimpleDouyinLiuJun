package logic

import (
	"context"
	"database/sql"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/user/rpc/internal/model"

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
		return &pb.SaveUserResp{
			Success: false,
		}, err
	}

	logx.Error("save user success")
	return &pb.SaveUserResp{
		Success: true,
		Id:      snowId,
	}, nil
}
