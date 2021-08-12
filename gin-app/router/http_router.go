package router

import (
	"gin-app/controller"
	"gin-app/library/middleware"
	"gin-app/library/server"
)

func HTTPRouter(app *server.Engine) {
	app.HTTPServer.NoRoute(middleware.RouterNotFound)

	// 路由组
	foobarGr := app.HTTPServer.Group("/foobar")
	{
		foobarGr.GET("/index", controller.Index)
	}
}
