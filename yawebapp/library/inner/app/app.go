package app

import (
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
	defaultApp.Run()
}
