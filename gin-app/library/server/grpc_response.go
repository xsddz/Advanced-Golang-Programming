package server

import (
	"encoding/json"
)

type GRPCResponse struct {
	pb interface{}
}

func NewGRPCResponse(pb interface{}) *GRPCResponse {
	return &GRPCResponse{pb: pb}
}

func (res *GRPCResponse) PB() interface{} {
	return res.pb
}

func (res *GRPCResponse) Success(ctx *WebContext, data interface{}) {
	j, _ := json.Marshal(map[string]interface{}{
		"code":    0,
		"message": "ok",
		"data":    data,
	})

	json.Unmarshal(j, res.pb)
}

func (res *GRPCResponse) Error(ctx *WebContext, e error) {
	code, message := ParseError(e)
	j, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
	})

	json.Unmarshal(j, res.pb)
}
