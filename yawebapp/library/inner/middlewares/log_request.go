package middlewares

import (
	"fmt"
	"time"
	"yawebapp/library/inner/server"

	"github.com/gin-gonic/gin"
)

func LogRequest() gin.HandlerFunc {
	log := server.NewLogger()
	return func(ctx *gin.Context) {
		appContext := server.NewWebContext(ctx)
		// Start timer
		start := time.Now()
		path := appContext.Request.URL.Path
		raw := appContext.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		// 处理请求
		ctx.Next()

		param := gin.LogFormatterParams{
			Request:      appContext.Request,
			StatusCode:   appContext.Writer.Status(),
			ClientIP:     appContext.ClientIP(),
			Method:       appContext.Request.Method,
			Path:         path,
			ErrorMessage: appContext.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:     appContext.Writer.Size(),
			Keys:         appContext.Keys,
		}
		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		log.Info(*appContext, defaultLogFormatter(param))
	}
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("%v |%s %3d %s| %13v | %15s |%s %-7s %s %#v |%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
