package foobar

import (
	"fmt"
	"strconv"
)

func fullAdder(a int, b int, c int, base int) (sum int, carry int) {
	s := a + b + c
	sum = s % base
	carry = s / base
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

// BigAdder 大整数加法
func BigAdder(x string, y string) string {
	return bigMaker(fullAdder, 10, x, y)
}

/*****************/

func halfMult(a int, b int, c int, base int) (sum int, carry int) {
	s := a*b + c
	sum = s % base
	carry = s / base
	return
}

func fullMult(a int, y string, rightZeroNum int) (sum string) {
	v, carry := 0, 0
	for i := len(y) - 1; i >= 0; i-- {
		vy, _ := strconv.ParseInt(string(y[i]), 10, 32)
		v, carry = halfMult(a, int(vy), carry, 10)
		sum = fmt.Sprintf("%d%s", v, sum)
	}
	if carry > 0 {
		sum = fmt.Sprintf("%d%s", carry, sum)
	}
	for i := 0; i < rightZeroNum; i++ {
		sum = fmt.Sprintf("%s%d", sum, 0)
	}
	return sum
}

// BigMult 大整数乘法
func BigMult(x string, y string) string {
	sumCollect := []string{}
	for i := len(x) - 1; i >= 0; i-- {
		vx, _ := strconv.ParseInt(string(x[i]), 10, 32)
		sumCollect = append(sumCollect, fullMult(int(vx), y, len(x)-i-1))
	}

	sum := "0"
	for _, item := range sumCollect {
		sum = BigAdder(sum, item)
	}

	return sum
}
