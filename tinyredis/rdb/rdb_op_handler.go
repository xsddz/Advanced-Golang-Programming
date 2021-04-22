package rdb

import (
	"fmt"
	"io"
)

const (
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
	RDB_TYPE_MODULE_2      = 7 /* Module value with annotations for parsing without the generating module being loaded. */
	/* NOTE: WHEN ADDING NEW RDB TYPE, UPDATE rdbIsObjectType() BELOW */

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

)

type opHandler func(*rio) error

var opTypeHandlerMap = map[byte]opHandler{
	// segment optype
	RDB_OPCODE_EOF:      opcodeHandlerEOF,
	RDB_OPCODE_AUX:      opcodeHandlerAux,
	RDB_OPCODE_SELECTDB: opcodeHandlerSelectDB,
	RDB_OPCODE_RESIZEDB: opcodeHandlerResizeDB,

	// data optype
	RDB_TYPE_STRING: rdbTypeHandlerString,
}

func opcodeHandlerEOF(rdb *rio) error {
	/* EOF: End of file, exit the main loop. */
	rdbConvertPrint("cyan", []byte{}, fmt.Sprintf("[opcodeHandlerEOF] end of data."))
	return io.EOF
}

func opcodeHandlerAux(rdb *rio) error {
	/* AUX: generic string-string fields. Use to add state to RDB
	 * which is backward compatible. Implementations of RDB loading
	 * are requierd to skip AUX fields they don't understand.
	 *
	 * An AUX field is composed of two strings: key and value. */
	key := rdb.rdbLoadStringObject() // read key
	val := rdb.rdbLoadStringObject() // read val
	rdbConvertPrint("cyan", []byte{}, fmt.Sprintf("[opcodeHandlerAux] key: %v, val: %v", key, val))
	return nil
}

func opcodeHandlerSelectDB(rdb *rio) error {
	/* SELECTDB: Select the specified database. */
	dbid, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	rdbConvertPrint("cyan", []byte{}, fmt.Sprintf("[opcodeHandlerSelectDB] select db: %v", dbid))
	return nil
}

func opcodeHandlerResizeDB(rdb *rio) error {
	/* RESIZEDB: Hint about the size of the keys in the currently
	 * selected data base, in order to avoid useless rehashing. */
	dbSize, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	expireSize, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	rdbConvertPrint("cyan", []byte{}, fmt.Sprintf("[opcodeHandlerResizeDB] db_size:%v, expire_size: %v", dbSize, expireSize))
	return nil
}

func rdbTypeHandlerString(rdb *rio) error {
	key := rdb.rdbLoadStringObject()          // read key
	val := rdb.rdbLoadObject(RDB_TYPE_STRING) // read obj
	rdbConvertPrint("cyan", []byte{}, fmt.Sprintf("[rdbTypeHandlerString] key: %v, val: %v", key, val))
	return nil
}
