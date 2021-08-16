package app

import (
	"fmt"
	"strconv"
	"strings"
	"yawebapp/library/inner/config"
	"yawebapp/library/inner/storage"
)

type AppConf struct {
	UseApollo   bool   `toml:"use_apollo"`
	AppName     string `toml:"app_name"`
	AppEnv      string `toml:"app_env"`
	DBDriver    string `toml:"db_driver"`
	CacheDriver string `toml:"cache_driver"`
}

var (
	appConf      *AppConf
	appDBConf    = make(map[string]*storage.DBConf)
	appCacheConf = make(map[string]*storage.RedisConf)
)

func loadAppConf() {
	var c AppConf
	err := config.LoadConf(ConfPath()+"/app.toml", &c)
	if err != nil {
		panic(fmt.Sprint("init app conf error: ", err))
	}
	appConf = &c

	if appConf.UseApollo {
		appConf.AppName = Apollo().Get("APP_NAME")
		appConf.AppEnv = Apollo().Get("APP_ENV")
		appConf.DBDriver = Apollo().Get("DB_DRIVER")
		appConf.CacheDriver = Apollo().Get("CACHE_DRIVER")
	}

	loadDBConf()
	loadCacheConf()
}

func loadDBConf() {
	err := config.LoadConf(ConfPath()+"/db.toml", &appDBConf)
	if err != nil {
		panic(fmt.Sprint("init db conf error: ", err))
	}

	if appConf.UseApollo {
		port, _ := strconv.Atoi(Apollo().Get("DB_PORT"))
		defaultConf := &storage.DBConf{
			ClusterName:       "Default",
			MaxOpenConns:      24,
			MaxIdleConns:      24,
			ConnMaxLifetimeMs: 18000,
			BalanceStrategy:   "",
			Charset:           "utf8mb4",
			Username:          Apollo().Get("DB_USERNAME"),
			Password:          Apollo().Get("DB_PASSWORD"),
			DefaultDB:         Apollo().Get("DB_DATABASE"),
			Hosts: []storage.DBConfHost{
				{
					IP:   Apollo().Get("DB_HOST"),
					Port: port,
				},
			},
		}

		appDBConf = map[string]*storage.DBConf{
			"Default": defaultConf,
		}
	}
}

func loadCacheConf() {
	err := config.LoadConf(ConfPath()+"/cache.toml", &appCacheConf)
	if err != nil {
		panic(fmt.Sprint("init cache conf error: ", err))
	}

	if appConf.UseApollo {
		defaultDB, _ := strconv.Atoi(Apollo().Get("REDIS_SENTINEL_DB"))
		hosts := []storage.RedisConfHost{}
		for _, host := range strings.Split(Apollo().Get("REDIS_SENTINEL_HOST"), ",") {
			hp := strings.Split(host, ":")
			p, _ := strconv.Atoi(hp[1])
			hosts = append(hosts, storage.RedisConfHost{IP: hp[0], Port: p})
		}
		defaultConf := &storage.RedisConf{
			MasterName: Apollo().Get("REDIS_SENTINEL_SERVICE"),
			Passowrd:   Apollo().Get("REDIS_SENTINEL_PASSWORD"),
			DefaultDB:  defaultDB,
			Hosts:      hosts,
		}

		appCacheConf = map[string]*storage.RedisConf{
			"Default": defaultConf,
		}
	}
}
