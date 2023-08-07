package genModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentForUserModel = (*customCommentForUserModel)(nil)

type (
	// CommentForUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentForUserModel.
	CommentForUserModel interface {
		commentForUserModel
	}

	customCommentForUserModel struct {
		*defaultCommentForUserModel
	}
)

// NewCommentForUserModel returns a model for the database table.
func NewCommentForUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentForUserModel {
	return &customCommentForUserModel{
		defaultCommentForUserModel: newCommentForUserModel(conn, c, opts...),
	}
}
