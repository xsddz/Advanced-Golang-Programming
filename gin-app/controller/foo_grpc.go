package controller

import (
	"context"
	"gin-app/entity/pb"
)

// FooController 定义我们的服务
type FooController struct{}

// Bar 实现Bar方法
func (s *FooController) Bar(ctx context.Context, req *pb.BarRequest) (*pb.BarResponse, error) {
	res := pb.BarResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return &res, nil
}
