package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userinfoFieldNames = builder.RawFieldNames(&Userinfo{})
	userinfoRows       = strings.Join(userinfoFieldNames, ",")

	// 此处去掉id以便于使用雪花算法生成id，加上is_delete以便于软删除
	userinfoRowsExpectAutoSet   = strings.Join(stringx.Remove(userinfoFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`is_delete`", "`update_time`", "`updated_at`"), ",")
	userinfoRowsWithPlaceHolder = strings.Join(stringx.Remove(userinfoFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheLiujunUserUserinfoIdPrefix = "cache:liujunUser:userinfo:id:"
)

type (
	userinfoModel interface {
		Insert(ctx context.Context, data *Userinfo) (sql.Result, error)
		FindOne(ctx context.Context, id int64, userId int64) (*Userdetail, error)
		Update(ctx context.Context, data *Userinfo) error
		Delete(ctx context.Context, id int64) error
		FindByIds(ctx context.Context, ids []int64, userId int64) (*[]*Userdetail, error)
		CheckOne(ctx context.Context, username string, password string) (*int64, error)
		FindUserById(ctx context.Context, id int64) (*Userinfo, error)
		FindUserListByIdList(ctx context.Context, userIdList *[]int64) (*[]*Userinfo, error)
		GetPasswordByUsername(ctx context.Context, username string) (*PasswordAndId, error)
	}

	defaultUserinfoModel struct {
		sqlc.CachedConn
		table string
	}
	Userinfo struct {
		Id              int64          `db:"id"`               // 主键
		Username        sql.NullString `db:"username"`         // 账号
		Password        sql.NullString `db:"password"`         // 密码
		Avatar          sql.NullString `db:"avatar"`           // 头像
		BackgroundImage sql.NullString `db:"background_image"` // 头像
		Signature       sql.NullString `db:"signature"`        // 个人简介
		CreateTime      time.Time      `db:"create_time"`      // 该条记录创建时间
		UpdateTime      time.Time      `db:"update_time"`      // 该条最后一次更新时间
		IsDelete        int64          `db:"is_delete"`        // 逻辑删除
		Name            sql.NullString `db:"name"`
	}
	Userdetail struct {
		Id              int64          `db:"id"`               // 主键
		Username        sql.NullString `db:"username"`         // 账号
		Avatar          sql.NullString `db:"avatar"`           // 头像
		BackgroundImage sql.NullString `db:"background_image"` // 头像
		Signature       sql.NullString `db:"signature"`        // 个人简介
		IsFollow        bool           `db:"is_follow"`
		Name            sql.NullString `db:"name"`
	}

	PasswordAndId struct {
		Password string `db:"password"`
		Id       int64  `db:"id"`
	}
)

func newUserinfoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultUserinfoModel {
	return &defaultUserinfoModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`userinfo`",
	}
}

func (m *defaultUserinfoModel) withSession(session sqlx.Session) *defaultUserinfoModel {
	return &defaultUserinfoModel{
		CachedConn: m.CachedConn.WithSession(session),
		table:      "`userinfo`",
	}
}

func (m *defaultUserinfoModel) FindUserListByIdList(ctx context.Context, userIdList *[]int64) (*[]*Userinfo, error) {
	var resp []*Userinfo
	// []int64 转换为 “,” 分隔的 string
	var strArr = make([]string, len(*userIdList))
	for k, v := range *userIdList {
		strArr[k] = fmt.Sprintf("%d", v)
	}
	var IdListStr = strings.Join(strArr, ",")
	query := fmt.Sprintf("select %s from %s where `id` in (%s) and `is_delete`!= '1'", userinfoRows, m.table, IdListStr)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserinfoModel) Delete(ctx context.Context, id int64) error {
	liujunUserUserinfoIdKey := fmt.Sprintf("%s%v", cacheLiujunUserUserinfoIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, liujunUserUserinfoIdKey)
	return err
}
func (m *defaultUserinfoModel) FindByIds(ctx context.Context, ids []int64, userId int64) (*[]*Userdetail, error) {
	var resp []*Userdetail
	var idStrings []string
	for _, id := range ids {
		idStrings = append(idStrings, "select u.id,u.username,u.avatar,u.background_image,u.signature,u.name,"+
			"EXISTS (SELECT 1 FROM follows WHERE user_id = "+strconv.FormatInt(userId, 10)+" AND follow_id = u.id) AS is_follow"+
			" from userinfo u where u.id = "+strconv.FormatInt(id, 10)+" and u.is_delete = 0")
	}
	combined := strings.Join(idStrings, " UNION ALL ")
	err := m.QueryRowsNoCacheCtx(ctx, &resp, combined)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultUserinfoModel) FindOne(ctx context.Context, id int64, userId int64) (*Userdetail, error) {
	liujunUserUserinfoIdKey := fmt.Sprintf("%s%v%v", cacheLiujunUserUserinfoIdPrefix, id, userId)
	var resp Userdetail
	err := m.QueryRowCtx(ctx, &resp, liujunUserUserinfoIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select u.id,u.username,u.avatar,u.background_image,u.signature,u.name," +
			"EXISTS (SELECT 1 FROM follows WHERE user_id = ? AND follow_id = u.id) AS is_follow" +
			" from userinfo u where u.id = ? and u.is_delete = 0")
		return conn.QueryRowCtx(ctx, v, query, userId, id)
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

func (m *defaultUserinfoModel) GetPasswordByUsername(ctx context.Context, username string) (*PasswordAndId, error) {

	var resp PasswordAndId
	query := fmt.Sprintf("select `password`,`id` from %s where `username` = ? ", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, username)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserinfoModel) FindUserById(ctx context.Context, id int64) (*Userinfo, error) {
	var resp Userinfo
	query := fmt.Sprintf("select %s from %s where `id` = ? ", userinfoRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, id)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserinfoModel) CheckOne(ctx context.Context, username, password string) (*int64, error) {

	var resp int64
	query := fmt.Sprintf("SELECT id FROM %s WHERE username = ? AND password = ?", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, username, password)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserinfoModel) Insert(ctx context.Context, data *Userinfo) (sql.Result, error) {
	liujunUserUserinfoIdKey := fmt.Sprintf("%s%v", cacheLiujunUserUserinfoIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, userinfoRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Id, data.Username, data.Password, data.Avatar, data.BackgroundImage, data.Signature, data.Name)
	}, liujunUserUserinfoIdKey)
	return ret, err
}

func (m *defaultUserinfoModel) Update(ctx context.Context, data *Userinfo) error {
	liujunUserUserinfoIdKey := fmt.Sprintf("%s%v", cacheLiujunUserUserinfoIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userinfoRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Username, data.Password, data.Avatar, data.BackgroundImage, data.Signature, data.IsDelete, data.Name, data.Id)
	}, liujunUserUserinfoIdKey)
	return err
}

func (m *defaultUserinfoModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheLiujunUserUserinfoIdPrefix, primary)
}

func (m *defaultUserinfoModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userinfoRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserinfoModel) tableName() string {
	return m.table
}
