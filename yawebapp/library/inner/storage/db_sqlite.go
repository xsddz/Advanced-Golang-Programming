package storage

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"yawebapp/library/inner/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLite(dataPath string, conf DBConf, l logger.Interface) ([]*gorm.DB, error) {
	// ruler:: data/db_{{clustername}}/{{dbname}}.db
	dbFile := fmt.Sprintf("%v/db_%v/%v.db", dataPath, strings.ToLower(conf.ClusterName), conf.DefaultDB)
	dbFileDir := filepath.Dir(dbFile)
	if !utils.IsDir(dbFileDir) {
		utils.MakeDirP(dbFileDir)
	}

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: l,
	})
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

	return []*gorm.DB{db}, nil
}
