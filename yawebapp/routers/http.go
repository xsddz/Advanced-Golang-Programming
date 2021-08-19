package routers

import (
	"yawebapp/controllers/http/demo"
	"yawebapp/library/inner/app"
	"yawebapp/library/inner/middlewares"
)

func HTTPRouter() {
	httpServer := app.GetHTTPServer()

	httpServer.Use(middlewares.Trace)
	httpServer.Use(middlewares.Recover())

	httpServer.NoRoute(middlewares.RouterNotFound)

	// 路由组
	DemoGr := httpServer.Group("/demo")
	{
		DemoGr.GET("/gituser", demo.GitUser)
	}
}
