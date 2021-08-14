package foobar

// CollectBySlice -
func CollectBySlice(collect *[]int, n int) {
	if n == 0 {
		return
	}
	*collect = append(*collect, n)
	CollectBySlice(collect, n-1)
}
