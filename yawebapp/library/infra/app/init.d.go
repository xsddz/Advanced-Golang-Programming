package app

import (
	"fmt"
	"os"
	"yawebapp/library/infra/config"
	"yawebapp/library/infra/storage"
	"yawebapp/library/infra/trace"

	"github.com/go-redis/redis"
	"github.com/shima-park/agollo"
	"gorm.io/gorm"
)

func initLogger() *trace.Logger {
	loglevel := trace.LOG_INFO
	if Env() == ENV_DEV {
		loglevel = trace.LOG_ALL
	}

	conf := trace.LogConf{Level: loglevel, HasColor: OutputColor(), LogPath: LogPath(), AppName: Name()}
	switch loglevel {
	case trace.LOG_ALL, trace.LOG_TRACE, trace.LOG_DEBUG:
		conf.HasFileNum = true
	}

	log, err := trace.NewLogger(conf)
	if err != nil {
		panic(fmt.Sprint("[initLogger] init logger error: ", err))
	}

	if Env() == ENV_DEV {
		log.AddWriter(os.Stdout)
	}

	return log
}

func initAgollo() agollo.Agollo {
	agollo, err := config.NewAgollo(ConfApollo())
	if err != nil {
		panic(fmt.Sprint("[initAgollo] init agollo error: ", err))
	}

	return agollo
}

func initSQLite(clusterName string) []*gorm.DB {
	dbs, err := storage.NewSQLite(DataPath(), ConfDB(clusterName), Logger)
	if err != nil {
		panic(fmt.Sprint("[initSQLite] init sqlite error: ", err))
	}

	return dbs
}

func initMySQL(clusterName string) []*gorm.DB {
	dbs, err := storage.NewMySQL(ConfDB(clusterName), Logger)
	if err != nil {
		panic(fmt.Sprint("[initMySQL] init mysql error: ", err))
	}

	return dbs
}

func initLocalCache(clusterName string) *storage.LocalCache {
	cache, err := storage.NewLocalCache(ConfCache(clusterName), DataPath())
	if err != nil {
		panic(fmt.Sprint("[initLocalCache] init local cache error: ", err))
	}

	return cache
}

func initRedis(clusterName string) (clis []*redis.Client) {
	// 设置redis log
	storage.SetRedisLog(trace.NewRedisLogger(initGlobalLogger()))

	switch ConfCache(clusterName).Type {
	case "sentine", "Sentine", "SENTINE":
		clis, err := storage.NewRedisSentine(ConfCache(clusterName))
		if err != nil {
			panic(fmt.Sprint("[initRedis] init redis sentine error: ", err))
		}
		return clis
	default:
		clis, err := storage.NewRedis(ConfCache(clusterName))
		if err != nil {
			panic(fmt.Sprint("[initRedis] init redis error: ", err))
		}
		return clis
	}
}
