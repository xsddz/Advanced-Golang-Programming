package foobar

import "fmt"

func getNext(p string) []int {
	i, j, next := 1, 0, make([]int, len(p))
	for i < len(p) {
		if p[i] == p[j] {
			j++
			i++
			next[i] = j
		} else {
			j = next[i]
		}
	}
	return next
}

func KMP() {
	fmt.Println("====", getNext("abababca"))
}
