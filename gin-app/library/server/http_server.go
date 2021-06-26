package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

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
