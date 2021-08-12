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
	srv := gin.New()
	srv.Use(gin.Logger(), gin.Recovery())

	// 设置路由
	app.HTTPServer = srv
	s.router(app)

	// 启动服务
	address := ":8080"
	fmt.Printf("Listening and serving HTTP on %s\n", address)
	err := http.ListenAndServe(address, srv)
	if err != nil {
		fmt.Println("failed to serve HTTP:", err)
	}
}
