package app

import (
	"context"
	"fmt"
	"path/filepath"
	"yawebapp/library/inner/server"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	defaultApp *server.Engine
)

func Name() string {
	if appConf == nil {
		appConf.AppName = filepath.Base(RootPath())
	}
	return appConf.AppName
}

func Init() {
	loadConf()

	loadLogger()

	initDBDriver(appConf.DBDriver)
	// 提前验证一次默认数据库资源连接
	if err := DB(context.TODO()).Ping(); err != nil {
		panic(fmt.Sprint("ping default db error:", err))
	}

	initCacheDriver(appConf.CacheDriver)
	// 提前验证一次默认缓存资源连接
	if _, err := Cache().Ping().Result(); err != nil {
		panic(fmt.Sprint("ping default cache error:", err))
	}

	defaultApp = server.NewEngine(appConf.AppEnv)
}

func RegisterServer(s server.ServerI) {
	defaultApp.RegisterServer(s)
}

func GetHTTPServer() *gin.Engine {
	return defaultApp.GetHTTPServer()
}

func GetGRPCServer() *grpc.Server {
	return defaultApp.GetGRPCServer()
}

func Run() {
	defaultApp.Run()
}
