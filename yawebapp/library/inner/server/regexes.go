package server

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	// error format: [code] message
	errorRegexString = `^\[[-]?\d+\]`
)

var (
	errorRegex = regexp.MustCompile(errorRegexString)
)

func ParseError(e error) (code int, message string) {
	codeArr := errorRegex.FindStringSubmatch(e.Error())
	if len(codeArr) == 0 {
		return -1, e.Error()
	}

	code, _ = strconv.Atoi(strings.Trim(codeArr[0], "]["))
	message = strings.Replace(e.Error(), codeArr[0], "", 1)
	return code, strings.TrimLeft(message, " ")
}
