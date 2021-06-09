package foobar

// ReturnDefer -
func ReturnDefer() int {
	a := 3
	defer func() {
		a = 4
	}()
	return a
}
