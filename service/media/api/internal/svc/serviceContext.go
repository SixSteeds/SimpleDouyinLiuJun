package svc

import (
	gloabmiddleware "doushen_by_liujun/internal/middleware"
	"doushen_by_liujun/service/media/api/internal/config"
	"doushen_by_liujun/service/media/rpc/media"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	RedisClient       *redis.Redis
	JwtAuthMiddleware rest.Middleware
	MediaRpcClient    media.Media
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		RedisClient:       redis.MustNewRedis(c.RedisConf),
		MediaRpcClient:    media.NewMedia(zrpc.MustNewClient(c.MediaRpcConf)),
		JwtAuthMiddleware: gloabmiddleware.NewJwtAuthMiddleware().Handle,
	}
}
