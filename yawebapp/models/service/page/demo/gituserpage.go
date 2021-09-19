package demo

import (
	"encoding/json"
	"fmt"
	"yawebapp/entities/entitydemo"
	"yawebapp/library/infra/app"
	"yawebapp/library/infra/server"
	"yawebapp/models/dao"

	"gorm.io/gorm"
)

type GitUserPage struct {
	ctx *server.WebContext
}

func NewGitUserPage(ctx *server.WebContext) *GitUserPage {
	return &GitUserPage{ctx: ctx}
}

func (p *GitUserPage) Execute(req entitydemo.ReqGitUser, res *entitydemo.ResGitUser) error {
	key := "dzh:test"
	val, err := app.Cache().Get(key)
	app.Logger.Debug(*p.ctx, "get", key, val, err)
	if err == app.ErrorNotExist {
		app.Logger.Info(*p.ctx, "key does not exist")
	} else if err != nil {
		app.Logger.Error(*p.ctx, err.Error())
	} else {
		app.Logger.Debug(*p.ctx, "get", key, val)
	}

	var data dao.CommonConfig
	result := app.DB(p.ctx).First(&data, "name = ?", req.Name)
	if result.Error == nil {
		j, _ := json.Marshal(data)
		res.Info = fmt.Sprintf("你好：%v", string(j))
		return nil
	}
	if result.Error != gorm.ErrRecordNotFound {
		err := fmt.Errorf("查询数据出错：%v", result.Error.Error())
		app.Logger.Error(p.ctx, err.Error())
		return err
	}
	return result.Error
}
