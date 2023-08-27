package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource string
	}
	Cache                cache.CacheConf
	RedisConf            redis.RedisConf
	ContentRpcConf       zrpc.RpcClientConf
	DefaultAvatar        []string
	DefaultBackgroundImg []string
}
