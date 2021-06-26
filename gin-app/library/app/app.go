package app

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type ServerHandle func(*sync.WaitGroup)

type Engine struct {
	servers []ServerHandle
}

func Init() *Engine {
	return &Engine{}
}

func (e *Engine) RegisterServer(s ServerHandle) {
	e.servers = append(e.servers, s)
}

func (e *Engine) Run() {
	var wg sync.WaitGroup
	wg.Add(len(e.servers))

	for _, sh := range e.servers {
		go sh(&wg)
	}

	wg.Wait()
	fmt.Println("the end.")
}

func RCPServerMaker(routerSetter func(*grpc.Server)) ServerHandle {
	return func(wg *sync.WaitGroup) {
		defer func() {
			fmt.Println("rpc server end.")
			wg.Done()
		}()

		// 初始化 RPC app
		var opts []grpc.ServerOption
		app := grpc.NewServer(opts...)

		// 设置路由
		routerSetter(app)

		// 启动服务
		address := ":8081"
		fmt.Printf("Listening and serving GRPC on %s\n", address)
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := app.Serve(lis); err != nil {
			log.Fatalf("failed to serve GRPC: %v", err)
		}
	}
}

func HTTPServerMaker(routerSetter func(app *gin.Engine)) ServerHandle {
	return func(wg *sync.WaitGroup) {
		defer func() {
			fmt.Println("http server end.")
			wg.Done()
		}()

		// 初始化gin app
		app := gin.New()
		app.Use(gin.Logger(), gin.Recovery())

		// 设置路由
		routerSetter(app)

		// 启动服务
		address := ":8080"
		fmt.Printf("Listening and serving HTTP on %s\n", address)
		err := http.ListenAndServe(address, app)
		if err != nil {
			fmt.Println("failed to serve HTTP:", err)
		}
	}
}
