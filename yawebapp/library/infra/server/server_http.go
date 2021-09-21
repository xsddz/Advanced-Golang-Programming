package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// HTTPConf -
type HTTPConf struct {
	Env  string `toml:"env"`
	Port string `toml:"port"`
}

// HTTPServer -
type HTTPServer struct {
	*gin.Engine

	address string
	routers []Router
}

// NewHTTPServer -
func NewHTTPServer(conf HTTPConf, routers ...Router) *HTTPServer {
	// gin模式设置
	switch conf.Env {
	case "prod", "Prod", "PROD":
		gin.SetMode("release")
	default:
		gin.SetMode("debug")
		// 开启终端颜色输出
		gin.ForceConsoleColor()
	}

	address := fmt.Sprintf(":%v", conf.Port)
	return &HTTPServer{gin.New(), address, routers}
}

// RegisterRouter 实现ServerI接口
func (s *HTTPServer) RegisterRouter(routers ...Router) {
	s.routers = append(s.routers, routers...)
}

// Run 实现ServerI接口
func (s *HTTPServer) Run(wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("http server end.")
		wg.Done()
	}()

	// gin路由设置
	for _, router := range s.routers {
		router()
	}

	// gin服务启动
	fmt.Printf("Listening and serving HTTP on %s\n", s.address)
	err := http.ListenAndServe(s.address, s)
	if err != nil {
		panic(fmt.Errorf("failed to serve HTTP: %v", err))
	}
}
