package controller

import "github.com/gin-gonic/gin"

func Foobar(ctx *gin.Context) {
	ctx.JSON(200, "demo app.")
}
