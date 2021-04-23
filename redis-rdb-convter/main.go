package main

import (
	"Advanced-Golang-Programming/tinyredislib/rdb"
	"flag"
	"fmt"
	"path/filepath"
)

var (
	rdbfile string
)

func init() {
	flag.StringVar(&rdbfile, "d", "", "rdb file")
	flag.Parse()
}

func main() {
	file, _ := filepath.Abs(rdbfile)
	fmt.Println("convert file:", file)
	rdb.RDBLoad(file)
}
