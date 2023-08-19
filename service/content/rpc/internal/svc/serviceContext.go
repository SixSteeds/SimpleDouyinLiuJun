package svc

import (
	"doushen_by_liujun/service/content/rpc/internal/config"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"doushen_by_liujun/service/user/rpc/user"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	CommentModel   model.CommentModel
	FavoriteModel  model.FavoriteModel
	VideoModel     model.VideoModel
	KqPusherClient *kq.Pusher

	UserRpcClient user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:         c,
		CommentModel:   model.NewCommentModel(sqlConn, c.Cache),
		FavoriteModel:  model.NewFavoriteModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		VideoModel:     model.NewVideoModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		KqPusherClient: kq.NewPusher(c.ContentKqPusherConf.Brokers, c.ContentKqPusherConf.Topic),

		UserRpcClient: user.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
