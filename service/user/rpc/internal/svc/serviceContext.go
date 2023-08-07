package svc

import (
	"doushen_by_liujun/service/user/rpc/internal/config"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	UserinfoModel genModel.UserinfoModel
	FollowsModel  genModel.FollowsModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:        c,
		UserinfoModel: genModel.NewUserinfoModel(sqlConn, c.Cache),
		FollowsModel:  genModel.NewFollowsModel(sqlConn, c.Cache),
	}
}
