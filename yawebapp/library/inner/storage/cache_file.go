package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"yawebapp/library/inner/utils"

	"github.com/go-redis/redis"
)

func NewFileCache(dataPath string, conf RedisConf) ([]*redis.Client, error) {
	// ruler:: data/cache_{{mastername}}/{{db}}.cache
	cacheFile := fmt.Sprintf("%v/cache_%v/%v.cache", dataPath, conf.MasterName, conf.DefaultDB)
	cacheFileDir := filepath.Dir(cacheFile)
	if !utils.IsDir(cacheFileDir) {
		utils.MakeDirP(cacheFileDir)
	}

	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var dataM map[string]string
	err = json.Unmarshal(data, &dataM)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// FileCacheClient 实现cacher接口
type FileCacheClient struct {
	file string
	data map[string]string
}

func (f *FileCacheClient) Dump() {
	data, _ := json.Marshal(f.data)
	ioutil.WriteFile(f.file, data, os.ModePerm)
}
