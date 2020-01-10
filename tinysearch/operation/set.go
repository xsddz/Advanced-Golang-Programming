package operation

import (
	"log"
	"tinysearch/util"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

// SetCommon 集合操作公共方法
func SetCommon(rconn redis.Conn, op string, names []string, ttl int, execute bool) string {
	cacheKey := util.MakeOpCacheRedisKey(uuid.Must(uuid.NewRandom()).String())

	// 定义为interface{}类型，在调用redis.Conn的Send和Do等方法时，与参数类型保持一致，否则会报错
	keys := []interface{}{cacheKey}
	for _, key := range names {
		keys = append(keys, key)
	}

	// fmt.Println(op, keys)
	if execute {
		// rconn.Send(op, cacheKey, "ridx:set:满意", "ridx:set:一种")
		rconn.Send(op, keys...)
		rconn.Send("EXPIRE", cacheKey, ttl)
		rconn.Flush()
		_, err := rconn.Receive()
		if err != nil {
			log.Fatal("Redis pipeline error", err)
		}
	} else {
		_, err := rconn.Do(op, keys)
		if err != nil {
			log.Fatal("Redis op error", err)
		}
		_, err = rconn.Do("EXPIRE", cacheKey, ttl)
		if err != nil {
			log.Fatal("Redis op error", err)
		}
	}

	return cacheKey
}

// SetIntersect 交集
func SetIntersect(conn redis.Conn, names []string, ttl int, execute bool) string {
	return SetCommon(conn, "sinterstore", names, ttl, execute)
}

// SetUnion 并集
func SetUnion(conn redis.Conn, names []string, ttl int, execute bool) string {
	return SetCommon(conn, "sunionstore", names, ttl, execute)
}

// SetDifference 差集
func SetDifference(conn redis.Conn, names []string, ttl int, execute bool) string {
	return SetCommon(conn, "sdiffstore", names, ttl, execute)
}
