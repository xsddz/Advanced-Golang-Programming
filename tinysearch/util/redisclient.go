package util

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

const (
	// address redis连接地址
	address = "10.14.15.214:8379"

	// tplOpCache 集合操作结果缓存key模版
	tplOpCache = "cache:set:%s"
	// tplDocWordsIdxKey 文档单词集索引key模版
	tplDocWordsIdxKey = "idx:string:%s"
	// tplDocWrodsRIdxKey 文档单词倒排索引key模版
	tplDocWrodsRIdxKey = "ridx:set:%s"
	// tplDocUpdateKey 文档更新时间排序集合key模版
	tplDocUpdateKey = "doc:upadte:zset"
	// tplDocVoteKey 文档票数排序集合key模版
	tplDocVoteKey = "doc:vote:zset"
)

// MakeOpCacheRedisKey 生成指定名称对应的缓存key
func MakeOpCacheRedisKey(name string) string {
	return fmt.Sprintf(tplOpCache, name)
}

// MakeDocWordsIndexRedisKey 生成指定name对应的索引key
func MakeDocWordsIndexRedisKey(name string) string {
	return fmt.Sprintf(tplDocWordsIdxKey, name)
}

// MakeDocWordsRIndexRedisKey 生成指定单词的倒排索引key
func MakeDocWordsRIndexRedisKey(word string) string {
	keys := MakeDocWordsRIndexRedisKeys([]string{word})
	return keys[0]
}

// MakeDocWordsRIndexRedisKeys 生成指定单词列表的倒排索引key列表
func MakeDocWordsRIndexRedisKeys(words []string) (keys []string) {
	for _, word := range words {
		keys = append(keys, fmt.Sprintf(tplDocWrodsRIdxKey, word))
	}
	return
}

// MakeDocUpdateRedisKey 生成文档更新时间排序集合key
func MakeDocUpdateRedisKey() string {
	return tplDocUpdateKey
}

// MakeDocVoteRedisKey 生成文档票数排序集合key
func MakeDocVoteRedisKey() string {
	return tplDocVoteKey
}

// NewRedisConn 新建redis客户端
func NewRedisConn() redis.Conn {
	rconn, err := redis.Dial("tcp", address)
	if err != nil {
		log.Fatal("Connect to redis error", err)
	}

	return rconn
}
