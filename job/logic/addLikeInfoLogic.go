package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

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
	// add your logic here
	fmt.Println("我被执行啦")
	return
}
