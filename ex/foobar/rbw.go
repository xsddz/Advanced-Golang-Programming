package foobar

/*
题目描述
给定一个包含红色、白色和蓝色，一共 n 个元素的数组，原地对它们进行排序，使得相同颜色的元素相邻，并按照红色、白色、蓝色顺序排列。
此题中，我们使用整数 0、 1 和 2 分别表示红色、白色和蓝色。
注意:
不能使用代码库中的排序函数来解决这道题。

样例
输入
[2,0,2,1,1,0]
输出
[0,0,1,1,2,2]
*/

// RWB -
func RWB(a []int) {
	ri, bi := -1, len(a)
	for i := 0; i < bi; {
		if a[i] == 0 {
			ri++
			a[ri], a[i] = a[i], a[ri]

			i++
			continue
		}
		if a[i] == 2 {
			bi--
			a[bi], a[i] = a[i], a[bi]
			continue
		}
		if a[i] == 1 {
			i++
		}
	}
}