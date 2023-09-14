package logic

import (
	"context"
	"errors"

	"database/sql"
	constants "doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"

	"github.com/zeromicro/go-zero/core/stores/builder"

	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLikeInfoLogic struct {
	ctx  context.Context
	conn sqlx.SqlConn
	rds  *redis.Redis
	logx.Logger
}

func NewAddLikeInfoLogic(ctx context.Context, conn sqlx.SqlConn, rds *redis.Redis) *AddLikeInfoLogic {
	return &AddLikeInfoLogic{
		ctx:    ctx,
		conn:   conn,
		rds:    rds,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddLikeInfoLogic) AddLikeInfo() {
	/*
		Author：    刘洋
		Function：  定时将 redis 中点赞数据同步到数据库
		Update：    08.21
	*/
	redisClient := l.rds
	videoLikedKeyPrefix := constants.LikeCacheVideoLikedPrefix

	// 1.获取 redis 中所有点赞信息的 key
	var keys []string
	var cur uint64 = 0
	var count int64 = 1000
	for {
		// redis.scan 从 cursor=0 开始，当返回 nextCursor=0 时说明遍历完
		keysResp, nextCur, err0 := redisClient.Scan(cur, videoLikedKeyPrefix+"*", count)
		if err0 != nil && err0 != redis.Nil {
			fmt.Println("【redis访问错误，同步数据到Mysql失败！】")
		}
		keys = append(keys, keysResp...)
		if nextCur == 0 {
			break
		}
		cur = nextCur
	}
	if len(keys) == 0 {
		fmt.Println("【Redis同步数据到Mysql完成】")
		return
	}
	fmt.Println(keys)
	// 2.根据所有点赞信息中的 key，获取所有(key,field,val)对
	hash := make(map[string]string)
	for _, item := range keys {
		// 2.1根据 videoId 查 Redis 得到 userId
		var cur uint64 = 0
		var count int64 = 1000
		arr := strings.Split(item, ":")
		var videoId = arr[2]
		for {
			// redis.scan 从 cursor=0 开始，当返回 nextCursor=0 时说明遍历完
			keysResp, nextCur, err0 := redisClient.Hscan(item, cur, "*", count)
			if err0 != nil && err0 != redis.Nil {
				fmt.Println("【redis访问错误，同步数据到Mysql失败！】")
			}
			for i := 0; i < len(keysResp); i += 2 {
				hash[videoId+":"+keysResp[i]] = keysResp[i+1]
			}
			if nextCur == 0 {
				break
			}
			cur = nextCur
		}
	}
	fmt.Println(hash)
	// 3.根据所有 videoId(key), userId(field) 查询并新增 favorite 表
	type Favorite struct {
		Id         int64     `db:"id"`          // 主键
		VideoId    int64     `db:"video_id"`    // 视频id
		UserId     int64     `db:"user_id"`     // 视频id
		CreateTime time.Time `db:"create_time"` // 该条记录创建时间
		UpdateTime time.Time `db:"update_time"` // 该条最后一次更新时间
		IsDelete   int64     `db:"is_delete"`   // 逻辑删除
	}
	var (
		favoriteFieldNames          = builder.RawFieldNames(&Favorite{})
		favoriteRows                = strings.Join(favoriteFieldNames, ",")
		favoriteRowsExpectAutoSet   = strings.Join(stringx.Remove(favoriteFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
		favoriteRowsWithPlaceHolder = strings.Join(stringx.Remove(favoriteFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
	)
	var insertBatch string
	for k, v := range hash {
		arr := strings.Split(k, ":")
		videoId, _ := strconv.ParseInt(arr[0], 10, 64)
		userId, _ := strconv.ParseInt(arr[1], 10, 64)
		isDelete, _ := strconv.ParseInt(v, 10, 64)
		// 3.1 根据 userId,videoId 查询 favorite 表项
		var resp Favorite
		query := fmt.Sprintf("select %s from `favorite` where `user_id` = ? and `video_id` = ? limit 1", favoriteRows)
		err := l.conn.QueryRow(&resp, query, userId, videoId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			fmt.Println("【Redis同步数据, Mysql更新失败】")
		}
		if resp.Id != 0 {
			// 3.2 如果记录存在，则进行更新
			query := fmt.Sprintf("update `favorite` set %s where `id` = ?", favoriteRowsWithPlaceHolder)
			_, err := l.conn.Exec(query, videoId, userId, isDelete, resp.Id)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				fmt.Println("【Redis同步数据, Mysql更新失败】")
			}
		} else {
			// 3.3 如果记录不存在，则拼接字符串 insertBatch 等到后续批量插入
			snowflake, err1 := util.NewSnowflake(constants.OtherMachineId)
			if err1 != nil && !errors.Is(err, sql.ErrNoRows) {
				fmt.Println("【Redis同步数据, Mysql更新失败时snowflake生成id失败】")
			}
			snowId := snowflake.Generate() //雪花算法生成id
			insertBatch += "(" + strconv.FormatInt(snowId, 10) + "," + arr[0] + "," + arr[1] + "," + v + ")" + ","
		}

	}
	if len(insertBatch) != 0 {
		insertBatch = insertBatch[:len(insertBatch)-1] // 去掉最后一个逗号
		query := fmt.Sprintf("insert into `favorite` (%s) values %s ", favoriteRowsExpectAutoSet, insertBatch)
		_, err := l.conn.Exec(query)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			fmt.Println("【Redis同步数据, Mysql更新失败】")
		}
		fmt.Println(insertBatch)
	}
	fmt.Println("【Redis同步数据到Mysql完成】")
	return
}
