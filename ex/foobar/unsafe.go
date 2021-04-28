package foobar

import (
	"fmt"
	"unsafe"
)

var foo struct{}

var caseA struct {
	a bool
	b int16
	c int
}
var caseB struct {
	a bool
	c int
	b int16
}
var caseC struct {
	c int
	a bool
	b int16
}
var caseD struct {
	c int
	b int16
	a bool
}

// MachineEndian -
func MachineEndian() {
	a := int16(0x1234)
	b := int8(a)
	if b == 0x34 {
		fmt.Println("little endian.")
	} else {
		fmt.Println("big endian.")
	}
}

// MachineBit -
func MachineBit() int {
	bit := 32 << (^uint(0) >> 63)
	fmt.Printf("machine bit: %v\n", bit)
	return bit
}

// USMem -
func USMem() {
	fmt.Printf("====size of struct{}: %v\n", unsafe.Sizeof(foo))

	sizeX, alignX := unsafe.Sizeof(caseA), unsafe.Alignof(caseA)
	fmt.Printf("====typeX: %T, sizeX: %v, alignX: %v\n", caseA, sizeX, alignX)
	sizeXA, alignXA, offsetXA := unsafe.Sizeof(caseA.a), unsafe.Alignof(caseA.a), unsafe.Offsetof(caseA.a)
	fmt.Printf("offsetXA: %v, sizeXA: %v, alignXA: %v\n", sizeXA, alignXA, offsetXA)
	sizeXB, alignXB, offsetXB := unsafe.Sizeof(caseA.b), unsafe.Alignof(caseA.b), unsafe.Offsetof(caseA.b)
	fmt.Printf("offsetXB: %v, sizeXB: %v, alignXB: %v\n", sizeXB, alignXB, offsetXB)
	sizeXC, alignXC, offsetXC := unsafe.Sizeof(caseA.c), unsafe.Alignof(caseA.c), unsafe.Offsetof(caseA.c)
	fmt.Printf("offsetXC: %v, sizeXC: %v, alignXC: %v\n", sizeXC, alignXC, offsetXC)

	sizeX, alignX = unsafe.Sizeof(caseB), unsafe.Alignof(caseB)
	fmt.Printf("====typeX: %T, sizeX: %v, alignX: %v\n", caseB, sizeX, alignX)
	sizeXA, alignXA, offsetXA = unsafe.Sizeof(caseB.a), unsafe.Alignof(caseB.a), unsafe.Offsetof(caseB.a)
	fmt.Printf("offsetXA: %v, sizeXA: %v, alignXA: %v\n", sizeXA, alignXA, offsetXA)
	sizeXB, alignXB, offsetXB = unsafe.Sizeof(caseB.b), unsafe.Alignof(caseB.b), unsafe.Offsetof(caseB.b)
	fmt.Printf("offsetXB: %v, sizeXB: %v, alignXB: %v\n", sizeXB, alignXB, offsetXB)
	sizeXC, alignXC, offsetXC = unsafe.Sizeof(caseB.c), unsafe.Alignof(caseB.c), unsafe.Offsetof(caseB.c)
	fmt.Printf("offsetXC: %v, sizeXC: %v, alignXC: %v\n", sizeXC, alignXC, offsetXC)

	sizeX, alignX = unsafe.Sizeof(caseC), unsafe.Alignof(caseC)
	fmt.Printf("====typeX: %T, sizeX: %v, alignX: %v\n", caseC, sizeX, alignX)
	sizeXA, alignXA, offsetXA = unsafe.Sizeof(caseC.a), unsafe.Alignof(caseC.a), unsafe.Offsetof(caseC.a)
	fmt.Printf("offsetXA: %v, sizeXA: %v, alignXA: %v\n", sizeXA, alignXA, offsetXA)
	sizeXB, alignXB, offsetXB = unsafe.Sizeof(caseC.b), unsafe.Alignof(caseC.b), unsafe.Offsetof(caseC.b)
	fmt.Printf("offsetXB: %v, sizeXB: %v, alignXB: %v\n", sizeXB, alignXB, offsetXB)
	sizeXC, alignXC, offsetXC = unsafe.Sizeof(caseC.c), unsafe.Alignof(caseC.c), unsafe.Offsetof(caseC.c)
	fmt.Printf("offsetXC: %v, sizeXC: %v, alignXC: %v\n", sizeXC, alignXC, offsetXC)

	sizeX, alignX = unsafe.Sizeof(caseD), unsafe.Alignof(caseD)
	fmt.Printf("====typeX: %T, sizeX: %v, alignX: %v\n", caseD, sizeX, alignX)
	sizeXA, alignXA, offsetXA = unsafe.Sizeof(caseD.a), unsafe.Alignof(caseD.a), unsafe.Offsetof(caseD.a)
	fmt.Printf("offsetXA: %v, sizeXA: %v, alignXA: %v\n", sizeXA, alignXA, offsetXA)
	sizeXB, alignXB, offsetXB = unsafe.Sizeof(caseD.b), unsafe.Alignof(caseD.b), unsafe.Offsetof(caseD.b)
	fmt.Printf("offsetXB: %v, sizeXB: %v, alignXB: %v\n", sizeXB, alignXB, offsetXB)
	sizeXC, alignXC, offsetXC = unsafe.Sizeof(caseD.c), unsafe.Alignof(caseD.c), unsafe.Offsetof(caseD.c)
	fmt.Printf("offsetXC: %v, sizeXC: %v, alignXC: %v\n", sizeXC, alignXC, offsetXC)
}

var bgColorMap = map[string]string{
	"Black":   "\033[40m",
	"Red":     "\033[41m",
	"Green":   "\033[42m",
	"Yellow":  "\033[43m",
	"Blue":    "\033[44m",
	"Magenta": "\033[45m",
	"Cyan":    "\033[46m",
	"White":   "\033[47m",
}

func colorBlock(color string, c string) {
	colorCtl, ok := bgColorMap[color]
	if !ok {
		colorCtl = ""
	}
	fmt.Printf(colorCtl+"%s\033[0m", c)
}
