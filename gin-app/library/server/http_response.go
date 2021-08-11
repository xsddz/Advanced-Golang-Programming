package server

import "github.com/gin-gonic/gin"

type HttpResponse struct {
	TraceID string      `json:"trace_id"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (res *HttpResponse) Success(ctx *gin.Context, data interface{}) {
	res.TraceID = ctx.GetString("trace_id")
	res.Code = 0
	res.Message = "ok"
	res.Data = data

	ctx.JSON(200, res)
}

func (res *HttpResponse) Error(ctx *gin.Context, ae AppErrorI) {
	res.TraceID = ctx.GetString("trace_id")
	res.Code = ae.Code()
	res.Message = ae.Message()

	ctx.JSON(200, res)
}
