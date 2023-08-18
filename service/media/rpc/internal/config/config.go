package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MediaKqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
