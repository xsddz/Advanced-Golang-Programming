package routers

import (
	"yawebapp/controllers/http/demo"
	"yawebapp/library/common/app"
	"yawebapp/library/common/middlewares"
)

func HTTPRouter() {
	httpServer := app.GetHTTPServer()
	httpServer.NoRoute(middlewares.RouterNotFound)

	// 路由组
	DemoGr := httpServer.Group("/demo")
	{
		DemoGr.GET("/gituser", demo.GitUser)
	}
}
