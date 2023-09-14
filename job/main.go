package main

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/job/logic"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

/*
@yearly (or @annually)	每年1月1日午夜跑步一次
@monthly	每个月第一天的午夜跑一次
@daily (or @midnight)	每天午夜跑一次
@hourly	每小时运行一次
@every <duration>	every duration    todo 尽量使用该模式使得代码简明
*/

func main() {
	// 数据库连接参数
	conn := sqlx.NewMysql("root:liujun@tcp(8.137.50.160:3306)/liujun_content?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")

	//redis连接
	conf := redis.RedisConf{
		Host: "127.0.0.1:8094",
		Pass: common.DefaultPass,
		Type: "node",
	}
	rds := redis.MustNewRedis(conf)
	ctx := context.Background()

	//定时任务
	c := cron.New()

	// 标准构建
	_, err := c.AddFunc("@every 5s", logic.NewAddLikeInfoLogic(ctx, conn, rds).AddLikeInfo)
	if err != nil {
		fmt.Println("点赞数据持久化定时任务添加失败")
	}
	_, err = c.AddFunc("@every 24h", logic.NewCheckLogLogic(ctx).CheckLog)
	if err != nil {
		fmt.Println("日志清理定时任务添加失败")
	}

	//定时任务启动
	c.Start()

	// 防止程序提前结束
	select {}
}
