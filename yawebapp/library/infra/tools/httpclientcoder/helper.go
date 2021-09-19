package httpclientcoder

import "strings"

func TrimSpace(s string) string {
	// 忽略空白、tab、换行字符
	s = strings.Trim(s, " \t\n\r")

	// 压缩中间的tab、空格字符
	s = strings.ReplaceAll(s, "\t", " ")
	s = regTrSpace.ReplaceAllString(s, "$1")

	return s
}
