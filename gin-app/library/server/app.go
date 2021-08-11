package server

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type ServerI interface {
	Run(*App, *sync.WaitGroup)
}

type WebContext struct {
	*gin.Context
}

type App struct {
	HTTPServer *gin.Engine
	GRPCServer *grpc.Server

	serverCollector []ServerI
}

func Init() *App {
	return &App{}
}

func (a *App) RegisterServer(s ServerI) {
	a.serverCollector = append(a.serverCollector, s)
}

func (a *App) Run() {
	fmt.Println("run app ...")
	var wg sync.WaitGroup
	wg.Add(len(a.serverCollector))

	for _, sh := range a.serverCollector {
		go sh.Run(a, &wg)
	}

	wg.Wait()
	fmt.Println("the end.")
}
