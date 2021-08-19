package app

import (
	"fmt"
	"os"
	"path/filepath"
	"yawebapp/library/inner/helper"
)

var appRootPath = ""

func RootPath() string {
	if appRootPath == "" {
		appRootPath = detectRootPath()
	}
	return appRootPath
}

func BinPath() string {
	return RootPath() + "/bin"
}

func ConfPath() string {
	return RootPath() + "/conf"
}

func LogPath() string {
	logPath := RootPath() + "/log"
	if !helper.IsExist(logPath) {
		helper.MakeDirP(logPath)
	}
	return logPath
}

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
