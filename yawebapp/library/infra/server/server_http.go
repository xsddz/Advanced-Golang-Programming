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
	conf   *HTTPConf
	router Router
}

// NewHTTPServer -
func NewHTTPServer(conf HTTPConf, router Router) *HTTPServer {
	return &HTTPServer{&conf, router}
}

// Run 实现ServerI接口
func (s *HTTPServer) Run(app *Engine, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("http server end.")
		wg.Done()
	}()

	// gin模式设置
	switch s.conf.Env {
	case "prod", "Prod", "PROD":
		gin.SetMode("release")
	default:
		gin.SetMode("debug")
		// 开启终端颜色输出
		gin.ForceConsoleColor()
	}

	// gin server初始化
	srv := gin.New()
	app.SetHTTPServer(srv)

	// gin路由设置
	s.router()

	// gin服务启动
	address := fmt.Sprintf(":%v", s.conf.Port)
	fmt.Printf("Listening and serving HTTP on %s\n", address)
	err := http.ListenAndServe(address, app.GetHTTPServer())
	if err != nil {
		panic(fmt.Errorf("failed to serve HTTP: %v", err))
	}
}
