package theBugs

import "fmt"

/*
 * 在线网站：https://c.runoob.com/compile/21
 * 下面这段golang程序，大部分情况没有输出
 *
 */

func sortIter(data []int, l int, r int) int {
	midVal := data[r]
	i := l - 1
	for j := i + 1; j < r; j++ {
		if data[j] < midVal {
			i++
			temp := data[j]
			data[j] = data[i]
			data[i] = temp
		}
	}
	data[r] = data[i+1]
	data[i+1] = midVal
	return i + 1
}

func sortProcess(data []int, l int, r int) {
	if l < r {
		midIndex := sortIter(data, l, r)
		sortProcess(data, l, midIndex-1)
		sortProcess(data, midIndex+1, r)
	}
}

func Sort(data []int) {
	sortProcess(data, 0, len(data)-1)
}

func Run() {
	data := []int{12, 22, 33, 7, 9, 45}
	fmt.Println("====data:", data)
	// Sort(data)
	// fmt.Println("====data:", data)
}
