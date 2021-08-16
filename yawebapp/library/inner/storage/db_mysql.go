package storage

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConf struct {
	ClusterName       string       `toml:"cluster_name"`
	MaxOpenConns      int          `toml:"max_open_conns"`
	MaxIdleConns      int          `toml:"max_idle_conns"`
	ConnMaxLifetimeMs int          `toml:"conn_max_lifetime_ms"`
	BalanceStrategy   string       `toml:"balance_strategy"`
	Charset           string       `toml:"charset"`
	Username          string       `toml:"username"`
	Password          string       `toml:"password"`
	DefaultDB         string       `toml:"default_db"`
	Hosts             []DBConfHost `toml:"host"`
}

type DBConfHost struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}

func (c *DBConf) DSN(index int) string {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		c.Username, c.Password, c.Hosts[index].IP, c.Hosts[index].Port, c.DefaultDB, c.Charset)
	return dsn
}

func NewMySQL(conf DBConf) ([]*gorm.DB, error) {
	var dbs []*gorm.DB
	for index := range conf.Hosts {
		db, err := gorm.Open(mysql.Open(conf.DSN(index)), &gorm.Config{})
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
