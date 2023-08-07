package svc

import (
	gloabmiddleware "doushen_by_liujun/internal/middleware"
	"doushen_by_liujun/service/content/api/internal/config"
	"doushen_by_liujun/service/content/rpc/content"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	RedisClient       *redis.Redis
	JwtAuthMiddleware rest.Middleware
	ContentRpcClient  content.Content
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		RedisClient:       redis.MustNewRedis(c.RedisConf),
		JwtAuthMiddleware: gloabmiddleware.NewJwtAuthMiddleware().Handle,
		ContentRpcClient:  content.NewContent(zrpc.MustNewClient(c.ContentRpcConf)),
	}
}
