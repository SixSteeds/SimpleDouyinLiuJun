package svc

import (
	gloabmiddleware "doushen_by_liujun/internal/middleware"
	"doushen_by_liujun/service/chat/api/internal/config"
	"doushen_by_liujun/service/chat/rpc/chat"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config            config.Config
	RedisClient       *redis.Redis
	JwtAuthMiddleware rest.Middleware
	ChatRpcClient     chat.Chat
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		RedisClient:       redis.MustNewRedis(c.RedisConf),
		JwtAuthMiddleware: gloabmiddleware.NewJwtAuthMiddleware().Handle,
		ChatRpcClient:     chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),
	}
}
