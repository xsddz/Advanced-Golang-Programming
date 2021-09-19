package app

import (
	"fmt"
	"yawebapp/library/inner/config"
	"yawebapp/library/inner/storage"
	"yawebapp/library/inner/trace"

	"github.com/go-redis/redis"
	"github.com/shima-park/agollo"
	"gorm.io/gorm"
)

func initLogger() *trace.Logger {
	return trace.NewLogger(trace.LogConf{
		FilePrefix: LogPath() + "/" + Name(),
		Rotate:     1,
		KeepDays:   7,
	})
}

func initAgollo() agollo.Agollo {
	var apolloConf config.ApolloConf
	config.LoadConf(ConfPath()+"/apollo.toml", &apolloConf)

	agollo, err := config.NewAgollo(apolloConf)
	if err != nil {
		panic(fmt.Sprint("init agollo error: ", err))
	}

	return agollo
}

func initSQLite(clusterName string) []*gorm.DB {
	if _, ok := appDBConf[clusterName]; !ok {
		panic(fmt.Sprint("[initSQLite] cluster conf not exist:", clusterName))
	}

	dbs, err := storage.NewSQLite(DataPath(), *appDBConf[clusterName], Logger)
	if err != nil {
		panic(fmt.Sprint("[initSQLite] init sqlite error: ", err))
	}

	return dbs
}

func initMySQL(clusterName string) []*gorm.DB {
	if _, ok := appDBConf[clusterName]; !ok {
		panic(fmt.Sprint("[initMySQL] cluster conf not exist:", clusterName))
	}

	dbs, err := storage.NewMySQL(*appDBConf[clusterName], Logger)
	if err != nil {
		panic(fmt.Sprint("[initMySQL] init mysql error: ", err))
	}

	return dbs
}

func initFileCache() []*redis.Client {
	if _, ok := appCacheConf["Default"]; !ok {
		panic(fmt.Sprint("[initFileCache] cluster conf not exist:", "Default"))
	}

	caches, err := storage.NewFileCache(DataPath(), *appCacheConf["Default"])
	if err != nil {
		panic(fmt.Sprint("[initFileCache] init filecache error: ", err))
	}

	return caches
}

func initRedis() []*redis.Client {
	if _, ok := appCacheConf["Default"]; !ok {
		panic(fmt.Sprint("[initRedis] cluster conf not exist:", "Default"))
	}

	var caches []*redis.Client
	var err error
	if appConf.AppEnv == "production" {
		caches, err = storage.NewRedisSentine(*appCacheConf["Default"])
	} else {
		caches, err = storage.NewRedis(*appCacheConf["Default"])
	}
	if err != nil {
		panic(fmt.Sprint("[initRedis] init redis error: ", err))
	}

	return caches
}
