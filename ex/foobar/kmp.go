package foobar

import "fmt"

func getPMT(p string) []int {
	j, k, next := 1, 0, make([]int, len(p)+1)
	next[0] = -1
	for j < len(p) {
		if k == -1 || p[j] == p[k] {
			j++
			k++
			next[j] = k
		} else {
			k = next[k]
		}
	}
	return next
}

func getNext(p string) []int {
	j, k, next := 1, 0, make([]int, len(p)+1)
	next[0] = -1
	for j < len(p) {
		if k == -1 || p[j] == p[k] {
			j++
			k++
			next[j] = k
		} else {
			k = next[k]
		}
	}
	return next
}

// KMP -
func KMP() {
	// s, p := "ABC ABCDAB ABCDABCDABDE", "ABCDABD"
	s, p := "aabaaabaaac", "aabaaac"
	fmt.Println("=======s:", s)
	fmt.Println("=======p:", p)
	fmt.Println("=====pmt:", getPMT(p))
	fmt.Println("====next:", getNext(p))

	// si, sl, pi, pl, nextP := 0, len(s), 0, len(p), getNext(p)
	// for si < sl && pi < pl {
	// 	if pi == -1 || s[si] == p[pi] {
	// 		si++
	// 		pi++
	// 	} else {
	// 		pi = nextP[pi]
	// 	}
	// }
	// if pi == pl {
	// 	fmt.Println("found:", si-pl)
	// } else {
	// 	fmt.Println("not found")
	// }
}
