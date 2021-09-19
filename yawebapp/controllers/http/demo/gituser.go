package demo

import (
	"yawebapp/entities/entitydemo"
	"yawebapp/library/infra/server"
	"yawebapp/models/service/page/demo"

	"github.com/gin-gonic/gin"
)

// GitUser -
func GitUser(ctx *gin.Context) {
	appContext := server.NewWebContextViaHTTP(ctx)
	response := server.NewHTTPResponse()

	// 业务逻辑请求参数初始化
	reqEntity := entitydemo.ReqGitUser{}
	resEntity := entitydemo.ResGitUser{}

	// 请求参数解析
	err := appContext.ShouldBind(&reqEntity)
	if err != nil {
		// 错误提前返回
		response.Error(appContext, err)
		return
	}

	// 执行业务逻辑
	err = demo.NewGitUserPage(appContext).Execute(reqEntity, &resEntity)
	if err != nil {
		// 错误提前返回
		response.Error(appContext, err)
		return
	}

	// 正常返回
	response.Success(appContext, resEntity)
}
