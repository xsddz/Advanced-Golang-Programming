package idutils

import (
	"math"
	"strings"
)

// A34Encode - 将10进制id编码为34进制字符串
func A34Encode(num int64, salt string) string {
	return EncodeGenerator(AlphabetShuffle(Alpha34, salt))(num)
}

// A34Decode - 将34进制字符串解码为10进制id
func A34Decode(input string, salt string) int64 {
	return DecoderGenerator(AlphabetShuffle(Alpha34, salt))(input)
}

// A58Encode - 将10进制id编码为58进制字符串
func A58Encode(num int64, salt string) string {
	return EncodeGenerator(AlphabetShuffle(Alpha58, salt))(num)
}

// A58Decode - 将58进制字符串解码为10进制id
func A58Decode(input string, salt string) int64 {
	return DecoderGenerator(AlphabetShuffle(Alpha58, salt))(input)
}

// EncodeGenerator - 进制编码生成器
func EncodeGenerator(alpha string) func(int64) string {
	alphaLen := int64(len(alpha))
	return func(num int64) string {
		hash := ""
		for num != 0 {
			hash = string(alpha[num%alphaLen]) + hash
			num = num / alphaLen
		}
		return hash
	}
}

// DecoderGenerator - 进制解码生成器
func DecoderGenerator(alpha string) func(input string) int64 {
	alphaLen := int64(len(alpha))
	return func(input string) int64 {
		num := 0.0
		length := len(input)
		for i := 0; i < length; i++ {
			index := strings.Index(alpha, string(input[i]))
			num += float64(index) * math.Pow(float64(alphaLen), float64(length-1-i)) // 倒序
		}
		return int64(num)
	}
}
