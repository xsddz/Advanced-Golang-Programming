package foobar

/*
给一个数组，有红白蓝三个颜色，分别代表0，1，2；
原地排序这个数组（不允许用额外的空间），使得最终的数组
- 最前面是红色，
- 中间是白色，
- 最后是蓝色
$input = [2,0,2,1,1,0];
$output = [0,0,1,1,2,2];
*/

// RWB -
func RWB(a []int) {
	ri, bi := -1, len(a)
	for i := 0; i < bi; i++ {
		if a[i] == 0 {
			ri++
			a[ri], a[i] = a[i], a[ri]
		}
		if a[i] == 2 {
			bi--
			a[bi], a[i] = a[i], a[bi]
		}
	}
}
