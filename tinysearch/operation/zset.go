package operation

import (
	"log"
	"tinysearch/util"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

// ZSetCommon 集合操作公共方法
func ZSetCommon(rconn redis.Conn, op string, keyWeight map[string]int, ttl int, execute bool) string {
	cacheKey := util.MakeOpCacheRedisKey(uuid.Must(uuid.NewRandom()).String())

	// 定义为interface{}类型，在调用redis.Conn的Send和Do等方法时，与参数类型保持一致，否则会报错
	keys := []interface{}{cacheKey, len(keyWeight)}
	weights := []interface{}{"WEIGHTS"}
	for key, weight := range keyWeight {
		keys = append(keys, key)
		weights = append(weights, weight)
	}
	keys = append(keys, weights...)

	// fmt.Println(op, keys)
	if execute {
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

// ZSetIntersect 交集
func ZSetIntersect(conn redis.Conn, keyWeight map[string]int, ttl int, execute bool) string {
	return ZSetCommon(conn, "zinterstore", keyWeight, ttl, execute)
}

// ZSetUnion 并集
func ZSetUnion(conn redis.Conn, keyWeight map[string]int, ttl int, execute bool) string {
	return ZSetCommon(conn, "zunionstore", keyWeight, ttl, execute)
}
