package router

import (
	"gin-app/controller"
	"gin-app/library/server"
	"gin-app/middleware"
)

func HTTPRouter(app *server.App) {
	app.ServerHTTP.NoRoute(middleware.RouterNotFound)
	app.ServerHTTP.GET("/", controller.Foobar)
}
