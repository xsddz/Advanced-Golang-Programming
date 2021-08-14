package foobar

// ShiftZero
// 请你补全下⾯的代码，输出两数之和，使下⾯的代码能正确运⾏
// 必须定义⼀个包名为 `main` 的包
// 9, 2, 3, 0, 3, 2, 9, 2, 0, 2, 2, 0, 2, 3
//将所有0移到数组末尾  O(1)空间  O(n)时间  非0稳定
func ShiftZero(nums []int) {
	// 在这⾥书写你的代码
	numsLen := len(nums)
	for currentZeroIndex, i := 0, 0; i < numsLen; i++ {
		if nums[i] == 0 {
			continue
		}

		nums[currentZeroIndex], nums[i] = nums[i], nums[currentZeroIndex]
		currentZeroIndex++
	}
}

// RemoveXian 设计一个函数，实现以下功能：将字符串 "li_xiang-qi_che" 转换为"LiXiangQiChe"
func RemoveXian(s string) string {
	ns, needtoupper := []byte{}, false
	for i := 0; i < len(s); i++ {
		if i == 0 {
			needtoupper = true
		}
		if s[i] == '_' || s[i] == '-' {
			needtoupper = true
			continue
		}

		c := s[i]
		if needtoupper {
			c = c - 32
			needtoupper = false
		}
		ns = append(ns, c)
	}
	return string(ns)
}
