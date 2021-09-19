package app_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"yawebapp/library/infra/app"
	"yawebapp/library/infra/storage"
)

// go test -run="TestLoadConf"
func TestLoadConf(t *testing.T) {
	// mock app
	rootPath, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("app root path error => %v", err)
	}
	app.SetRootPath(rootPath)
	app.Init()

	for _, cl := range caseTable {
		if err := cl(); err != nil {
			t.Fatalf("load conf error => %v", err)
		}
	}
}

var caseTable = []func() error{
	loadAppConf,
	loadDBConf,
	loadCacheConf,
}

func loadAppConf() error {
	var c app.AppConf
	err := app.LoadConf("app", &c)
	if err != nil {
		return err
	}
	fmt.Println("---->app:", c)
	return nil
}

func loadDBConf() error {
	var c map[string]storage.DBConf
	err := app.LoadConf("db", &c)
	if err != nil {
		return err
	}
	fmt.Println("---->db:", c)
	return nil
}

func loadCacheConf() error {
	var c map[string]*storage.RedisConf
	err := app.LoadConf("cache", &c)
	if err != nil {
		return err
	}
	fmt.Println("---->cache:", c)
	fmt.Println("---->cache:", c["Default"])
	return nil
}
