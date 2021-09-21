package server

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Router 可注册服务的路由方法约定
type Router func()

// ServerI 可注册服务的方法约定
type ServerI interface {
	// RegisterRouter 路由注册
	RegisterRouter(...Router)
	// Run 服务启动
	Run(*sync.WaitGroup)
}

////////////////////////////////////////////////////////////////////////////////

// Engine -
type Engine struct {
	servers     []string
	serverHoder map[string]ServerI
}

// NewEngine -
func NewEngine() *Engine {
	return &Engine{
		serverHoder: make(map[string]ServerI),
	}
}

// RegisterServer -
func (e *Engine) RegisterServer(name string, server ServerI) {
	if _, ok := e.serverHoder[name]; ok {
		return
	}

	e.servers = append(e.servers, name)
	e.serverHoder[name] = server
}

// HTTPServer -
func (e *Engine) HTTPServer() *gin.Engine {
	if server, ok := e.serverHoder["http"]; ok {
		return server.(*HTTPServer).Engine
	}
	return nil
}

// GRPCServer -
func (e *Engine) GRPCServer() *grpc.Server {
	if server, ok := e.serverHoder["grpc"]; ok {
		return server.(*GRPCServer).Server
	}
	return nil
}

// AddServerRouter -
func (e *Engine) AddServerRouter(name string, routers ...Router) error {
	if server, ok := e.serverHoder[name]; ok {
		server.RegisterRouter(routers...)
		return nil
	}

	return fmt.Errorf("must register %v server first", name)
}

// Run -
func (e *Engine) Run() {
	fmt.Println("\nrun app ...")
	var wg sync.WaitGroup
	wg.Add(len(e.servers))

	for _, name := range e.servers {
		go e.serverHoder[name].Run(&wg)
	}

	wg.Wait()
	fmt.Println("the end.")
}
