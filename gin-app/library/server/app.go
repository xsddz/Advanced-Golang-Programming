package server

import (
	"fmt"
	"gin-app/library/storage"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Router func(*Engine)

type ServerI interface {
	Run(*Engine, *sync.WaitGroup)
}

type Engine struct {
	HTTPServer *gin.Engine
	GRPCServer *grpc.Server

	serverCollector []ServerI
}

func NewAPP() *Engine {
	app := &Engine{}

	app.InitEnv()

	return app
}

func (app *Engine) InitEnv() {
	var err error
	SQLite, err = storage.NewSQLite()
	if err != nil {
		panic(fmt.Sprint("init sqlite error: ", err))
	}
	MySQL, err = storage.NewSQLite()
	if err != nil {
		panic(fmt.Sprint("init mysql error: ", err))
	}
}

func (app *Engine) RegisterServer(s ServerI) {
	app.serverCollector = append(app.serverCollector, s)
}

func (app *Engine) Run() {
	fmt.Println("run app ...")
	var wg sync.WaitGroup
	wg.Add(len(app.serverCollector))

	for _, sh := range app.serverCollector {
		go sh.Run(app, &wg)
	}

	wg.Wait()
	fmt.Println("the end.")
}
