package middleware

import (
	"gin-app/library/server"

	"github.com/gin-gonic/gin"
)

func RouterNotFound(ctx *gin.Context) {
	response := server.HttpResponse{}
	response.Error(ctx, server.NewAppError(400, "404 not found"))
	ctx.Abort()
}
