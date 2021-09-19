package helper_test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
)

// go test -bench="BenchmarkMD5V3"
func BenchmarkMD5V3(b *testing.B) {
	md5v3 := func(data []byte) string {
		sum := md5.Sum(data)
		return fmt.Sprintf("%x\n", sum)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		md5v3([]byte("12345678901234567890"))
	}
}

// go test -bench="BenchmarkMD5V2"
func BenchmarkMD5V2(b *testing.B) {
	md5v2 := func(data []byte) string {
		hash := md5.New()
		io.WriteString(hash, string(data))
		return hex.EncodeToString(hash.Sum(nil))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		md5v2([]byte("12345678901234567890"))
	}
}

// go test -bench="BenchmarkMD5"
func BenchmarkMD5(b *testing.B) {
	md5v1 := func(data []byte) string {
		sum := md5.Sum(data)
		return hex.EncodeToString(sum[:])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		md5v1([]byte("12345678901234567890"))
	}
}
