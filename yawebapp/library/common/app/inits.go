package app

import (
	"fmt"
	"strconv"
	"strings"
	"yawebapp/library/common/config"
	"yawebapp/library/common/storage"
)

var initTable = map[string]func(){
	"apollo": initApollo,
	"mysql":  initDB,
	"redis":  initRedis,
	"sqlite": initSQLite,
}

func initApollo() {
	var err error

	var apolloConf config.ApolloConf
	config.LoadConf(ConfPath()+"/apollo.toml", &apolloConf)

	APOLLO, err = config.NewAgollo(apolloConf)
	if err != nil {
		panic(fmt.Sprint("init agollo error: ", err))
	}
}

func initDB() {
	var err error

	conf := storage.MysqlConf{
		Host:                APOLLO.Get("DB_HOST"),
		Port:                APOLLO.Get("DB_PORT"),
		Username:            APOLLO.Get("DB_USERNAME"),
		Password:            APOLLO.Get("DB_PASSWORD"),
		Database:            APOLLO.Get("DB_DATABASE"),
		Charset:             "utf8mb4",
		MaxIdleConns:        10,
		MaxOpenConns:        100,
		ConnMaxLifetimeHour: 1,
	}

	DB, err = storage.NewMySQL(conf)
	if err != nil {
		panic(fmt.Sprint("init mysql error: ", err))
	}
}

func initRedis() {
	var err error

	hosts := strings.Split(APOLLO.Get("REDIS_SENTINEL_HOST"), ",")
	database, _ := strconv.Atoi(APOLLO.Get("REDIS_SENTINEL_DB"))
	conf := storage.RedisConf{
		MasterName: APOLLO.Get("REDIS_SENTINEL_SERVICE"),
		Hosts:      hosts,
		Passowrd:   APOLLO.Get("REDIS_SENTINEL_PASSWORD"),
		DB:         database,
	}

	REDIS, err = storage.NewRedis(conf)
	if err != nil {
		panic(fmt.Sprint("init mysql error: ", err))
	}
}

func initSQLite() {
	var err error
	SQLITE, err = storage.NewSQLite(DataPath())
	if err != nil {
		panic(fmt.Sprint("init sqlite error: ", err))
	}
}
