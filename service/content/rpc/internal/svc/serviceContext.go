package svc

import (
	"doushen_by_liujun/service/content/rpc/internal/config"
	genModel "doushen_by_liujun/service/content/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	CommentModel        genModel.CommentModel
	CommentForUserModel genModel.CommentForUserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		CommentModel:        genModel.NewCommentModel(sqlConn, c.Cache),
		CommentForUserModel: genModel.NewCommentForUserModel(sqlConn, c.Cache),
		Config:              c,
	}
}
