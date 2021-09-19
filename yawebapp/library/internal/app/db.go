package app

import (
	"context"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

var (
	dbDriver = "sqlite"
	dbTable  = make(map[string][]*gorm.DB)
)

func initDBDriver(driver string) {
	if driver != "mysql" && driver != "sqlite" {
		driver = "sqlite"
	}
	dbDriver = driver
}

// db 通过嵌套gorm.DB，利用gorm能力提供统一的数据库操作方法。对于新的数据库，可以实现对应的gorm驱动
type db struct {
	*gorm.DB
}

func (d *db) Ping() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *db) WithContext(ctx context.Context) *db {
	d.DB = d.DB.WithContext(ctx)
	return d
}

// DB 提供延迟初始化对应资源连接的能力，及连接复用能力
// 利用 ... 提供可选参数能力，只取第一个或走默认的
func DB(ctx context.Context, clusterNames ...string) *db {
	var clusterName string
	if len(clusterNames) == 0 {
		clusterName = "Default"
	} else {
		clusterName = clusterNames[0]
	}

	if _, ok := dbTable[clusterName]; !ok {
		var dbs []*gorm.DB
		switch dbDriver {
		case "sqlite":
			dbs = initSQLite(clusterName)
		case "mysql":
			dbs = initMySQL(clusterName)
		}
		dbTable[clusterName] = dbs
	}

	return randDB(dbTable[clusterName]).WithContext(ctx)
}

func randDB(dbs []*gorm.DB) *db {
	rand.Seed(time.Now().UnixNano())
	return &db{dbs[rand.Intn(len(dbs))]}
}
