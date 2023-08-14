package svc

import (
	"doushen_by_liujun/service/chat/rpc/internal/config"
	"doushen_by_liujun/service/chat/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	ChatMessageModel model.ChatMessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:           c,
		ChatMessageModel: model.NewChatMessageModel(sqlConn, c.Cache),
	}
}
