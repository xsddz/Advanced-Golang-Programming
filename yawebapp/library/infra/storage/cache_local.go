package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"yawebapp/library/infra/helper"
)

// LocalCache -
type LocalCache struct {
	file string
	data map[string]string
}

func NewLocalCache(conf RedisConf, dataPath string) (*LocalCache, error) {
	// ruler:: data/cache/{{mastername}}_{{db}}.cache
	cacheFile := fmt.Sprintf("%v/cache/%v_%v.cache", dataPath, conf.MasterName, conf.DefaultDB)
	if err := helper.TouchFile(cacheFile); err != nil {
		return nil, err
	}

	dataRaw, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	if len(dataRaw) > 0 {
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
			return nil, err
		}
	}

	return &LocalCache{cacheFile, data}, nil
}

func (c *LocalCache) Dump() error {
	data, _ := json.Marshal(c.data)
	return ioutil.WriteFile(c.file, data, os.ModePerm)
}

func (c *LocalCache) Get(key string) string {
	if val, ok := c.data[key]; ok {
		return val
	}
	return ""
}

func (c *LocalCache) Set(key string, val string) {
	c.data[key] = val
}
