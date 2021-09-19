package helper

import (
	"os"
)

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

func MakeDirP(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
