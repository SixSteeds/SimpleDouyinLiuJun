package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRpcConf          zrpc.RpcClientConf
	ContentRpcConf       zrpc.RpcClientConf
	RedisConf            redis.RedisConf
	LoginLogKqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
