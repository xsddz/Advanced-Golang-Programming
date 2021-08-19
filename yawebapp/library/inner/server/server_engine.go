package server

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Router func()

type ServerI interface {
	Run(*Engine, *sync.WaitGroup)
}

type Engine struct {
	env string

	HTTPServer *gin.Engine
	GRPCServer *grpc.Server

	serverCollector []ServerI
}

func NewEngine(env string) *Engine {
	return &Engine{env: env}
}

func (app *Engine) RegisterServer(s ServerI) {
	app.serverCollector = append(app.serverCollector, s)
}

func (app *Engine) GetHTTPServer() *gin.Engine {
	return app.HTTPServer
}

func (app *Engine) GetGRPCServer() *grpc.Server {
	return app.GRPCServer
}

func (app *Engine) Run() {
	fmt.Println("run app ...")
	var wg sync.WaitGroup
	wg.Add(len(app.serverCollector))

	for _, srv := range app.serverCollector {
		go srv.Run(app, &wg)
	}

	wg.Wait()
	fmt.Println("the end.")
}
