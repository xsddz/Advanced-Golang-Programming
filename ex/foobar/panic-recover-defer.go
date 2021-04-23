package foobar

import (
	"errors"
	"fmt"
)

func PRD() {
	defer func() {
		fmt.Println("====PRD")
		if r := recover(); r != nil {
			fmt.Printf("-----recover: %T, %v\n", r, r)
		}
	}()

	panic(errors.New("bbb"))
}
