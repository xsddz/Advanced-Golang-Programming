package foobar

import "fmt"

func init() {
	fmt.Println("init:foobar")
}

func sortProc(data []int, l int, r int) (mid int) {
	midVal := data[r]
	i := l - 1
	for j := i + 1; j < r; j++ {
		if data[j] < midVal {
			i++
			data[i], data[j] = data[j], data[i]
		}
	}
	data[r], data[i+1] = data[i+1], midVal
	return i + 1
}

func sortIter(data []int, l int, r int) {
	if l < r {
		mid := sortProc(data, l, r)
		sortIter(data, l, mid-1)
		sortIter(data, mid+1, r)
	}
}

// Sort -
func Sort(data []int) {
	dataLen := len(data)
	if dataLen < 2 {
		return
	}

	sortIter(data, 0, dataLen-1)
}
