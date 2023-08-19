package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	MediaRpcConf      zrpc.RpcClientConf
	RedisConf         redis.RedisConf
	MediaKqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
