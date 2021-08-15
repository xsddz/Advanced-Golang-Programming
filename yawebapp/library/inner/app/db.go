package app

import (
	"gorm.io/gorm"
)

var (
	dbDriver = "sqlite"
	dbTable  = make(map[string]*db)
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

// DB 提供延迟初始化对应资源连接的能力，及连接复用能力
func DB(database string) *db {
	if d, ok := dbTable[database]; ok {
		return d
	}

	var d *db
	switch dbDriver {
	case "sqlite":
		d = &db{initSQLite(database)}
	case "mysql":
		d = &db{initMySQL(database)}
	}
	dbTable[database] = d

	return dbTable[database]
}
