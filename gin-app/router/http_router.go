package router

import (
	"gin-app/controller"
	"gin-app/library/server"
	"gin-app/middleware"
)

func HTTPRouter(app *server.App) {
	app.HTTPServer.NoRoute(middleware.RouterNotFound)
	app.HTTPServer.GET("/", controller.Index)
}
