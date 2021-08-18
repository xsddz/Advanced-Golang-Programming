package app

import (
	"fmt"
	"yawebapp/library/inner/server"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	defaultApp *server.Engine
)

func Name() string {
	if appConf == nil {
		return ""
	}
	return appConf.AppName
}

func Init() {
	loadAppConf()

	initDBDriver(appConf.DBDriver)
	initCacheDriver(appConf.CacheDriver)

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
	// 提前验证一次默认资源
	if err := DB().Ping(); err != nil {
		panic(fmt.Sprint("ping default db error:", err))
	}
	if _, err := Cache().Ping().Result(); err != nil {
		panic(fmt.Sprint("ping default cache error:", err))
	}

	defaultApp.Run()
}
