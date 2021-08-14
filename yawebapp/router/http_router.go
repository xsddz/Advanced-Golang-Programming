package router

import (
	"yawebapp/Illuminate/app"
	"yawebapp/Illuminate/middleware"
	"yawebapp/controllers/http/demo"
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
