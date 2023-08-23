package svc

import (
	"doushen_by_liujun/service/user/rpc/internal/config"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	UserinfoModel model.UserinfoModel
	FollowsModel  model.FollowsModel
	RedisClient   *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:        c,
		UserinfoModel: model.NewUserinfoModel(sqlConn, c.Cache),
		FollowsModel:  model.NewFollowsModel(sqlConn, c.Cache),
		RedisClient:   redis.MustNewRedis(c.RedisConf),
	}
}
