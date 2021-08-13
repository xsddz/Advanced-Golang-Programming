package config

import "github.com/BurntSushi/toml"

func LoadConf(filename string, obj interface{}) error {
	if _, err := toml.DecodeFile(filename, obj); err != nil {
		return err
	}
	return nil
}
