package page

import (
	"gin-app/entity"
	"gin-app/library/server"
)

type IndexPage struct {
	ctx *server.WebContext
}

func NewIndexPage(ctx *server.WebContext) *IndexPage {
	return &IndexPage{
		ctx: ctx,
	}
}

func (p *IndexPage) Execute(req entity.ReqIndex, res *entity.ResIndex) server.AppError {

	res.Words = "你好，" + req.Name

	return server.AppErrorNone
}
