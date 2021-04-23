package rdb

import "fmt"

func rdbSegmentPrint() {
	fmt.Println("--------------------------------------------------------------------------------")
}

func rdbConvertHeaderPrint() {
	fmt.Printf("%-40s  %-40s\n", "rdb raw in hex", "means")
}

func rdbConvertPrint(color string, raw []byte, msg string) {
	if color == "red" {
		fmt.Printf("\033[31m")
	} else if color == "green" {
		fmt.Printf("\033[32m")
	} else if color == "yellow" {
		fmt.Printf("\033[33m")
	} else if color == "cyan" {
		fmt.Printf("\033[36m")
	} else {
	}

	for _, byte := range raw {
		fmt.Printf("%02X ", byte)
	}
	for i := 0; i < 42-len(raw)*3; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("%s\033[0m\n", msg)
}
