package app

import (
	"fmt"
	"os"
	"path/filepath"
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
	if !IsExist(logPath) {
		os.MkdirAll(logPath, os.ModePerm)
	}
	return logPath
}

func DataPath() string {
	dataPath := RootPath() + "/data"
	if !IsExist(dataPath) {
		os.MkdirAll(dataPath, os.ModePerm)
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
	if IsDir(dir + "/conf") {
		return dir
	}

	// use current dir
	currDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("detect root path for current error: %v", err))
	}
	return currDir
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if !os.IsNotExist(err) && fileInfo.IsDir() {
		return true
	}
	return false
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
