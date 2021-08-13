package app

import "errors"

var (
	ErrorNone  = errors.New("[0] ok")
	ErrorParam = errors.New("[1] 参数错误")
)
