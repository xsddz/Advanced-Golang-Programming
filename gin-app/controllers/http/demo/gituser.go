package demo

import (
	"gin-app/entity/demoentity"
	"gin-app/library/server"
	"gin-app/models/page/demopage"

	"github.com/gin-gonic/gin"
)

// GitUser -
func GitUser(ctx *gin.Context) {
	appContext := server.NewWebContext(ctx)
	response := server.NewHTTPResponse()

	// 业务逻辑请求参数初始化
	reqEntity := demoentity.ReqGitUser{}
	resEntity := demoentity.ResGitUser{}

	// 请求参数解析
	err := appContext.ShouldBind(&reqEntity)
	if err != nil {
		// 错误提前返回
		response.Error(appContext, err)
		return
	}

	// 执行业务逻辑
	err = demopage.NewGitUserPage(appContext).Execute(reqEntity, &resEntity)
	if err != nil {
		// 错误提前返回
		response.Error(appContext, err)
		return
	}

	// 正常返回
	response.Success(appContext, resEntity)
}
