package app

import (
	"yawebapp/library/inner/storage"
)

var (
	cacheDriver = "file"
	cacheTable  = make(map[string]*cache)
)

func initCacheDriver(driver string) {
	if driver != "redis" && driver != "file" {
		driver = "file"
	}
	dbDriver = driver
}

// cache 通过 storage.Cacher 接口提供统一的缓存操作方法，支持redis、文件等缓存
type cache struct {
	storage.Cacher
}

// Cache 提供延迟初始化对应资源连接的能力，及连接复用能力
func Cache() *cache {
	if c, ok := cacheTable[cacheDriver]; ok {
		return c
	}

	var c *cache
	switch cacheDriver {
	case "file":
		c = &cache{initFileCache()}
	case "redis":
		c = &cache{initRedis()}
	}
	cacheTable[cacheDriver] = c

	return cacheTable[cacheDriver]
}
