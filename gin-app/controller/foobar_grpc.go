package controller

import (
	"context"
	"encoding/json"
	"gin-app/entity"
	"gin-app/entity/pb"
	"gin-app/library/server"
	"gin-app/service"
)

// FoobarController 定义服务
type FoobarController struct{}

// Index 实现 Index 方法
func (s *FoobarController) Index(ctx context.Context, req *pb.RequestIndex) (*pb.ResponseIndex, error) {
	appContext := &server.WebContext{}
	response := pb.ResponseIndex{}

	// 业务逻辑请求参数初始化
	reqEntity := entity.ReqIndex{}
	resEntity := entity.ResIndex{}

	// 请求参数转换
	reqJson, _ := json.Marshal(req)
	err := json.Unmarshal(reqJson, &reqEntity)
	if err != nil {
		response.Code = -1
		response.Message = "参数错误：" + err.Error()
		return &response, nil
	}
	// 必传参数
	if reqEntity.Name == "" {
		response.Code = -1
		response.Message = "参数错误：name"
		return &response, nil
	}

	// 执行业务逻辑
	appErr := service.NewIndexService(appContext).Execute(reqEntity, &resEntity)

	// 返回
	response.Code = int32(appErr.Code())
	response.Message = appErr.Message()
	// 返回数据转换
	resJson, _ := json.Marshal(resEntity)
	json.Unmarshal(resJson, &response.Data)
	return &response, nil
}
