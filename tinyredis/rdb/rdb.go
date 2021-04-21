package rdb

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

const (
	/* rdbLoad...() functions flags. */
	RDB_LOAD_NONE  int = 0
	RDB_LOAD_ENC       = (1 << 0)
	RDB_LOAD_PLAIN     = (1 << 1)
	RDB_LOAD_SDS       = (1 << 2)

	/* When a length of a string object stored on disk has the first two bits
	 * set, the remaining six bits specify a special encoding for the object
	 * accordingly to the following defines: */
	RDB_ENC_INT8  uint64 = 0 /* 8 bit signed integer */
	RDB_ENC_INT16        = 1 /* 16 bit signed integer */
	RDB_ENC_INT32        = 2 /* 32 bit signed integer */
	RDB_ENC_LZF          = 3 /* string compressed with FASTLZ */

	/* Defines related to the dump file format. To store 32 bits lengths for short
	 * keys requires a lot of space, so we check the most significant 2 bits of
	 * the first byte to interpreter the length:
	 *
	 * 00|XXXXXX => if the two MSB are 00 the len is the 6 bits of this byte
	 * 01|XXXXXX XXXXXXXX =>  01, the len is 14 byes, 6 bits + 8 bits of next byte
	 * 10|000000 [32 bit integer] => A full 32 bit len in net byte order will follow
	 * 10|000001 [64 bit integer] => A full 64 bit len in net byte order will follow
	 * 11|OBKIND this means: specially encoded object will follow. The six bits
	 *           number specify the kind of object that follows.
	 *           See the RDB_ENC_* defines.
	 *
	 * Lengths up to 63 are stored using a single byte, most DB keys, and may
	 * values, will fit inside. */
	RDB_6BITLEN  byte = 0
	RDB_14BITLEN      = 1
	RDB_32BITLEN      = 0x80
	RDB_64BITLEN      = 0x81
	RDB_ENCVAL        = 3
	// RDB_LENERR   = UINT64_MAX

	/* Special RDB opcodes (saved/loaded with rdbSaveType/rdbLoadType). */
	RDB_OPCODE_MODULE_AUX    byte = 247 /* Module auxiliary data. */
	RDB_OPCODE_IDLE               = 248 /* LRU idle time. */
	RDB_OPCODE_FREQ               = 249 /* LFU frequency. */
	RDB_OPCODE_AUX                = 250 /* RDB aux field. */
	RDB_OPCODE_RESIZEDB           = 251 /* Hash table resize hint. */
	RDB_OPCODE_EXPIRETIME_MS      = 252 /* Expire time in milliseconds. */
	RDB_OPCODE_EXPIRETIME         = 253 /* Old expire time in seconds. */
	RDB_OPCODE_SELECTDB           = 254 /* DB number of the following keys. */
	RDB_OPCODE_EOF                = 255 /* End of the RDB file. */

	/* Map object types to RDB object types. Macros starting with OBJ_ are for
	 * memory storage and may change. Instead RDB types must be fixed because
	 * we store them on disk. */
	RDB_TYPE_STRING   byte = 0
	RDB_TYPE_LIST          = 1
	RDB_TYPE_SET           = 2
	RDB_TYPE_ZSET          = 3
	RDB_TYPE_HASH          = 4
	RDB_TYPE_ZSET_2        = 5 /* ZSET version 2 with doubles stored in binary. */
	RDB_TYPE_MODULE        = 6
	RDB_TYPE_MODULE_2      = 7 /* Module value with annotations for parsing without
	   the generating module being loaded. */
/* NOTE: WHEN ADDING NEW RDB TYPE, UPDATE rdbIsObjectType() BELOW */
)

func rdbConvertPrint(color string, raw []byte, msg string) {
	if color == "red" {
		fmt.Printf("\x1B[31m")
	} else if color == "green" {
		fmt.Printf("\x1B[32m")
	} else if color == "yellow" {
		fmt.Printf("\x1B[33m")
	} else {
	}

	if len(raw) > 0 {
		fmt.Printf("[rdb raw in hex]:")
	} else {
		fmt.Printf("                 ")
	}

	for _, byte := range raw {
		fmt.Printf(" %02X", byte)
	}
	for i := 0; i < 34-len(raw)*3; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("%s\x1B[0m\n", msg)
}

// RDBLoad -
func RDBLoad(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("open rdb file err:", err)
		return
	}
	defer file.Close()

	rdb := &rio{
		fp: file,
	}
	rdbLoadRio(rdb)
}

func rdbLoadRio(rdb *rio) {
	// read magic string and rdb version
	buf := make([]byte, 9)
	err := rioRead(rdb, buf)
	if err != nil {
		fmt.Println("[rdbLoadRio] read rdb file err:", err)
		return
	}
	rdbConvertPrint("green", buf, fmt.Sprintf("=> [rdbLoadRio] magic string and rdb version: %v", string(buf)))
	rdbver, _ := strconv.Atoi(string(buf[6:]))

	for {
		// read opcode
		opcode, err := rdbLoadType(rdb)
		if err != nil {
			break
		}

		if opcode == RDB_OPCODE_AUX {
			/* AUX: generic string-string fields. Use to add state to RDB
			 * which is backward compatible. Implementations of RDB loading
			 * are requierd to skip AUX fields they don't understand.
			 *
			 * An AUX field is composed of two strings: key and value. */
			rdbLoadStringObject(rdb) // read key
			rdbLoadStringObject(rdb) // read val
			continue                 /* Read next opcode. */
		} else if opcode == RDB_OPCODE_EOF {
			/* EOF: End of file, exit the main loop. */
			break
		} else if opcode == RDB_OPCODE_SELECTDB {
			/* SELECTDB: Select the specified database. */
			dbid, _, err := rdbLoadLen(rdb)
			if err != nil {
				break
			}
			rdbConvertPrint("green", []byte{}, fmt.Sprintf("=> [rdbLoadRio] select db: %d", dbid))
			continue /* Read next opcode. */
		} else if opcode == RDB_OPCODE_RESIZEDB {
			/* RESIZEDB: Hint about the size of the keys in the currently
			 * selected data base, in order to avoid useless rehashing. */
			dbSize, _, err := rdbLoadLen(rdb)
			if err != nil {
				break
			}
			expireSize, _, err := rdbLoadLen(rdb)
			if err != nil {
				break
			}
			rdbConvertPrint("", []byte{}, fmt.Sprintf("=> [rdbLoadRio] db_size:%d, expire_size: %d", dbSize, expireSize))
			continue /* Read next opcode. */
		}

		rdbLoadStringObject(rdb)   // read key
		rdbLoadObject(opcode, rdb) // read obj
	}

	if rdbver > 5 {
		buf := make([]byte, 8)
		err = rioRead(rdb, buf)
		if err != nil {
			fmt.Println("[rdbLoadRio] read rdb file err:", err)
			return
		}
		rdbConvertPrint("green", buf, fmt.Sprintf("=> [rdbLoadRio] CRC 64 checksum: %08b", buf))
	}
}

func rdbLoadType(rdb *rio) (byte, error) {
	buf := make([]byte, 1)
	err := rioRead(rdb, buf)
	if err != nil {
		fmt.Println("[rdbLoadType] read rdb file err:", err)
		return 0, err
	}
	rdbConvertPrint("red", buf, fmt.Sprintf("=> [rdbLoadType] rdb opcode: %X, %d", buf[0], buf[0]))
	return buf[0], nil
}

func rdbLoadStringObject(rdb *rio) {
	rdbGenericLoadStringObject(rdb, RDB_LOAD_NONE)
}

func rdbLoadEncodedStringObject(rdb *rio) {
	rdbGenericLoadStringObject(rdb, RDB_LOAD_ENC)
}

func rdbLoadObject(vType byte, rdb *rio) {
	if vType == RDB_TYPE_STRING {
		rdbLoadEncodedStringObject(rdb)
	} else if vType == RDB_TYPE_LIST {
	} else if vType == RDB_TYPE_SET {
	} else if vType == RDB_TYPE_ZSET_2 || vType == RDB_TYPE_ZSET {
	} else if vType == RDB_TYPE_HASH {
	} else {
	}
}

func rdbGenericLoadStringObject(rdb *rio, flags int) {
	len, isEnc, err := rdbLoadLen(rdb)
	if err != nil {
		fmt.Println("[rdbGenericLoadStringObject] rdbLoadLen() err:", err)
		return
	}
	if isEnc {
		buf := make([]byte, 4)
		if len == RDB_ENC_INT8 {
			err = rioRead(rdb, buf[:1])
			if err != nil {
				fmt.Println("[rdbGenericLoadStringObject] rdbLoadLen() err:", err)
				return
			}
			val := int(buf[0])
			rdbConvertPrint("yellow", buf[:1], fmt.Sprintf("=> [rdbGenericLoadStringObject] encode int8: %08b, %v", buf[:1], val))
			return
		} else if len == RDB_ENC_INT16 {
			err = rioRead(rdb, buf[:2])
			if err != nil {
				fmt.Println("[rdbGenericLoadStringObject] rdbLoadLen() err:", err)
				return
			}
			val := int16(buf[0]) | (int16(buf[1]) << 8)
			rdbConvertPrint("yellow", buf[:2], fmt.Sprintf("=> [rdbGenericLoadStringObject] encode int16: %08b, %v", buf[:2], val))
			return
		} else if len == RDB_ENC_INT32 {
			err = rioRead(rdb, buf)
			if err != nil {
				fmt.Println("[rdbGenericLoadStringObject] rdbLoadLen() err:", err)
				return
			}
			val := int32(buf[0]) | (int32(buf[1]) << 8) | (int32(buf[2]) << 16) | (int32(buf[3]) << 24)
			rdbConvertPrint("yellow", buf, fmt.Sprintf("=> [rdbGenericLoadStringObject] encode int32: %08b, %v", buf, val))
			return
		} else if len == RDB_ENC_LZF {
			clen, _, err := rdbLoadLen(rdb)
			if err != nil {
				fmt.Println("[rdbGenericLoadStringObject] rdbLoadLen() err:", err)
				return
			}
			len, _, err := rdbLoadLen(rdb)
			if err != nil {
				fmt.Println("[rdbGenericLoadStringObject] rdbLoadLen() err:", err)
				return
			}
			buf := make([]byte, clen)
			err = rioRead(rdb, buf)
			if err != nil {
				fmt.Println("[rdbGenericLoadStringObject] rdbLoadLen() err:", err)
				return
			}
			rdbConvertPrint("yellow", buf, fmt.Sprintf("=> [rdbGenericLoadStringObject] encode LZF compressed string: %d, %v, %v", len, buf, string(buf)))
			return
		} else {
			fmt.Printf("[rdbGenericLoadStringObject] unknown RDB string encoding type %d", len)
			return
		}
	}

	buf := make([]byte, len)
	err = rioRead(rdb, buf)
	if err != nil {
		fmt.Println("[rdbGenericLoadStringObject] read rdb file err:", err)
		return
	}
	rdbConvertPrint("yellow", buf, fmt.Sprintf("=> [rdbGenericLoadStringObject]: %v", string(buf)))
	return
}

func rdbLoadLen(rdb *rio) (len uint64, isEnc bool, err error) {
	buf := make([]byte, 2)
	err = rioRead(rdb, buf[:1])
	if err != nil {
		fmt.Println("[rdbLoadLen] read rdb file err:", err)
		return 0, false, err
	}
	// rdbConvertPrint("", buf[:1], fmt.Sprintf("=> [rdbLoadLen] type: %08b", buf[0]))

	lenType := (buf[0] & 0xC0) >> 6
	if lenType == RDB_ENCVAL {
		/* Read a 6 bit encoding type. */
		isEnc = true
		len = uint64(buf[0] & 0x3F)
		rdbConvertPrint("", buf[:1], fmt.Sprintf("=> [rdbLoadLen] 6 bit encoding type: %08b, %v", buf[:1], len))
		return
	} else if lenType == RDB_6BITLEN {
		/* Read a 6 bit len. */
		len = uint64(buf[0] & 0x3F)
		rdbConvertPrint("", buf[:1], fmt.Sprintf("=> [rdbLoadLen] 6 bit len: %08b, %v", buf[:1], len))
		return
	} else if lenType == RDB_14BITLEN {
		/* Read a 14 bit len. */
		err = rioRead(rdb, buf[1:])
		if err != nil {
			fmt.Println("[rdbLoadLen] read rdb file err:", err)
			return 0, false, err
		}
		len = uint64(((uint(buf[0]) & 0x3F) << 8) | uint(buf[1]))
		rdbConvertPrint("", buf, fmt.Sprintf("=> [rdbLoadLen] 14 bit len: %08b, %v", buf, len))
		return
	} else if buf[0] == RDB_32BITLEN {
		/* Read a 32 bit len. */
		buf32 := make([]byte, 4)
		err = rioRead(rdb, buf32)
		if err != nil {
			fmt.Println("[rdbLoadLen] read rdb file err:", err)
			return 0, false, err
		}

		// network byte ordering => big endian systems
		len = uint64(binary.BigEndian.Uint32(buf32))

		rdbConvertPrint("", buf32, fmt.Sprintf("=> [rdbLoadLen] 32 bit len: %08b, %032b, %d", buf32, len, len))
		return
	} else if buf[0] == RDB_64BITLEN {
		/* Read a 64 bit len. */
		buf64 := make([]byte, 8)
		err = rioRead(rdb, buf64)
		if err != nil {
			fmt.Println("[rdbLoadLen] read rdb file err:", err)
			return 0, false, err
		}

		// network byte ordering => big endian systems
		len = binary.BigEndian.Uint64(buf64)

		rdbConvertPrint("", buf64, fmt.Sprintf("=> [rdbLoadLen] 64 bit len: %08b, %032b, %d", buf64, len, len))

		return
	} else {
		return 0, false, fmt.Errorf("unknow length encoding: %02b", lenType)
	}
}
