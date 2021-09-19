package middlewares

import (
	"fmt"
	"yawebapp/library/infra/app"
	"yawebapp/library/infra/server"

	"github.com/gin-gonic/gin"
)

func Recover(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			appContext := server.NewWebContextViaHTTP(ctx)

			app.Logger.Critical(appContext, fmt.Sprint(err))

			response := server.HttpResponse{}
			response.Error(appContext, app.Error500)
			ctx.Abort()
		}
	}()

	ctx.Next()
}
