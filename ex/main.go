package main

import (
	"Advanced-Golang-Programming/ex/foobar"
	"fmt"
)

func init() {
	fmt.Println("init:main")
}

func main() {
	foobar.MachineEndian()
	foobar.MachineBit()
	foobar.USMem()
}
