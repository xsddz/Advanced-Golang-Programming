package server

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Router 服务的路由方法约定
type Router func()

// ServerI 服务标准
type ServerI interface {
	// Run 运行服务
	Run(*Engine, *sync.WaitGroup)
}

type Engine struct {
	serverHoder []ServerI

	httpServer *gin.Engine
	grpcServer *grpc.Server
}

func NewEngine() *Engine                           { return &Engine{} }
func (app *Engine) RegisterServer(srv ServerI)     { app.serverHoder = append(app.serverHoder, srv) }
func (app *Engine) GetHTTPServer() *gin.Engine     { return app.httpServer }
func (app *Engine) SetHTTPServer(srv *gin.Engine)  { app.httpServer = srv }
func (app *Engine) GetGRPCServer() *grpc.Server    { return app.grpcServer }
func (app *Engine) SetGRPCServer(srv *grpc.Server) { app.grpcServer = srv }

func (app *Engine) Run() {
	fmt.Println("run app ...")
	var wg sync.WaitGroup
	wg.Add(len(app.serverHoder))

	for _, srv := range app.serverHoder {
		go srv.Run(app, &wg)
	}

	wg.Wait()
	fmt.Println("the end.")
}
