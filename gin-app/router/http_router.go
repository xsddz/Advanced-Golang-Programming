package router

import (
	"gin-app/controller"
	"gin-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetGinRouter(app *gin.Engine) {
	app.NoRoute(middleware.RouterNotFound)

	app.GET("/", controller.Demo)
}
