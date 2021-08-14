package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"tinyredis/lib/rdb"
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
