package page

import (
	"fmt"
	"gin-app/entity"
	"gin-app/library/server"
	"gin-app/model/service/data"
)

type IndexPage struct {
	ctx *server.WebContext

	dataIndex *data.IndexData
}

func NewIndexPage(ctx *server.WebContext) *IndexPage {
	return &IndexPage{
		ctx:       ctx,
		dataIndex: data.NewIndexData(ctx),
	}
}

func (p *IndexPage) Execute(req entity.ReqIndex, res *entity.ResIndex) server.AppError {

	data, appErr := p.dataIndex.FindByName(req.Name)
	if appErr != server.AppErrorNone {
		return appErr
	}
	res.Words = fmt.Sprintf("你好：%v, %v", data.Name, data.CreateTs)

	return server.AppErrorNone
}
