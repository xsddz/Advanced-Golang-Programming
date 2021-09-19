package middlewares

import (
	"yawebapp/library/infra/app"
	"yawebapp/library/infra/server"

	"github.com/gin-gonic/gin"
)

func RouterNotFound(ctx *gin.Context) {
	appContext := server.NewWebContextViaHTTP(ctx)
	response := server.HttpResponse{}
	response.Error(appContext, app.Error404)
	ctx.Abort()
}
