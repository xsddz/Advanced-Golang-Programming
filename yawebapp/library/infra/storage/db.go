package storage

import "fmt"

// DBConf -
type DBConf struct {
	Driver            string `toml:"driver"`
	ClusterName       string `toml:"cluster_name"`
	MaxOpenConns      int    `toml:"max_open_conns"`
	MaxIdleConns      int    `toml:"max_idle_conns"`
	ConnMaxLifetimeMs int    `toml:"conn_max_lifetime_ms"`
	BalanceStrategy   string `toml:"balance_strategy"`
	Charset           string `toml:"charset"`
	Username          string `toml:"username"`
	Password          string `toml:"password"`
	DefaultDB         string `toml:"default_db"`
	HostPorts         string `toml:"hosts"`
	Hosts             []DBConfHost
}

// DBConfHost -
type DBConfHost struct {
	IP   string
	Port int
}

// DSN -
func (c *DBConf) DSN(index int) string {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		c.Username, c.Password, c.Hosts[index].IP, c.Hosts[index].Port, c.DefaultDB, c.Charset)
	return dsn
}
