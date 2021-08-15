package app

import (
	"yawebapp/library/inner/storage"
)

var (
	cacheDriver = "file"
	cacheTable  = make(map[string]storage.Cacher)
)

func initCacheDriver(driver string) {
	if driver != "redis" && driver != "file" {
		driver = "file"
	}
	dbDriver = driver
}

// cache 通过storage.Cacher接口提供统一的缓存操作方法，支持redis、文件等缓存
type cache struct {
	driver    string
	redis     *storage.Redis
	fileCache *storage.FileCache
}

func Cache() storage.Cacher {
	if c, ok := cacheTable[cacheDriver]; ok {
		return c
	}

	var c *cache
	switch cacheDriver {
	case "file":
		c = &cache{
			driver:    cacheDriver,
			fileCache: initFileCache(),
		}
	case "redis":
		c = &cache{
			driver: cacheDriver,
			redis:  initRedis(),
		}
	}
	cacheTable[cacheDriver] = c

	return cacheTable[cacheDriver]
}
