package idutils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	defaultAutoIncrID         uint64
	defaultAutoIncrIDInitOnce sync.Once
	defaultNextIDLock         sync.Mutex
)

// DefaultNextID -
//     理论上一秒内最多生成 1024 * 16384 个id，这里我们在uint64的数字空间上限下，
//     随机生成一个数字，然后将其映射过来。一般情况下，我们可以简单认为这是够用的。
func DefaultNextID() int64 {
	defaultAutoIncrIDInitOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
		defaultAutoIncrID = rand.Uint64() // 将自增id初始化为随机值，降低在代码分布式部署情况下，同一秒内产生重复id概率
	})

	defaultNextIDLock.Lock()
	shardID := defaultAutoIncrID / 16384
	id := IDGenerator(time.Now(), shardID, defaultAutoIncrID)
	defaultAutoIncrID++
	defaultNextIDLock.Unlock()
	return id
}

// IDGenerator -
// 64bit = 1bit + 39bit + 10bit + 14bit
//         39bit for time:
//             year 0-(9999-1970=8029) => 2^13 => 13bit
//             month 0-12              => 2^4  => 4bit
//             day 0-31                => 2^5  => 5bit
//             hour 0-23               => 2^5  => 5bit
//             minute 0-59             => 2^6  => 6bit
//             seconds 0-59            => 2^6  => 6bit
//         10bit for shard ID:
//             eg: uid % 1024
//                 orgid % 1024
//                 ...
//         14bit for auto-incrementing sequence:
//             eg: incrID % 16384
// thus, on each shard, we have 16384 ids in every seconds.
func IDGenerator(currTime time.Time, shardID uint64, incrID uint64) int64 {
	num := uint64(0)

	// add time
	t := uint64(currTime.Year() - 1970)
	t = (t << 4) + uint64(currTime.Month())
	t = (t << 5) + uint64(currTime.Day())
	t = (t << 5) + uint64(currTime.Hour())
	t = (t << 6) + uint64(currTime.Minute())
	t = (t << 6) + uint64(currTime.Second())
	num += (t << (10 + 14))

	// add shard ID
	num += ((shardID % 1024) << 14)

	// add auto-incrementing sequence
	num += (incrID % 16384)

	return int64(num)
}
