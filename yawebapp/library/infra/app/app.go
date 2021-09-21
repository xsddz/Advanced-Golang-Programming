package app

import (
	"fmt"
	"strconv"
	"yawebapp/library/infra/server"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Init 初始化app
func Init() {
	// 1. 加载所有配置
	loadConf()

	// 2. 初始化全局日志
	initGlobalLogger()

	// 3. 资源检查
	resourcesCheck()

	// 4. 注册服务
	registerServer()
}

////////////////////////////////////////////////////////////////////////////////

var defaultEngine *server.Engine

func registerServer() {
	defaultEngine = server.NewEngine()

	for name, conf := range ConfApp().Server {
		fmt.Printf("register server: [%v] \t......", name)
		if ok, err := strconv.ParseBool(conf.Enable); err != nil || !ok {
			fmt.Printf("... [SKIP]\n")
			continue
		}

		switch name {
		case "http":
			httpconf := server.HTTPConf{Env: ConfApp().AppEnv, Port: conf.Port}
			defaultEngine.RegisterServer("http", server.NewHTTPServer(httpconf))
		case "grpc":
			grpcconf := server.GRPCConf{Port: conf.Port}
			defaultEngine.RegisterServer("grpc", server.NewGRPCServer(grpcconf))
		default:
			panic(fmt.Errorf("unsupport server name: %v", name))
		}
		fmt.Printf("... [OK]\n")
	}
}

// HTTPServer 获取http服务实例
func HTTPServer() *gin.Engine {
	if srv := defaultEngine.HTTPServer(); srv != nil {
		return srv
	}

	panic(fmt.Errorf("must register http server first"))
}

// GRPCServer 获取grpc服务实例
func GRPCServer() *grpc.Server {
	if srv := defaultEngine.GRPCServer(); srv != nil {
		return srv
	}

	panic(fmt.Errorf("must register grpc server first"))
}

// AddServerRouter 添加服务路由
func AddServerRouter(name string, router ...server.Router) {
	if err := defaultEngine.AddServerRouter(name, router...); err != nil {
		panic(err)
	}
}

// Run 启动服务
func Run() {
	defaultEngine.Run()
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
func Name() string { return ConfApp().AppName }
