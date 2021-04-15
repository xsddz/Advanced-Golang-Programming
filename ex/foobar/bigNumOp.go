package foobar

import (
	"fmt"
	"strconv"
)

func halfAdder(a int, b int, base int) (sum int, carry int) {
	c := a + b
	sum = c % base
	carry = c / base
	return
}

func fullAdder(a int, b int, c int, base int) (sum int, carry int) {
	sab, cab := halfAdder(a, b, base)
	sum, cabc := halfAdder(sab, c, base)
	carry = cab + cabc
	return
}

type opCall func(a int, b int, c int, base int) (sum int, carry int)

func bigMaker(op opCall, base int, x string, y string) string {
	rtn, carry := []int{}, 0

	xLen, yLen := len(x), len(y)
	if xLen < yLen {
		x, y = y, x
		xLen, yLen = yLen, xLen
	}

	for i, j := xLen-1, yLen-1; i >= 0; i, j = i-1, j-1 {
		vx, _ := strconv.ParseInt(string(x[i]), 10, 32)
		sum := 0
		if j < 0 {
			sum, carry = op(int(vx), 0, carry, base)
		} else {
			vy, _ := strconv.ParseInt(string(y[j]), 10, 32)
			sum, carry = op(int(vx), int(vy), carry, base)
		}
		rtn = append(rtn, sum)
	}

	if carry > 0 {
		rtn = append(rtn, carry)
	}

	s := ""
	for i := len(rtn) - 1; i >= 0; i-- {
		s += fmt.Sprint(rtn[i])
	}
	return s
}

// BigAdder -
func BigAdder(x string, y string) string {
	return bigMaker(fullAdder, 10, x, y)
}

func fullMult(a int, b int, c int, base int) (sum int, carry int) {
	mab := a*b + c
	sum = mab % base
	carry = mab / base
	return
}
