package server

import (
	"regexp"
	"strconv"
	"strings"
)

// 错误码格式定义: [code] message
const errRegexString = `^\[[-]?\d+\]`

// 错误码格式解析正则
var errRegex = regexp.MustCompile(errRegexString)

// ParseError 从约定的错误格式中解析错误码和错误信息
func ParseError(e error) (code int, message string) {
	codeArr := errRegex.FindStringSubmatch(e.Error())
	if len(codeArr) == 0 {
		return -1, e.Error()
	}

	code, _ = strconv.Atoi(strings.Trim(codeArr[0], "]["))
	message = strings.Replace(e.Error(), codeArr[0], "", 1)
	return code, strings.TrimLeft(message, " ")
}
