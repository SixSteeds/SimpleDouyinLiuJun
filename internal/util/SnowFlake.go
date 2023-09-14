package util

import (
	"doushen_by_liujun/internal/common"
	"errors"
	"sync"
	"time"
)

// A Snowflake struct holds the basic information needed for a snowflake generator worker
type Snowflake struct {
	sync.Mutex
	timestamp int64
	workerid  int64
	sequence  int64
}

// NewSnowflake NewNode returns a new snowflake worker that can be used to generate snowflake IDs
func NewSnowflake(workerid int64) (*Snowflake, error) {
	// 检查workerid是否在有效范围内（0到1023）
	if workerid < 0 || workerid > common.WorkeridMax {
		return nil, errors.New("workerid must be between 0 and 1023")
	}
	// 创建并返回Snowflake对象
	return &Snowflake{
		timestamp: 0,
		workerid:  workerid,
		sequence:  0,
	}, nil
}

// Generate creates and returns a unique snowflake ID
func (s *Snowflake) Generate() int64 {
	// 为了保证线程安全，对整个函数进行加锁
	s.Lock()
	// 获取当前时间的毫秒数
	now := time.Now().UnixNano() / 1000000
	// 如果当前时间与上次生成ID的时间相同，则自增序列号
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & common.SequenceMask
		// 如果自增序列号达到了最大值，则等待直到下一毫秒再生成ID
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 如果当前时间与上次生成ID的时间不同，则重置序列号为0
		s.sequence = 0
	}
	// 更新时间戳为当前时间
	s.timestamp = now
	// 生成64位整数ID
	r := (now-common.Twepoch)<<common.TimestampShift | (s.workerid << common.WorkeridShift) | (s.sequence)
	// 解锁，返回生成的ID
	s.Unlock()
	return r
}
