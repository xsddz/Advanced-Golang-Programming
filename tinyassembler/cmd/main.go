package main

import (
	"Advanced-Golang-Programming/tinyassembler"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: tinyassembler filename.asm\n\n")
		return
	}

	targetFileDir := filepath.Dir(os.Args[1])
	sourceFile := filepath.Base(os.Args[1])
	targetFile := strings.TrimSuffix(sourceFile, filepath.Ext(sourceFile)) + ".hack"

	tinyassembler.Run(os.Args[1], filepath.Join(targetFileDir, targetFile))
}
