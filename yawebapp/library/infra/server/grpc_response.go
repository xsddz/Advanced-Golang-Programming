package server

import (
	"encoding/json"
)

// GRPCResponse -
type GRPCResponse struct {
	pb interface{}
}

// NewGRPCResponse -
func NewGRPCResponse(pb interface{}) *GRPCResponse {
	return &GRPCResponse{pb: pb}
}

// PB -
func (res *GRPCResponse) PB() interface{} {
	return res.pb
}

// Success -
func (res *GRPCResponse) Success(ctx *WebContext, data interface{}) {
	j, _ := json.Marshal(map[string]interface{}{
		"trace_id": ctx.Value("trace_id"),
		"code":     0,
		"message":  "ok",
		"data":     data,
	})

	json.Unmarshal(j, res.pb)
}

// Error -
func (res *GRPCResponse) Error(ctx *WebContext, e error) {
	code, message := ParseError(e)
	j, _ := json.Marshal(map[string]interface{}{
		"trace_id": ctx.Value("trace_id"),
		"code":     code,
		"message":  message,
	})

	json.Unmarshal(j, res.pb)
}
