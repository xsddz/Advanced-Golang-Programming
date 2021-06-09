package foobar

import "fmt"

func ReturnDefer(nums []int) {
	defer func() {
		fmt.Println("======defer before:", nums)
		nums[0] = 1223
		fmt.Println("======defer  after:", nums)
	}()
	nums[0] = 1111
}
