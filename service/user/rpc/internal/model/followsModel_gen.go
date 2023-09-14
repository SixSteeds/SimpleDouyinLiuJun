package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	followsFieldNames = builder.RawFieldNames(&Follows{})
	followsRows       = strings.Join(followsFieldNames, ",")

	followsRowsExpectAutoSet   = strings.Join(stringx.Remove(followsFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	followsRowsWithPlaceHolder = strings.Join(stringx.Remove(followsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheLiujunUserFollowsIdPrefix = "cache:liujunUser:follows:id:"
)

type (
	followsModel interface {
		Insert(ctx context.Context, data *Follows) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Follows, error)
		Update(ctx context.Context, data *Follows) error
		Delete(ctx context.Context, id int64) error
		DeleteByUserIdAndFollowId(ctx context.Context, userId string, followId string) error
		FindByUserId(ctx context.Context, userId int64) (*[]*FollowUser, error)
		FindByFollowId(ctx context.Context, userId int64) (*[]*FollowUser, error)
		FindFriendsByUserId(ctx context.Context, userId int64) (*[]*FollowUser, error)
		FindFollowsCount(ctx context.Context, userId int64) (int, error)
		FindFollowersCount(ctx context.Context, userId int64) (int, error)
		CheckIsFollowed(ctx context.Context, userId int64, followId int64) (bool, error)
	}

	defaultFollowsModel struct {
		sqlc.CachedConn
		table string
	}

	Follows struct {
		Id         int64     `db:"id"`          // 主键
		UserId     string    `db:"user_id"`     // 关注的人
		FollowId   string    `db:"follow_id"`   // 被关注的人
		CreateTime time.Time `db:"create_time"` // 该条记录创建时间
		UpdateTime time.Time `db:"update_time"` // 该条最后一次信息修改时间
		IsDelete   int64     `db:"is_delete"`   // 逻辑删除
	}

	FollowUser struct {
		Id              int64          `db:"id"` // 主键
		UserName        sql.NullString `db:"username"`
		Avator          sql.NullString `db:"avator"`
		BackgroundImage sql.NullString `db:"background_image"`
		Signature       sql.NullString `db:"signature"`
		IsFollow        bool           `db:"is_follow"`
	}
)

func newFollowsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultFollowsModel {
	return &defaultFollowsModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`follows`",
	}
}

func (m *defaultFollowsModel) withSession(session sqlx.Session) *defaultFollowsModel {
	return &defaultFollowsModel{
		CachedConn: m.CachedConn.WithSession(session),
		table:      "`follows`",
	}
}

func (m *defaultFollowsModel) Delete(ctx context.Context, id int64) error {
	liujunUserFollowsIdKey := fmt.Sprintf("%s%v", cacheLiujunUserFollowsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, liujunUserFollowsIdKey)
	return err
}
func (m *defaultFollowsModel) DeleteByUserIdAndFollowId(ctx context.Context, userId string, followId string) error {
	liujunUserFollowsIdKey := fmt.Sprintf("%s%v", cacheLiujunUserFollowsIdPrefix, userId+followId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("DELETE FROM %s WHERE `user_id` = ? AND `follow_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, userId, followId)
	}, liujunUserFollowsIdKey)
	return err
}

func (m *defaultFollowsModel) FindOne(ctx context.Context, id int64) (*Follows, error) {
	liujunUserFollowsIdKey := fmt.Sprintf("%s%v", cacheLiujunUserFollowsIdPrefix, id)
	var resp Follows
	err := m.QueryRowCtx(ctx, &resp, liujunUserFollowsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? ", followsRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFollowsModel) FindByUserId(ctx context.Context, id int64) (*[]*FollowUser, error) {
	var resp []*FollowUser
	query := fmt.Sprintf("select u.id,u.username,u.avatar,u.background_image,u.signature," +
		"TRUE AS is_follow" +
		" from userinfo u,follows f where f.user_id = ? and f.follow_id = u.id and u.id <> f.user_id and u.is_delete = 0")
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, id)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFollowsModel) FindByFollowId(ctx context.Context, id int64) (*[]*FollowUser, error) {
	var resp []*FollowUser
	query := fmt.Sprintf("select u.id,u.username,u.avatar,u.background_image,u.signature," +
		"EXISTS (SELECT 1 FROM follows WHERE user_id = ? AND follow_id = u.id) AS is_follow" +
		" from userinfo u,follows f where f.follow_id = ? and f.user_id = u.id and u.id <> f.follow_id and u.is_delete = 0")
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, id, id)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFollowsModel) FindFriendsByUserId(ctx context.Context, id int64) (*[]*FollowUser, error) {
	var resp []*FollowUser
	query := fmt.Sprintf("select u.id,u.username,u.avatar,u.background_image,u.signature," +
		"TRUE AS is_follow" +
		" from userinfo u,follows f,follows f2 where f.follow_id = f2.user_id and f2.follow_id = f.user_id  and u.id <> f.user_id and f.user_id = ? and f.follow_id = u.id and u.is_delete = 0")
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, id)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultFollowsModel) FindFollowsCount(ctx context.Context, id int64) (int, error) {
	var count int
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, id)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (m *defaultFollowsModel) FindFollowersCount(ctx context.Context, id int64) (int, error) {
	var count int
	query := fmt.Sprintf("select count(*) from %s where `follow_id` = ?", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, id)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (m *defaultFollowsModel) CheckIsFollowed(ctx context.Context, userid int64, followid int64) (bool, error) {
	var count int
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `follow_id` = ?", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, userid, followid)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (m *defaultFollowsModel) Insert(ctx context.Context, data *Follows) (sql.Result, error) {
	liujunUserFollowsIdKey := fmt.Sprintf("%s%v", cacheLiujunUserFollowsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, followsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Id, data.UserId, data.FollowId, data.IsDelete)
	}, liujunUserFollowsIdKey)
	return ret, err
}

func (m *defaultFollowsModel) Update(ctx context.Context, data *Follows) error {
	liujunUserFollowsIdKey := fmt.Sprintf("%s%v", cacheLiujunUserFollowsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, followsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.UserId, data.FollowId, data.IsDelete, data.Id)
	}, liujunUserFollowsIdKey)
	return err
}

func (m *defaultFollowsModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheLiujunUserFollowsIdPrefix, primary)
}

func (m *defaultFollowsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", followsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultFollowsModel) tableName() string {
	return m.table
}
