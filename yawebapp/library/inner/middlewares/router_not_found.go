package middlewares

import (
	"errors"
	"yawebapp/library/inner/server"

	"github.com/gin-gonic/gin"
)

func RouterNotFound(ctx *gin.Context) {
	appContext := server.NewWebContextViaHTTP(ctx)
	response := server.HttpResponse{}
	response.Error(appContext, errors.New("[404] not found"))
	ctx.Abort()
}
