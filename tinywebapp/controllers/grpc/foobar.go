package grpc

import (
	"context"
	"gin-app/entity/demoentity"
	"gin-app/entity/pb"
	"gin-app/library/server"
	"gin-app/models/page/demopage"
)

// Foobar 定义服务
type Foobar struct{}

// Index 实现 Index 方法
func (s *Foobar) Index(ctx context.Context, req *pb.RequestGitUser) (*pb.ResponseGitUser, error) {
	appContext := server.NewWebContext(nil)
	response := server.NewGRPCResponse(&pb.ResponseGitUser{})

	// 业务逻辑请求参数初始化
	reqEntity := demoentity.ReqGitUser{}
	resEntity := demoentity.ResGitUser{}

	// 请求参数解析
	err := appContext.ShouldBindGRPC(req, &reqEntity)
	if err != nil {
		// 错误提前返回
		response.Error(appContext, err)
		return response.PB().(*pb.ResponseGitUser), nil
	}

	// 执行业务逻辑
	err = demopage.NewGitUserPage(appContext).Execute(reqEntity, &resEntity)
	if err != nil {
		// 错误提前返回
		response.Error(appContext, err)
		return response.PB().(*pb.ResponseGitUser), nil
	}

	// 正常返回
	response.Success(appContext, resEntity)
	return response.PB().(*pb.ResponseGitUser), nil
}
