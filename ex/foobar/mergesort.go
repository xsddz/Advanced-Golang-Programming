package foobar

// MergeSort -
func MergeSort(nums []int) {
	msIter(nums, 0, len(nums))
}

func msIter(nums []int, l int, r int) {
	if len(nums[l:r]) < 2 {
		return
	}

	mid := (l + r) / 2
	msIter(nums, l, mid)
	msIter(nums, mid, r)
	msProc(nums, l, mid, r)
}

func msProc(nums []int, l int, mid int, r int) {
	// fmt.Println("++++proc:", l, mid, r, nums[l:mid], nums[mid:r])
	i, j := l, mid
	for i < mid && j < r {
		if nums[i] <= nums[j] {
			i++
		} else {
			// 取出num[j]，num[i:j-1]整体后移一位
			temp := nums[j]
			for k := j; k > i; k-- {
				nums[k] = nums[k-1]
			}
			nums[i] = temp

			mid++ // mid 也相应后移一位

			i++
			j++
		}
	}
	// fmt.Println("+++++++proc:", i, mid, j, r)
}
