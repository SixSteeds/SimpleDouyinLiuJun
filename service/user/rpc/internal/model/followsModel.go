package genModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FollowsModel = (*customFollowsModel)(nil)

type (
	// FollowsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowsModel.
	FollowsModel interface {
		followsModel
	}

	customFollowsModel struct {
		*defaultFollowsModel
	}
)

// NewFollowsModel returns a model for the database table.
func NewFollowsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FollowsModel {
	return &customFollowsModel{
		defaultFollowsModel: newFollowsModel(conn, c, opts...),
	}
}
