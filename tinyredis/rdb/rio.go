package rdb

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

/*------------------------------------------------------------------------------
 * 标准定义如下
 * ---------------------------------------------------------------------------*/

type rioI interface {
	// 从rdb文件最多读取buf大小的byte数据
	rioRead(buf []byte) error

	// 从rdb文件头读取 magic string 和 rdb version
	rdbLoadMagicString() (magic string, ver int)

	// 读取类型操作编码
	rdbLoadType() (opType byte, err error)
	// 读取长度编码
	rdbLoadLen() (len uint64, isEnc bool, err error)
	// 读取对象
	rdbLoadStringObject() (val interface{})
	rdbLoadEncodedStringObject() (val interface{})
	rdbGenericLoadStringObject(flags int) interface{}
	rdbLoadIntegerObject(enctype int, flags int) (val interface{})
	rdbLoadLzfStringObject(flags int) (val interface{})
	rdbLoadObject(vType byte) (val interface{})

	rdbLoadTime() (t uint32)
	rdbLoadMillisecondTime(rdbver int) (t uint64)

	// 从rdb文件尾读取校验和
	rdbLoadCRC64Checksum(rdbver int)
}

/*------------------------------------------------------------------------------
 * 实现如下
 * ---------------------------------------------------------------------------*/

const (
	/* rdbLoad...() functions flags. */
	RDB_LOAD_NONE  = 0
	RDB_LOAD_ENC   = (1 << 0)
	RDB_LOAD_PLAIN = (1 << 1)
	RDB_LOAD_SDS   = (1 << 2)

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
	RDB_6BITLEN  = 0
	RDB_14BITLEN = 1
	RDB_32BITLEN = 0x80
	RDB_64BITLEN = 0x81
	RDB_ENCVAL   = 3

	/* When a length of a string object stored on disk has the first two bits
	 * set, the remaining six bits specify a special encoding for the object
	 * accordingly to the following defines: */
	RDB_ENC_INT8  = 0 /* 8 bit signed integer */
	RDB_ENC_INT16 = 1 /* 16 bit signed integer */
	RDB_ENC_INT32 = 2 /* 32 bit signed integer */
	RDB_ENC_LZF   = 3 /* string compressed with FASTLZ */
)

type rio struct {
	rdbver int
	fp     *os.File
}

func newRio(filename string) *rio {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintln("open rdb file err:", err))
	}

	return &rio{
		fp: file,
	}
}

func (r *rio) releaseRio() {
	r.fp.Close()
}

func (r *rio) rioSetRDBVersion(ver int) {
	r.rdbver = ver
}

func (r *rio) rioRead(buf []byte) error {
	// Read reads up to len(buf) bytes from the File.
	// It returns the number of bytes read and any error encountered.
	// At end of file, Read returns 0, io.EOF.
	_, e := r.fp.Read(buf)
	if e != nil {
		return e
	}
	return nil
}

func (r *rio) rdbLoadMagicString() (magic string, ver int) {
	// read magic string and rdb version
	// magic string: 5 byte
	// version: 4 byte
	buf := make([]byte, 9)
	err := r.rioRead(buf)
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadMagicString] read rdb file err:", err))
	}
	rdbConvertPrint("cyan", buf, fmt.Sprintf("=> [rdbLoadMagicString] magic string and rdb version: %v", string(buf)))

	magic = string(buf[:6])
	ver, _ = strconv.Atoi(string(buf[6:]))
	return
}

func (r *rio) rdbLoadType() (byte, error) {
	buf := make([]byte, 1)
	err := r.rioRead(buf)
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadType] read rdb file err:", err))
	}
	rdbConvertPrint("red", buf, fmt.Sprintf("=> [rdbLoadType] rdb optype: %X, %d", buf[0], buf[0]))
	return buf[0], nil
}

func (r *rio) rdbLoadStringObject() interface{} {
	return r.rdbGenericLoadStringObject(RDB_LOAD_NONE)
}

func (r *rio) rdbLoadEncodedStringObject() interface{} {
	return r.rdbGenericLoadStringObject(RDB_LOAD_ENC)
}

func (r *rio) rdbGenericLoadStringObject(flags int) interface{} {
	len, isEnc, err := r.rdbLoadLen()
	if err != nil {
		panic(fmt.Sprintln("[rdbGenericLoadStringObject] rdbLoadLen() err:", err))
	}

	if isEnc {
		if len == RDB_ENC_INT8 ||
			len == RDB_ENC_INT16 ||
			len == RDB_ENC_INT32 {
			return r.rdbLoadIntegerObject(int(len), flags)
		} else if len == RDB_ENC_LZF {
			return r.rdbLoadLzfStringObject(flags)
		} else {
			panic(fmt.Sprintf("[rdbGenericLoadStringObject] unknown RDB string encoding type %d\n", len))
		}
	}

	buf := make([]byte, len)
	err = r.rioRead(buf)
	if err != nil {
		panic(fmt.Sprintln("[rdbGenericLoadStringObject] read rdb file err:", err))
	}
	val := string(buf)
	rdbConvertPrint("", buf, fmt.Sprintf("=> [rdbGenericLoadStringObject]: %v", val))
	return val
}

func (r *rio) rdbLoadLen() (len uint64, isEnc bool, err error) {
	buf := make([]byte, 2)
	err = r.rioRead(buf[:1])
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadLen] read rdb file err:", err))
	}

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
		err = r.rioRead(buf[1:])
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadLen] read rdb file err:", err))
		}
		len = uint64(((uint(buf[0]) & 0x3F) << 8) | uint(buf[1]))
		rdbConvertPrint("", buf, fmt.Sprintf("=> [rdbLoadLen] 14 bit len: %08b, %v", buf, len))
		return
	} else if buf[0] == RDB_32BITLEN {
		/* Read a 32 bit len. */
		buf32 := make([]byte, 4)
		err = r.rioRead(buf32)
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadLen] read rdb file err:", err))
		}

		// network byte ordering => big endian systems
		len = uint64(binary.BigEndian.Uint32(buf32))

		rdbConvertPrint("", buf32, fmt.Sprintf("=> [rdbLoadLen] 32 bit len: %08b, %032b, %d", buf32, len, len))
		return
	} else if buf[0] == RDB_64BITLEN {
		/* Read a 64 bit len. */
		buf64 := make([]byte, 8)
		err = r.rioRead(buf64)
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadLen] read rdb file err:", err))
		}

		// network byte ordering => big endian systems
		len = binary.BigEndian.Uint64(buf64)

		rdbConvertPrint("", buf64, fmt.Sprintf("=> [rdbLoadLen] 64 bit len: %08b, %032b, %d", buf64, len, len))

		return
	} else {
		panic(fmt.Sprintf("[rdbLoadLen] unknow length encoding: %02b, %08b\n", lenType, buf[0]))
	}
}

func (r *rio) rdbLoadIntegerObject(enctype int, flags int) (val interface{}) {
	buf := make([]byte, 4)
	if enctype == RDB_ENC_INT8 {
		err := r.rioRead(buf[:1])
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadIntegerObject] rdbLoadLen() err:", err))
		}
		val = int(buf[0])
		rdbConvertPrint("", buf[:1], fmt.Sprintf("=> [rdbLoadIntegerObject] encode int8: %08b, %v", buf[:1], val))
		return
	} else if enctype == RDB_ENC_INT16 {
		err := r.rioRead(buf[:2])
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadIntegerObject] rdbLoadLen() err:", err))
		}
		val = int16(buf[0]) | (int16(buf[1]) << 8)
		rdbConvertPrint("", buf[:2], fmt.Sprintf("=> [rdbLoadIntegerObject] encode int16: %08b, %v", buf[:2], val))
		return
	} else if enctype == RDB_ENC_INT32 {
		err := r.rioRead(buf)
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadIntegerObject] rdbLoadLen() err:", err))
		}
		val = int32(buf[0]) | (int32(buf[1]) << 8) | (int32(buf[2]) << 16) | (int32(buf[3]) << 24)
		rdbConvertPrint("", buf, fmt.Sprintf("=> [rdbLoadIntegerObject] encode int32: %08b, %v", buf, val))
		return
	} else {
		panic(fmt.Sprintf("[rdbLoadIntegerObject] unknown RDB integer encoding type %d\n", enctype))
	}
}

func (r *rio) rdbLoadLzfStringObject(flags int) interface{} {
	clen, _, err := r.rdbLoadLen()
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadLzfStringObject] rdbLoadLen() err:", err))
	}
	len, _, err := r.rdbLoadLen()
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadLzfStringObject] rdbLoadLen() err:", err))
	}

	buf := make([]byte, clen)
	err = r.rioRead(buf)
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadLzfStringObject] rdbLoadLen() err:", err))
	}

	val := string(buf)
	// todo:: lzf_decompress()
	rdbConvertPrint("", buf, fmt.Sprintf("=> [rdbLoadLzfStringObject] encode LZF compressed string: %d, %v, %v", len, buf, string(buf)))
	return val
}

func (r *rio) rdbLoadObject(vType int) interface{} {
	if vType == RDB_TYPE_STRING {
		/* Read string value */
		return r.rdbLoadEncodedStringObject()
	} else if vType == RDB_TYPE_LIST {
		/* Read list value */
		len, _, err := r.rdbLoadLen()
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadObject] rdbLoadLen() err:", err))
		}
		vals := []interface{}{}
		for i := uint64(0); i < len; i++ {
			vals = append(vals, r.rdbLoadEncodedStringObject())
		}
		return vals
	} else if vType == RDB_TYPE_SET {
		/* Read Set value */
		len, _, err := r.rdbLoadLen()
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadObject] rdbLoadLen() err:", err))
		}
		vals := []interface{}{}
		for i := uint64(0); i < len; i++ {
			vals = append(vals, r.rdbGenericLoadStringObject(RDB_LOAD_SDS))
		}
		return vals
	} else if vType == RDB_TYPE_ZSET_2 || vType == RDB_TYPE_ZSET {
		/* Read list/set value. */
		len, _, err := r.rdbLoadLen()
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadObject] rdbLoadLen() err:", err))
		}
		vals := []interface{}{}
		for i := uint64(0); i < len; i++ {
			val := r.rdbGenericLoadStringObject(RDB_LOAD_SDS)
			var score interface{}
			if vType == RDB_TYPE_ZSET_2 {
				score = r.rdbLoadBinaryDoubleValue()
			} else {
				score = r.rdbLoadDoubleValue()
			}
			vals = append(vals, fmt.Sprintf("%v:%v", val, score))
		}
		return vals
	} else if vType == RDB_TYPE_HASH {
		panic(fmt.Sprintln("can't handle RDB_TYPE_HASH type object."))
	} else if vType == RDB_TYPE_LIST_QUICKLIST {
		len, _, err := r.rdbLoadLen()
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadObject] rdbLoadLen() err:", err))
		}
		vals := []interface{}{}
		for i := uint64(0); i < len; i++ {
			vals = append(vals, r.rdbGenericLoadStringObject(RDB_LOAD_PLAIN))
		}
		return vals
	} else if vType == RDB_TYPE_HASH_ZIPMAP ||
		vType == RDB_TYPE_LIST_ZIPLIST ||
		vType == RDB_TYPE_SET_INTSET ||
		vType == RDB_TYPE_ZSET_ZIPLIST ||
		vType == RDB_TYPE_HASH_ZIPLIST {
		val := r.rdbGenericLoadStringObject(RDB_LOAD_PLAIN)
		// handle val ....
		return val
	} else if vType == RDB_TYPE_STREAM_LISTPACKS {
		panic(fmt.Sprintln("can't handle RDB_TYPE_STREAM_LISTPACKS type object."))
	} else if vType == RDB_TYPE_MODULE ||
		vType == RDB_TYPE_MODULE_2 {
		panic(fmt.Sprintln("can't handle RDB_TYPE_MODULE|RDB_TYPE_MODULE_2 type object."))
	} else {
		panic(fmt.Sprintf("[rdbLoadObject] unknown RDB encoding type %08b, %d\n", vType, vType))
	}
}

func (r *rio) rdbLoadDoubleValue() interface{} {
	len := make([]byte, 1)
	err := r.rioRead(len)
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadBinaryDoubleValue] rioRead() err:", err))
	}
	if len[0] == 255 {
		return "R_NegInf"
	} else if len[0] == 254 {
		return "R_PosInf"
	} else if len[0] == 253 {
		return "R_Nan"
	} else {
		buf := make([]byte, int(len[0]))
		err = r.rioRead(buf)
		if err != nil {
			panic(fmt.Sprintln("[rdbLoadBinaryDoubleValue] rioRead() err:", err))
		}
		return string(buf)
	}
}

func (r *rio) rdbLoadBinaryDoubleValue() interface{} {
	buf := make([]byte, 8)
	err := r.rioRead(buf)
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadBinaryDoubleValue] rioRead() err:", err))
	}
	return string(buf)
}

/* This is only used to load old databases stored with the RDB_OPCODE_EXPIRETIME
 * opcode. New versions of Redis store using the RDB_OPCODE_EXPIRETIME_MS
 * opcode. */
func (r *rio) rdbLoadTime() uint32 {
	buf := make([]byte, 4)
	err := r.rioRead(buf)
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadTime] rioRead() err:", err))
	}
	return binary.BigEndian.Uint32(buf)
}

/* This function loads a time from the RDB file. It gets the version of the
 * RDB because, unfortunately, before Redis 5 (RDB version 9), the function
 * failed to convert data to/from little endian, so RDB files with keys having
 * expires could not be shared between big endian and little endian systems
 * (because the expire time will be totally wrong). The fix for this is just
 * to call memrev64ifbe(), however if we fix this for all the RDB versions,
 * this call will introduce an incompatibility for big endian systems:
 * after upgrading to Redis version 5 they will no longer be able to load their
 * own old RDB files. Because of that, we instead fix the function only for new
 * RDB versions, and load older RDB versions as we used to do in the past,
 * allowing big endian systems to load their own old RDB files. */
func (r *rio) rdbLoadMillisecondTime(rdbver int) uint64 {
	buf := make([]byte, 8)
	err := r.rioRead(buf)
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadMillisecondTime] rioRead() err:", err))
	}
	if rdbver > 9 { /* Convert in big endian if the system is BE. */
		return binary.BigEndian.Uint64(buf)
	}
	// little endian order
	return binary.LittleEndian.Uint64(buf)
}

func (r *rio) rdbLoadCRC64Checksum(rdbver int) {
	// Starting with RDB Version 5
	// CRC 64 checksum: 8 byte
	if rdbver <= 5 {
		return
	}

	buf := make([]byte, 8)
	err := r.rioRead(buf)
	if err != nil {
		panic(fmt.Sprintln("[rdbLoadCRC64Checksum] rioRead() err:", err))
	}
	rdbConvertPrint("cyan", buf, fmt.Sprintf("=> [rdbLoadCRC64Checksum] rdbver:%d, CRC 64 checksum: %08b", rdbver, buf))
}
