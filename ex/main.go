package main

import (
	"Advanced-Golang-Programming/ex/foobar"
	"fmt"
)

func init() {
	fmt.Println("init:main")
}

func main() {
	defer func() {
		fmt.Println("=====main")
		if r := recover(); r != nil {
			fmt.Printf("++++recover: %T, %v\n", r, r)
		}
	}()

	foobar.PRD()

	for {
	}
}
