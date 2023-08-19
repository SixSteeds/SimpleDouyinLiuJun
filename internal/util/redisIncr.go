package util

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func IncreaseOrSet(ctx context.Context, redisClient *redis.Client, key string) (int64, error) {
	txf := func(tx *redis.Tx) error {
		// 在事务中获取键的值
		cmd := tx.Get(ctx, key)
		oldValue, err := cmd.Result()
		if err != nil && err != redis.Nil {
			return err
		}

		// 键不存在，设置初始值为1
		if err == redis.Nil {
			cmd := tx.Set(ctx, key, 1, 0)
			_, err := cmd.Result()
			if err != nil {
				return err
			}
			return nil
		}

		// 键存在，将值加1
		cmd = tx.Incr(ctx, key)
		_, err = cmd.Result()
		if err != nil {
			return err
		}

		return nil
	}

	err := redisClient.Watch(ctx, txf, key)
	if err == redis.TxFailedErr {
		// 事务失败，可以选择重试或处理冲突
		return 0, err
	} else if err != nil {
		// 发生了其他错误
		return 0, err
	}

	newValue, err := redisClient.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}

	return newValue, nil
}
