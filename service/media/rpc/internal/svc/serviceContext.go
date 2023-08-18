package svc

import (
	"doushen_by_liujun/service/media/rpc/internal/config"
	"doushen_by_liujun/service/media/rpc/internal/model"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	KqPusherClient *kq.Pusher
	MediaModel     model.VideoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{

		Config:     c,
		MediaModel: model.NewVideoModel(sqlConn, c.Cache),

		KqPusherClient: kq.NewPusher(c.MediaKqPusherConf.Brokers, c.MediaKqPusherConf.Topic),
	}
}
