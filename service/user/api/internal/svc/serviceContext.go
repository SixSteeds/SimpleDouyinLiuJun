package svc

import (
	gloabmiddleware "doushen_by_liujun/internal/middleware"
	"doushen_by_liujun/service/content/rpc/content"
	"doushen_by_liujun/service/user/api/internal/config"
	"doushen_by_liujun/service/user/api/internal/middleware"
	"doushen_by_liujun/service/user/rpc/user"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                 config.Config
	UserAgentMiddleware    rest.Middleware
	UserRpcClient          user.User
	JwtAuthMiddleware      rest.Middleware
	RedisClient            *redis.Redis
	LoginLogKqPusherClient *kq.Pusher
	ContentRpcClient       content.Content
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                 c,
		UserAgentMiddleware:    middleware.NewUserAgentMiddleware().Handle,
		UserRpcClient:          user.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
		RedisClient:            redis.MustNewRedis(c.RedisConf),
		JwtAuthMiddleware:      gloabmiddleware.NewJwtAuthMiddleware().Handle,
		LoginLogKqPusherClient: kq.NewPusher(c.LoginLogKqPusherConf.Brokers, c.LoginLogKqPusherConf.Topic),
		ContentRpcClient:       content.NewContent(zrpc.MustNewClient(c.ContentRpcConf)),
	}
}
