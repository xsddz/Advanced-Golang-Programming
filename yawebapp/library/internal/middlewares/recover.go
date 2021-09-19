package middlewares

import (
	"errors"
	"fmt"
	"yawebapp/library/inner/app"
	"yawebapp/library/inner/server"

	"github.com/gin-gonic/gin"
)

func Recover(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			app.Logger.Critical(ctx, fmt.Sprint(err))

			appContext := server.NewWebContextViaHTTP(ctx)
			response := server.HttpResponse{}
			response.Error(appContext, errors.New("[500] internal server error"))
			ctx.Abort()
		}
	}()

	ctx.Next()
}
