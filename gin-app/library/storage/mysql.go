package storage

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlConf struct {
	Host                string
	Port                string
	Username            string
	Password            string
	Database            string
	Charset             string
	MaxIdleConns        int
	MaxOpenConns        int
	ConnMaxLifetimeHour int
}

func (c *MysqlConf) DSN() string {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
	return dsn
}

func NewMySQL(conf MysqlConf) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeHour) * time.Hour)

	return db, err
}
