package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	timestampBits  = uint(41) // 时间戳占用位数
	machineID      = uint(10)
	sequenceBits   = uint(12)                          // 序列所占的位数
	timestampMax   = int64(-1 ^ (-1 << timestampBits)) // 时间戳最大值
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits))  // 支持的最大序列id数量
	timestampShift = sequenceBits + machineID          // 时间戳左移位数
)

type Snowflake struct {
	sync.Mutex       // 锁
	timestamp  int64 // 时间戳 ，毫秒
	sequence   int64 // 序列号
	epoch      int64
}

type option struct {
	epoch int64
}

type WithOptionFunc func(op *option)

func New(fns ...WithOptionFunc) *Snowflake {

	defOption := &option{
		epoch: 1577808000000, // default timestamp, Millisecond, 2020-01-01 00:00:00, expire at 69 years later
	}

	for _, fn := range fns {
		fn(defOption)
	}

	return &Snowflake{
		epoch: defOption.epoch,
	}
}

func WithEpoch(epoch int64) WithOptionFunc {
	return func(op *option) {
		op.epoch = epoch
	}
}

func (s *Snowflake) NextVal() (int64, error) {
	s.Lock()
	defer s.Unlock()
	now := time.Now().UnixNano() / 1000000 // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - s.epoch
	if t > timestampMax {
		return 0, fmt.Errorf("epoch must be between 0 and %d", timestampMax-1)
	}
	s.timestamp = now
	r := (t)<<timestampShift | (s.sequence)
	return r, nil
}
