package common

import (
	"crypto/md5"
	"encoding/base64"
	"strings"
)

// GenerateUsername -
func GenerateUsername(account string) (name string) {
	// TODO: generate uniq random meaningful name
	hash := md5.Sum([]byte(account))
	name = base64.StdEncoding.EncodeToString(hash[:])
	name = strings.TrimRight(name, "=")
	return
}
