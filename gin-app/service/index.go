package service

import (
	"fmt"
	"gin-app/entity"
	"gin-app/library/server"
	"gin-app/model"

	"gorm.io/gorm"
)

type IndexService struct {
	ctx *server.WebContext
}

func NewIndexService(ctx *server.WebContext) *IndexService {
	return &IndexService{ctx: ctx}
}

func (s *IndexService) Execute(req entity.ReqIndex, res *entity.ResIndex) server.Error {
	err := server.SQLite.AutoMigrate(&model.Index{})
	if err != nil {
		return server.NewError(-1, "AutoMigrate"+err.Error())
	}

	var data model.Index
	result := server.SQLite.First(&data, "name = ?", req.Name)
	if result.Error == nil {
		res.Words = fmt.Sprintf("你好：%v, %v", data.Name, data.CreateTs)
		return server.ErrorNone
	}

	if result.Error != gorm.ErrRecordNotFound {
		return server.NewError(-1, "查询数据出错："+result.Error.Error())
	}

	data = model.Index{Name: req.Name}
	result = server.SQLite.Create(&data)
	if result.Error != nil {
		return server.NewError(-1, "创建数据出错："+result.Error.Error())
	}
	res.Words = fmt.Sprintf("你好：%v, %v", data.Name, data.CreateTs)
	return server.ErrorNone
}
