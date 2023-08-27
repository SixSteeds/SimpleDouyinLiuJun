package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ContentRpcConf               zrpc.RpcClientConf
	UserRpcConf                  zrpc.RpcClientConf
	RedisConf                    redis.RedisConf
	UploadPersistentKqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
