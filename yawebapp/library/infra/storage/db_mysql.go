package storage

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMySQL -
func NewMySQL(conf DBConf, l logger.Interface) ([]*gorm.DB, error) {
	var dbs []*gorm.DB
	for index := range conf.Hosts {
		db, err := gorm.Open(mysql.Open(conf.DSN(index)), &gorm.Config{Logger: l})
		if err != nil {
			return nil, err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		// SetMaxOpenConns 设置最大的并发打开连接数，设置这个数小于等于0则表示没有限制，也就是默认设置。
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
		// SetMaxIdleConns 设置最大的空闲连接数，设置小于等于0的数意味着不保留空闲连接。
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		// SetConnMaxLifetime 设置连接的最大生命周期，设置为0的话意味着没有最大生命周期，连接总是可重用(默认行为)。
		sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeMs) * time.Millisecond)

		dbs = append(dbs, db)
	}

	return dbs, nil
}
