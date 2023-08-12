package svc

import (
	gloabmiddleware "doushen_by_liujun/internal/middleware"
	"doushen_by_liujun/service/user/api/internal/config"
	"doushen_by_liujun/service/user/api/internal/middleware"
	"doushen_by_liujun/service/user/rpc/user"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	UserAgentMiddleware rest.Middleware
	UserRpcClient       user.User
	JwtAuthMiddleware   rest.Middleware
	RedisClient         *redis.Redis
	KqPusherClient      *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		UserAgentMiddleware: middleware.NewUserAgentMiddleware().Handle,
		UserRpcClient:       user.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
		RedisClient:         redis.MustNewRedis(c.RedisConf),
		JwtAuthMiddleware:   gloabmiddleware.NewJwtAuthMiddleware().Handle,
		KqPusherClient:      kq.NewPusher(c.LoginKqPusherConf.Brokers, c.LoginKqPusherConf.Topic),
	}
}
