package app

import (
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

var (
	cacheDriver = "file"
	cacheTable  = make(map[string][]*redis.Client)
)

func initCacheDriver(driver string) {
	if driver != "redis" && driver != "file" {
		driver = "file"
	}
	cacheDriver = driver
}

// cache 通过嵌套redis.Client，利用 redis.Client 接口提供统一的缓存操作方法，支持redis、文件等缓存
type cache struct {
	*redis.Client
}

// Cache 提供延迟初始化对应资源连接的能力，及连接复用能力
func Cache() *cache {
	if _, ok := cacheTable[cacheDriver]; !ok {
		var c []*redis.Client
		switch cacheDriver {
		case "file":
			c = initFileCache()
		case "redis":
			c = initRedis()
		}
		cacheTable[cacheDriver] = c
	}

	return randCache(cacheTable[cacheDriver])
}

func randCache(caches []*redis.Client) *cache {
	rand.Seed(time.Now().UnixNano())
	return &cache{caches[rand.Intn(len(caches))]}
}
