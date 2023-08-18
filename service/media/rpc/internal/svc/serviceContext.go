package svc

import (
	"doushen_by_liujun/service/media/rpc/internal/config"
	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config         config.Config
	KqPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		KqPusherClient: kq.NewPusher(c.MediaKqPusherConf.Brokers, c.MediaKqPusherConf.Topic),
	}
}
