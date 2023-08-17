package svc

import (
	"doushen_by_liujun/service/user/rpc/internal/config"
	"doushen_by_liujun/service/user/rpc/internal/model"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	UserinfoModel  model.UserinfoModel
	FollowsModel   model.FollowsModel
	KqPusherClient *kq.Pusher
	RedisClient    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:         c,
		UserinfoModel:  model.NewUserinfoModel(sqlConn, c.Cache),
		FollowsModel:   model.NewFollowsModel(sqlConn, c.Cache),
		KqPusherClient: kq.NewPusher(c.LoginKqPusherConf.Brokers, c.LoginKqPusherConf.Topic),
		RedisClient:    redis.MustNewRedis(c.RedisConf),
	}
}
