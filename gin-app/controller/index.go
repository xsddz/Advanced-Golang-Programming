package controller

import (
	"gin-app/entity"
	"gin-app/library/server"
	"gin-app/service"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	ctx.Set("trace_id", "111111")
	appContext := &server.WebContext{Context: ctx}
	response := server.HttpResponse{}

	// 业务逻辑请求参数初始化
	reqEntity := entity.ReqIndex{}
	resEntity := entity.ResIndex{}

	// 请求参数解析
	err := ctx.ShouldBind(&reqEntity)
	if err != nil {
		response.Error(appContext, server.NewError(-1, "参数错误："+err.Error()))
		return
	}

	// 执行业务逻辑
	appErr := service.NewIndexService(appContext).Execute(reqEntity, &resEntity)
	if appErr != server.ErrorNone {
		// 错误返回
		response.Error(appContext, appErr)
		return
	}

	// 正常返回
	response.Success(appContext, resEntity)
}
