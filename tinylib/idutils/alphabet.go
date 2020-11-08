package idutils

var (
	// Alpha34 -
	Alpha34 = "123456789abcdefghijkmnopqrstuvwxyz"
	// Alpha58 -
	Alpha58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// AlphabetShuffle - 根据salt将字母表打乱
func AlphabetShuffle(alphabet string, salt string) string {
	// 未设置salt时，不打乱字母表
	if salt == "" {
		return alphabet
	}

	alphabetBytes := []byte(alphabet)
	alphabetLen := len(alphabetBytes)
	saltBytes := []byte(salt)
	saltLen := len(salt)
	for i, v, p := alphabetLen-1, 0, 0; i >= 1; i-- {
		intergerVal := 0
		v = v % saltLen
		intergerVal = int(saltBytes[v])
		p += intergerVal
		j := (intergerVal + v + p) % i

		alphabetBytes[i], alphabetBytes[j] = alphabetBytes[j], alphabetBytes[i]

		v++
	}
	return string(alphabetBytes)
}
