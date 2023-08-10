package svc

import (
	"doushen_by_liujun/service/content/rpc/internal/config"
	genModel "doushen_by_liujun/service/content/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	CommentModel  genModel.CommentModel
	FavoriteModel genModel.FavoriteModel
	VideoModel    genModel.VideoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:        c,
		CommentModel:  genModel.NewCommentModel(sqlConn, c.Cache),
		FavoriteModel: genModel.NewFavoriteModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		VideoModel:    genModel.NewVideoModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
