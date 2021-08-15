package app

import (
	"yawebapp/library/inner/server"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	defaultApp *server.Engine
)

func Init() {
	initDBDriver(Apollo().Get("DB_DRIVER"))
	initCacheDriver(Apollo().Get("CACHE_DRIVER"))

	defaultApp = server.NewEngine(Apollo().Get("APP_ENV"))
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
