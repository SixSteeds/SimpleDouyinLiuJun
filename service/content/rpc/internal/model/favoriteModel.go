package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FavoriteModel = (*customFavoriteModel)(nil)

type (
	// FavoriteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFavoriteModel.
	FavoriteModel interface {
		favoriteModel
	}

	customFavoriteModel struct {
		*defaultFavoriteModel
	}
)

// NewFavoriteModel returns a model for the database table.
func NewFavoriteModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FavoriteModel {
	return &customFavoriteModel{
		defaultFavoriteModel: newFavoriteModel(conn, c, opts...),
	}
}
