package helper

import (
	"os"
	"path/filepath"
)

// IsDir 是否为目录
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if !os.IsNotExist(err) && fileInfo.IsDir() {
		return true
	}
	return false
}

// IsExist 检测目录/文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// MakeDirP 同makedir -p
func MakeDirP(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// TouchFile 同touch file
func TouchFile(file string) error {
	file, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	if IsExist(file) {
		return nil
	}

	fileDir := filepath.Dir(file)
	if !IsDir(fileDir) {
		if err = MakeDirP(fileDir); err != nil {
			return err
		}
	}

	_, err = os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666)
	return err
}
