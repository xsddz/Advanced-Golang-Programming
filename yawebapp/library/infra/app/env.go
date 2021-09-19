package app

import (
	"fmt"
	"os"
	"path/filepath"
	"yawebapp/library/infra/helper"
)

var appRootPath = ""

// SetRootPath 设置app的根目录
func SetRootPath(path string) {
	appRootPath = path
}

// RootPath 获取app的根目录
func RootPath() string {
	if appRootPath == "" {
		appRootPath = detectRootPath()
	}
	return appRootPath
}

// BinPath 获取app的bin目录
func BinPath() string {
	return RootPath() + "/bin"
}

// ConfPath 获取app的conf目录
func ConfPath() string {
	return RootPath() + "/conf"
}

// LogPath 获取app的log目录
func LogPath() string {
	logPath := RootPath() + "/log"
	if !helper.IsExist(logPath) {
		helper.MakeDirP(logPath)
	}
	return logPath
}

// DataPath 获取app的data目录
func DataPath() string {
	dataPath := RootPath() + "/data"
	if !helper.IsExist(dataPath) {
		helper.MakeDirP(dataPath)
	}
	return dataPath
}

func detectRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("detect root path error: %v", err))
	}

	// check bin
	if filepath.Base(dir) == "bin" {
		return filepath.Dir(dir)
	}

	// check conf
	if helper.IsDir(dir + "/conf") {
		return dir
	}

	// use current dir
	currDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("detect root path for current error: %v", err))
	}
	return currDir
}
