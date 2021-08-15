package app

import (
	"fmt"
	"strconv"
	"strings"
	"yawebapp/library/inner/config"
	"yawebapp/library/inner/storage"

	"github.com/shima-park/agollo"
	"gorm.io/gorm"
)

func initAgollo() agollo.Agollo {
	var apolloConf config.ApolloConf
	config.LoadConf(ConfPath()+"/apollo.toml", &apolloConf)

	agollo, err := config.NewAgollo(apolloConf)
	if err != nil {
		panic(fmt.Sprint("init agollo error: ", err))
	}

	return agollo
}

func initSQLite(database string) *gorm.DB {
	sqlite, err := storage.NewSQLite(DataPath())
	if err != nil {
		panic(fmt.Sprint("init sqlite error: ", err))
	}

	return sqlite
}

func initMySQL(database string) *gorm.DB {
	conf := storage.MysqlConf{
		Host:                Apollo().Get("DB_HOST"),
		Port:                Apollo().Get("DB_PORT"),
		Username:            Apollo().Get("DB_USERNAME"),
		Password:            Apollo().Get("DB_PASSWORD"),
		Database:            Apollo().Get("DB_DATABASE"),
		Charset:             "utf8mb4",
		MaxIdleConns:        10,
		MaxOpenConns:        100,
		ConnMaxLifetimeHour: 1,
	}

	db, err := storage.NewMySQL(conf)
	if err != nil {
		panic(fmt.Sprint("init mysql error: ", err))
	}

	return db
}

func initFileCache() *storage.FileCache {
	fileCache, err := storage.NewFileCache(DataPath())
	if err != nil {
		panic(fmt.Sprint("init filecache error: ", err))
	}

	return fileCache
}

func initRedis() *storage.Redis {
	hosts := strings.Split(Apollo().Get("REDIS_SENTINEL_HOST"), ",")
	database, _ := strconv.Atoi(Apollo().Get("REDIS_SENTINEL_DB"))
	conf := storage.RedisConf{
		MasterName: Apollo().Get("REDIS_SENTINEL_SERVICE"),
		Hosts:      hosts,
		Passowrd:   Apollo().Get("REDIS_SENTINEL_PASSWORD"),
		DB:         database,
	}

	redis, err := storage.NewRedis(conf)
	if err != nil {
		panic(fmt.Sprint("init redis error: ", err))
	}

	return redis
}
