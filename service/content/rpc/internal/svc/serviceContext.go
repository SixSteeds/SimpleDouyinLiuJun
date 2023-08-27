package svc

import (
	"doushen_by_liujun/service/content/rpc/internal/config"
	"doushen_by_liujun/service/content/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	CommentModel  model.CommentModel
	FavoriteModel model.FavoriteModel
	VideoModel    model.VideoModel
	RedisClient   *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:        c,
		CommentModel:  model.NewCommentModel(sqlConn, c.Cache),
		FavoriteModel: model.NewFavoriteModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		VideoModel:    model.NewVideoModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		RedisClient:   redis.MustNewRedis(c.RedisConf),
	}
}
