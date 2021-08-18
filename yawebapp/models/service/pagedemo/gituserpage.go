package pagedemo

import (
	"encoding/json"
	"fmt"
	"yawebapp/entities/entitydemo"
	"yawebapp/library/inner/app"
	"yawebapp/library/inner/server"
	"yawebapp/models/dao"

	"github.com/go-redis/redis"
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
	val, err := app.Cache().Get(key).Result()
	app.Logger.Debug(*p.ctx, fmt.Sprint(val, err))
	if err == redis.Nil {
		app.Logger.Info(*p.ctx, "key does not exist")
	} else if err != nil {
		app.Logger.Error(*p.ctx, err)
	} else {
		app.Logger.Debug(*p.ctx, "key:"+val)
	}

	var data dao.CommonConfig
	result := app.DB().First(&data, "name = ?", req.Name)
	if result.Error == nil {
		j, _ := json.Marshal(data)
		res.Info = fmt.Sprintf("你好：%v", string(j))
		return nil
	}
	if result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询数据出错：%v", result.Error.Error())
	}
	return result.Error
}
