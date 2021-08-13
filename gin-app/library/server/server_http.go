package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	router Router
}

func NewHTTPServer(r Router) *HTTPServer {
	return &HTTPServer{
		router: r,
	}
}

func (s *HTTPServer) Run(app *Engine, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("http server end.")
		wg.Done()
	}()

	// 初始化gin server
	switch app.env {
	case "dev":
		gin.SetMode("debug")
	case "test":
		gin.SetMode("debug")
	case "prod":
		gin.SetMode("release")
	default:
		gin.SetMode("debug")
	}
	srv := gin.New()
	srv.Use(gin.Logger(), gin.Recovery())
	app.HTTPServer = srv

	// 设置路由
	s.router()

	// 启动服务
	address := ":8080"
	fmt.Printf("Listening and serving HTTP on %s\n", address)
	err := http.ListenAndServe(address, app.HTTPServer)
	if err != nil {
		fmt.Println("failed to serve HTTP:", err)
	}
}
