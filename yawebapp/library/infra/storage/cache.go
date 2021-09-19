package storage

// RedisConf -
type RedisConf struct {
	Driver     string `toml:"driver"`
	Type       string `toml:"type"`
	MasterName string `toml:"master_name"`
	Passowrd   string `toml:"password"`
	DefaultDB  string `toml:"default_db"`
	HostPorts  string `toml:"hosts"`
	Hosts      []RedisConfHost
}

// RedisConfHost -
type RedisConfHost struct {
	IP   string
	Port int
}
