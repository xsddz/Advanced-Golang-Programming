package slice

import "fmt"

func appendSlice(a []int) {
	a = append(a, 21)
}

func caseCopy() {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	b := a[1:3]

	fmt.Println("[caseCopy]=====before a:", a)
	fmt.Println("[caseCopy]=====before b:", b)
	appendSlice(b)
	fmt.Println("[caseCopy]=====after a:", a)
	fmt.Println("[caseCopy]=====after b:", b)
}

func caseCopy2() {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	b := a[1:3:3]

	fmt.Println("[caseCopy2]=====before a:", a)
	fmt.Println("[caseCopy2]=====before b:", b)
	appendSlice(b)
	fmt.Println("[caseCopy2]=====after a:", a)
	fmt.Println("[caseCopy2]=====after b:", b)
}

// Run -
func Run() {
	caseCopy()
	caseCopy2()
}
