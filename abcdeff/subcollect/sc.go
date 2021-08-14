package subcollect

import (
	"fmt"
)

/*
两个整型数组，需要一个函数，判断小数组是否是大数组的连续子集，如果是返回下标，否则返回-1
举例：
A={1,2,3,1,2,4,5},  B={4,5} return 5
A={1,2,3,1,2,4,5},  B={1,2} return 0
A={1,2,3,1,2,4,5},  B={1,4,5} return -1
*/

func isSub(a []int, b []int) int {
	lenA, lenB := len(a), len(b)
	if lenA < lenB {
		return -1
	}
	if lenA == lenB {
		if a[0] == b[0] {
			return 0
		}
		return -1
	}

	eqIndex := 0
	for {
		eqIndex = findFirst(a, eqIndex, b[0])
		if eqIndex == -1 {
			return -1
		}

		if eq(a, eqIndex, b) {
			return eqIndex
		}
		eqIndex++
	}
}

func eq(a []int, index int, b []int) bool {
	for i, j := index, 0; j < len(b); i, j = i+1, j+1 {
		if a[i] != b[j] {
			return false
		}
	}
	return true
}

func findFirst(a []int, startIndex int, v int) int {
	for i := startIndex; i < len(a); i++ {
		if a[i] == v {
			return i
		}
	}
	return -1
}

// Run -
func Run() {
	a, b := []int{1, 2, 3, 1, 2, 4, 5}, []int{4, 5}
	fmt.Println(a, b, isSub(a, b))

	a, b = []int{1, 2, 3, 1, 2, 4, 5}, []int{1, 2}
	fmt.Println(a, b, isSub(a, b))

	a, b = []int{1, 2, 3, 1, 2, 4, 5}, []int{1, 4, 5}
	fmt.Println(a, b, isSub(a, b))
}
