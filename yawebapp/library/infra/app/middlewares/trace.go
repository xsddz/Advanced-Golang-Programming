package middlewares

import (
	"fmt"
	"time"
	"yawebapp/library/infra/app"
	"yawebapp/library/infra/helper"

	"github.com/gin-gonic/gin"
)

// Trace 追踪http请求
func Trace(ctx *gin.Context) {
	setTraceID(ctx)
	param := startRequest(ctx)

	// 处理请求
	ctx.Next()

	endRequest(ctx, param)
}

// setTraceID set trace_id from request header or generate new
var setTraceID = func(ctx *gin.Context) {
	traceID := ctx.GetHeader("TRACE_ID")
	if traceID == "" {
		traceID = helper.NextRequestID()
	}
	ctx.Set("trace_id", traceID)
}

var startRequest = func(ctx *gin.Context) *gin.LogFormatterParams {
	// Start timer
	return &gin.LogFormatterParams{
		TimeStamp: time.Now(),
	}
}

var endRequest = func(ctx *gin.Context, param *gin.LogFormatterParams) {
	start := param.TimeStamp
	// Stop timer
	param.TimeStamp = time.Now()
	param.Latency = param.TimeStamp.Sub(start)

	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}

	param.Request = ctx.Request
	param.ClientIP = ctx.ClientIP()
	param.Method = ctx.Request.Method
	param.Path = path
	param.StatusCode = ctx.Writer.Status()
	param.BodySize = ctx.Writer.Size()
	param.Keys = ctx.Keys
	param.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()

	app.Logger.Info(ctx, formatRequestParam(param))
}

// formatRequestParam is the default log format function Logger middleware uses.
var formatRequestParam = func(param *gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("[%v] %s[%d]%s %s[%v]%s [%s] %s%s%s %#v %s",
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		helper.ColorBlueBold, param.StatusCode, helper.ColorReset,
		helper.ColorYellowBold, param.Latency, helper.ColorReset,
		param.ClientIP,
		helper.ColorGreenBold, param.Method, helper.ColorReset,
		param.Path,
		param.ErrorMessage,
	)
}
