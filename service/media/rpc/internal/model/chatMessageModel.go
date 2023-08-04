package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatMessageModel = (*customChatMessageModel)(nil)

type (
	// ChatMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatMessageModel.
	ChatMessageModel interface {
		chatMessageModel
	}

	customChatMessageModel struct {
		*defaultChatMessageModel
	}
)

// NewChatMessageModel returns a model for the database table.
func NewChatMessageModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ChatMessageModel {
	return &customChatMessageModel{
		defaultChatMessageModel: newChatMessageModel(conn, c, opts...),
	}
}
