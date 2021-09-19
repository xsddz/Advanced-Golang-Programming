package app

import (
	"yawebapp/library/infra/server"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var defaultApp *server.Engine

// Init -
func Init() {
	// 1. 加载所有配置
	loadConf()

	// 2. 初始化全局日志
	initGlobalLogger()

	// 3. 资源检查
	resourcesCheck()

	// 4. 初始化默认app服务
	defaultApp = server.NewEngine()
}

func RegisterServer(srv server.ServerI) { defaultApp.RegisterServer(srv) }
func GetHTTPServer() *gin.Engine        { return defaultApp.GetHTTPServer() }
func GetGRPCServer() *grpc.Server       { return defaultApp.GetGRPCServer() }
func Run()                              { defaultApp.Run() }

////////////////////////////////////////////////////////////////////////////////

func NewHTTPServer(router server.Router) *server.HTTPServer {
	conf := server.HTTPConf{Env: ConfApp().AppEnv, Port: ConfApp().Server["http"].Port}
	return server.NewHTTPServer(conf, router)
}

func NewGRPCServer(router server.Router) *server.GRPCServer {
	conf := server.GRPCConf{Port: ConfApp().Server["grpc"].Port}
	return server.NewGRPCServer(conf, router)
}

////////////////////////////////////////////////////////////////////////////////

type AppEnv int

const (
	_ AppEnv = iota
	ENV_DEV
	ENV_TEST
	ENV_ON_TEST
	ENV_PROD
)

// Env 获取当前app的env
func Env() AppEnv {
	switch ConfApp().AppEnv {
	case "dev", "Dev", "DEV":
		return ENV_DEV
	case "test", "Test", "TEST":
		return ENV_TEST
	case "prod", "Prod", "PROD":
		return ENV_PROD
	default:
		return ENV_DEV
	}
}

// Name 获取当前app的name
func Name() string {
	return ConfApp().AppName
}
