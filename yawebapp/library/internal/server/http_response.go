package server

type HttpResponse struct {
	TraceID string      `json:"trace_id"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewHTTPResponse() *HttpResponse {
	return &HttpResponse{}
}

func (res *HttpResponse) Success(ctx *WebContext, data interface{}) {
	res.TraceID = ctx.GetString("trace_id")
	res.Code = 0
	res.Message = "ok"
	res.Data = data

	ctx.JSON(200, res)
}

func (res *HttpResponse) Error(ctx *WebContext, err error) {
	code, message := ParseError(err)

	res.TraceID = ctx.GetString("trace_id")
	res.Code = code
	res.Message = message

	ctx.JSON(200, res)
}
