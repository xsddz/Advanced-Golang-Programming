package router

import (
	"gin-app/controllers/http/demo"
	"gin-app/library/app"
	"gin-app/library/middleware"
)

func HTTPRouter() {
	httpServer := app.GetHTTPServer()
	httpServer.NoRoute(middleware.RouterNotFound)

	// 路由组
	DemoGr := httpServer.Group("/demo")
	{
		DemoGr.GET("/gituser", demo.GitUser)
	}
}
