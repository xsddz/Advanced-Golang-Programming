package middleware

import (
	"errors"
	"gin-app/library/server"

	"github.com/gin-gonic/gin"
)

func RouterNotFound(ctx *gin.Context) {
	response := server.HttpResponse{}
	response.Error(&server.WebContext{Context: ctx}, errors.New("[404] not found"))
	ctx.Abort()
}
