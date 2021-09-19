package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func MD5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

func SHA1(data []byte) string {
	sum := sha1.Sum(data)
	return hex.EncodeToString(sum[:])
}
