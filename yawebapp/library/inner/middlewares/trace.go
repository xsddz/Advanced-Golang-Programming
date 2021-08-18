package middlewares

import (
	"yawebapp/library/inner/utils"

	"github.com/gin-gonic/gin"
)

func Trace(ctx *gin.Context) {
	traceID := ctx.Request.Header.Get("TRACE_ID")
	if traceID == "" {
		traceID = utils.GenrateRequestID()
	}
	ctx.Set("trace_id", traceID)

	// 处理请求
	ctx.Next()
}
