package data

import (
	"gin-app/library/server"
	"gin-app/model/dao"
	"time"
)

type IndexData struct {
	ctx *server.WebContext
}

func NewIndexData(ctx *server.WebContext) *IndexData {
	return &IndexData{ctx: ctx}
}

func (d *IndexData) FindByName(name string) (dao.IndexDao, server.AppError) {
	res := dao.IndexDao{Name: name, CreateTs: time.Now().Unix()}
	return res, server.AppErrorNone
}
