package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserinfoModel = (*customUserinfoModel)(nil)

type (
	// UserinfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserinfoModel.
	UserinfoModel interface {
		userinfoModel
	}

	customUserinfoModel struct {
		*defaultUserinfoModel
	}
)

// NewUserinfoModel returns a model for the database table.
func NewUserinfoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserinfoModel {
	return &customUserinfoModel{
		defaultUserinfoModel: newUserinfoModel(conn, c, opts...),
	}
}
