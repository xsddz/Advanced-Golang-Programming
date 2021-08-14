package foobar

// [1 2 3 4 5] =>1 [2 3 4 5 1]
// [1 2 3 4 5] =>2 [3 4 5 1 2]  4
// [1 2 3 4 5] =>3 [4 5 1 2 3]  4
// [1 2 3 4 5]  4

// FindX -
func FindX(nums []int, x int) int {
	return bs(nums, 0, len(nums), x)
}

func bs(nums []int, l int, r int, x int) int {
	if l > r {
		return -1
	}

	mid := (l + r) / 2
	if mid == r {
		return -1
	}

	if nums[mid] < x {
		findIndex := bs(nums, mid+1, r, x)
		if findIndex == -1 {
			findIndex = bs(nums, l, mid, x)
		}
		return findIndex
	} else if nums[mid] > x {
		findIndex := bs(nums, l, mid, x)
		if findIndex == -1 {
			findIndex = bs(nums, mid+1, r, x)
		}
		return findIndex
	} else {
		return mid
	}
}
