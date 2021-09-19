package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"yawebapp/library/infra/storage"

	"github.com/go-redis/redis"
)

var cacheTable = make(map[string]*cache)

// cache
type cache struct {
	r []*redis.Client
	l *storage.LocalCache
}

func (c *cache) Ping() error                                                { return nil }
func (c *cache) Get(key string) (string, error)                             { return "", nil }
func (c *cache) Set(key string, val string, expiration time.Duration) error { return nil }

// Cache 提供延迟初始化对应资源连接的能力，及连接复用能力
func Cache(clusterNames ...string) *cache {
	clusterName := "Default"
	if len(clusterNames) > 0 {
		clusterName = clusterNames[0]
	}

	if _, ok := cacheTable[clusterName]; !ok {
		var c cache
		switch ConfCache(clusterName).Driver {
		case "local":
			c.l = initLocalCache(clusterName)
		case "redis":
			c.r = initRedis(clusterName)
		default:
			j, _ := json.Marshal(ConfCache(clusterName))
			panic(fmt.Errorf("unsupport cache driver in conf: %s", j))
		}
		cacheTable[clusterName] = &c
	}

	return cacheTable[clusterName]
}

func Redis(clusterNames ...string) *redis.Client {
	clusterName := "Default"
	if len(clusterNames) > 0 {
		clusterName = clusterNames[0]
	}

	if _, ok := cacheTable[clusterName]; !ok {
		var c cache
		switch ConfCache(clusterName).Driver {
		case "redis":
			c.r = initRedis(clusterName)
		default:
			j, _ := json.Marshal(ConfCache(clusterName))
			panic(fmt.Errorf("need redis driver in conf: %s", j))
		}
		cacheTable[clusterName] = &c
	}

	return randRedis(cacheTable[clusterName].r)
}

func randRedis(clis []*redis.Client) *redis.Client {
	rand.Seed(time.Now().UnixNano())
	return clis[rand.Intn(len(clis))]
}
