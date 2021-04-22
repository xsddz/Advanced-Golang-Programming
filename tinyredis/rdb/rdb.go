package rdb

import (
	"fmt"
	"io"
)

// RDBLoad -
func RDBLoad(filename string) {
	rdb := newRio(filename)
	defer rdb.releaseRio()

	rdbLoadRio(rdb)
}

func rdbLoadRio(rdb *rio) {
	rdbSegmentPrint()
	rdbConvertHeaderPrint()

	// segment:: magic string
	rdbSegmentPrint()
	_, rdbver := rdb.rdbLoadMagicString()

	// segment:: data
	rdbSegmentPrint()
	for {
		// read type of operation
		optype, err := rdb.rdbLoadType()
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadRio] rdbLoadType err:", err))
		}

		// find a handler and excute
		handler, exist := opTypeHandlerMap[optype]
		if !exist {
			panic(fmt.Sprintln("[rdbLoadRio] rdbLoadType err:", err))
		}
		err = handler(rdb)
		if err != nil { // include io.EOF
			if err != io.EOF {
				panic(fmt.Sprintln("[rdbLoadRio] handler() err:", err))
			}
			break
		}
	}

	// segment:: crc 64 checksum
	rdbSegmentPrint()
	rdb.rdbLoadCRC64Checksum(rdbver)

	rdbSegmentPrint()
}
