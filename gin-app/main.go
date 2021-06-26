package main

import (
	"gin-app/router"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	router.SetGinRouter(app)
	app.Run()
}
