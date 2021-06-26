package middleware

import "github.com/gin-gonic/gin"

func RouterNotFound(ctx *gin.Context) {
	ctx.JSON(400, "404 not found.")
}
