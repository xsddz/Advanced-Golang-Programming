package main

import (
	"fmt"
	"math/rand"
	"time"
	"tinysearch/document"
	"tinysearch/operation"
	"tinysearch/util"

	"github.com/gomodule/redigo/redis"
)

func main() {
	rconn := util.NewRedisConn()
	defer rconn.Close()

	files := []string{
		"data/是的，我回来了.txt",
		"data/动机三和银杏叶落时分前传.txt",
		"data/动机三和银杏叶落时分.txt",
	}
	rand.Seed(time.Now().UnixNano())
	for _, f := range files {
		doc := document.ReadDoc(f)

		// 更新文档信息
		doc.VoteDoc(rconn, rand.Intn(10))
		doc.UpdateDocUpdateTime(rconn, time.Now().Unix()+rand.Int63n(1000))

		// 索引文档
		// indexing.IndexDocument(rconn, doc)
	}
	// 查看文档的更新信息
	v, err := redis.Strings(rconn.Do("ZRANGE", util.MakeDocVoteRedisKey(), 0, -1, "withscores"))
	fmt.Println(util.MakeDocVoteRedisKey(), v, err)
	v, err = redis.Strings(rconn.Do("ZRANGE", util.MakeDocUpdateRedisKey(), 0, -1, "withscores"))
	fmt.Println(util.MakeDocUpdateRedisKey(), v, err)

	// 搜索操作
	rk := operation.Search(rconn, "一种 一些", 30)
	v, err = redis.Strings(rconn.Do("SMEMBERS", rk))
	fmt.Println(rk, v, err)
	// 搜索分页操作
	count, items, cacheKey := operation.SearchAndZSort(rconn, "一种 一些", "", 30, 1, 0, 0, 1, true)
	fmt.Println(count, items, cacheKey)
	count, items, cacheKey = operation.SearchAndZSort(rconn, "一种 一些", cacheKey, 30, 1, 0, len(items), 20, true)
	fmt.Println(count, items, cacheKey)

	// setOp(rconn)
	// zsetOp(rconn)
}

// setOp 集合操作测试
func setOp(rconn redis.Conn) {
	fmt.Println("=================== set operation ===================")

	// items := []string{"满意", "一种", "一些"}
	items := []string{"一种", "一些"}
	// 查看倒排索引信息
	for _, word := range items {
		v, err := redis.Strings(rconn.Do("SMEMBERS", util.MakeDocWordsRIndexRedisKey(word)))
		fmt.Println(word, v, err)
	}

	operations := []func(redis.Conn, []string, int, bool) string{
		operation.SetIntersect,
		operation.SetUnion,
		operation.SetDifference,
	}
	for _, operation := range operations {
		rk := operation(rconn, util.MakeDocWordsRIndexRedisKeys(items), 120, true)
		v, err := redis.Strings(rconn.Do("SMEMBERS", rk))
		fmt.Println(rk, v, err)
	}

	fmt.Println("=================== end set operation ===================")
}

// zsetOp 排序集合操作测试
func zsetOp(rconn redis.Conn) {
	fmt.Println("=================== zset operation ===================")

	zsetkeys := []string{"zset1", "zset2"}

	for _, key := range zsetkeys {
		v, err := redis.Strings(rconn.Do("ZRANGE", key, 0, -1, "withscores"))
		fmt.Println("zset1", v, err)
	}

	names := map[string]int{
		"zset1": 2,
		"zset2": 3,
	}
	operations := []func(redis.Conn, map[string]int, int, bool) string{
		operation.ZSetIntersect,
		operation.ZSetUnion,
	}
	for _, operation := range operations {
		rk := operation(rconn, names, 120, true)
		v, err := redis.Strings(rconn.Do("ZRANGE", rk, 0, -1, "withscores"))
		fmt.Println(rk, v, err)
	}

	fmt.Println("=================== end zset operation ===================")
}
