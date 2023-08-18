// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
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
	videoFieldNames = builder.RawFieldNames(&Video{})
	videoRows       = strings.Join(videoFieldNames, ",")
	// 此处去掉id以便于使用雪花算法生成id，加上is_delete以便于软删除
	videoRowsExpectAutoSet   = strings.Join(stringx.Remove(videoFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`is_delete`", "`update_time`", "`updated_at`"), ",")
	videoRowsWithPlaceHolder = strings.Join(stringx.Remove(videoFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheLiujunContentVideoIdPrefix = "cache:liujunContent:video:id:"
)

type (
	videoModel interface {
		Insert(ctx context.Context, data *Video) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Video, error)
		Update(ctx context.Context, data *Video) error
		Delete(ctx context.Context, id int64) error
		GetFeedList(ctx context.Context, user_id int64, latest_time int64, size int64) ([]FeedVideo, error)
	}

	defaultVideoModel struct {
		sqlc.CachedConn
		table string
	}

	Video struct {
		Id         int64          `db:"id"`          // 主键
		UserId     int64          `db:"user_id"`     // 视频作者id
		PlayUrl    string         `db:"play_url"`    // 视频播放地址
		CoverUrl   sql.NullString `db:"cover_url"`   // 视频封面地址
		Title      sql.NullString `db:"title"`       // 视频标题
		CreateTime time.Time      `db:"create_time"` // 该条记录创建时间
		UpdateTime time.Time      `db:"update_time"` // 该条最后一次更新时间
		IsDelete   int64          `db:"is_delete"`   // 逻辑删除
	}

	FeedVideo struct {
		Id            int64  `db:"id"`             // 主键
		UserId        int64  `db:"user_id"`        // 视频作者id
		PlayUrl       string `db:"play_url"`       // 视频播放地址
		CoverUrl      string `db:"cover_url"`      // 视频封面地址
		Title         string `db:"title"`          // 视频标题
		FavoriteCount int64  `db:"favorite_count"` // 视频被收藏次数
		CommentCount  int64  `db:"comment_count"`  // 视频被评论次数
		IsFavorite    int64  `db:"is_favorite"`    // 是否被当前用户点赞
	}
)

func newVideoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultVideoModel {
	return &defaultVideoModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`video`",
	}
}

func (m *defaultVideoModel) withSession(session sqlx.Session) *defaultVideoModel {
	return &defaultVideoModel{
		CachedConn: m.CachedConn.WithSession(session),
		table:      "`video`",
	}
}

func (m *defaultVideoModel) Delete(ctx context.Context, id int64) error {
	liujunContentVideoIdKey := fmt.Sprintf("%s%v", cacheLiujunContentVideoIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, liujunContentVideoIdKey)
	return err
}

func (m *defaultVideoModel) FindOne(ctx context.Context, id int64) (*Video, error) {
	liujunContentVideoIdKey := fmt.Sprintf("%s%v", cacheLiujunContentVideoIdPrefix, id)
	var resp Video
	err := m.QueryRowCtx(ctx, &resp, liujunContentVideoIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", videoRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultVideoModel) getFeedList(ctx context.Context, user_id int64, latest_time int64, size int64) (*[]FeedVideo, error) {
	var resp []FeedVideo
	query := fmt.Sprintf("SELECT   "+
		"v.id,"+
		"v.user_id,"+
		"v.play_url,"+
		"v.cover_url,"+
		"v.title,"+
		"(SELECT COUNT(*) FROM favorite WHERE video_id = v.id) AS favorite_count,"+
		"(ELECT COUNT(*) FROM comment WHERE video_id = v.id) AS comment_count,"+
		"CASE WHEN EXISTS (SELECT 1 FROM favorite WHERE video_id = v.id AND user_id = ?) THEN true ELSE false END AS is_favorite"+
		"FROM video v"+
		"WHERE v.is_delete = 0 and v.update_time<?"+
		"ORDER BY v.create_time DESC"+
		"limit ?;", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, user_id, latest_time, size)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *defaultVideoModel) Insert(ctx context.Context, data *Video) (sql.Result, error) {
	liujunContentVideoIdKey := fmt.Sprintf("%s%v", cacheLiujunContentVideoIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, videoRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.PlayUrl, data.CoverUrl, data.Title, data.IsDelete)
	}, liujunContentVideoIdKey)
	return ret, err
}

func (m *defaultVideoModel) Update(ctx context.Context, data *Video) error {
	liujunContentVideoIdKey := fmt.Sprintf("%s%v", cacheLiujunContentVideoIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, videoRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.UserId, data.PlayUrl, data.CoverUrl, data.Title, data.IsDelete, data.Id)
	}, liujunContentVideoIdKey)
	return err
}

func (m *defaultVideoModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheLiujunContentVideoIdPrefix, primary)
}

func (m *defaultVideoModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", videoRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultVideoModel) tableName() string {
	return m.table
}
