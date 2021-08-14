package main

import (
	"abcdeff/foobar"
	"errors"
	"fmt"
)

func init() {
	fmt.Println("init:main")
}

func aaa() (err error) {
	defer func() {
		fmt.Println("defer err:", err)
	}()

	return errors.New("测试")
}

func main() {
	foobar.MachineEndian()
	foobar.MachineBit()
	foobar.USMem()

	aaa()
}
