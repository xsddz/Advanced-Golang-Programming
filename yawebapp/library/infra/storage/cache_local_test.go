package storage_test

import (
	"encoding/json"
	"testing"
	"yawebapp/library/infra/storage"
)

// go test -run "TestLocalCache"
func TestLocalCache(t *testing.T) {
	conf := storage.RedisConf{
		MasterName: "unittest",
		DefaultDB:  "1",
	}
	dataPath := "."
	cache, err := storage.NewLocalCache(conf, dataPath)
	if err != nil {
		j, _ := json.Marshal(conf)
		t.Fatalf("new local cache failed: %v, %s, %v", err, j, dataPath)
	}

	key, val := "foo", "bar"
	if got := cache.Get(key); got == "" {
		cache.Set(key, val)
	}
	if got := cache.Get(key); got != val {
		t.Fatalf("local cache get/set not match")
	}

	if err := cache.Dump(); err != nil {
		t.Fatalf("dump cache error: %v", err)
	}
}
