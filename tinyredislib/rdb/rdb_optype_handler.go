package rdb

import (
	"fmt"
	"io"
)

/*------------------------------------------------------------------------------
 * 标准定义如下
 * ---------------------------------------------------------------------------*/

const (
	/* Map object types to RDB object types. Macros starting with OBJ_ are for
	 * memory storage and may change. Instead RDB types must be fixed because
	 * we store them on disk. */
	RDB_TYPE_STRING   = 0
	RDB_TYPE_LIST     = 1
	RDB_TYPE_SET      = 2
	RDB_TYPE_ZSET     = 3
	RDB_TYPE_HASH     = 4
	RDB_TYPE_ZSET_2   = 5 /* ZSET version 2 with doubles stored in binary. */
	RDB_TYPE_MODULE   = 6
	RDB_TYPE_MODULE_2 = 7 /* Module value with annotations for parsing without
	   the generating module being loaded. */
	/* NOTE: WHEN ADDING NEW RDB TYPE, UPDATE rdbIsObjectType() BELOW */

	/* Object types for encoded objects. */
	RDB_TYPE_HASH_ZIPMAP      = 9
	RDB_TYPE_LIST_ZIPLIST     = 10
	RDB_TYPE_SET_INTSET       = 11
	RDB_TYPE_ZSET_ZIPLIST     = 12
	RDB_TYPE_HASH_ZIPLIST     = 13
	RDB_TYPE_LIST_QUICKLIST   = 14
	RDB_TYPE_STREAM_LISTPACKS = 15
	/* NOTE: WHEN ADDING NEW RDB TYPE, UPDATE rdbIsObjectType() BELOW */

	/* Special RDB opcodes (saved/loaded with rdbSaveType/rdbLoadType). */
	RDB_OPCODE_MODULE_AUX    = 247 /* Module auxiliary data. */
	RDB_OPCODE_IDLE          = 248 /* LRU idle time. */
	RDB_OPCODE_FREQ          = 249 /* LFU frequency. */
	RDB_OPCODE_AUX           = 250 /* RDB aux field. */
	RDB_OPCODE_RESIZEDB      = 251 /* Hash table resize hint. */
	RDB_OPCODE_EXPIRETIME_MS = 252 /* Expire time in milliseconds. */
	RDB_OPCODE_EXPIRETIME    = 253 /* Old expire time in seconds. */
	RDB_OPCODE_SELECTDB      = 254 /* DB number of the following keys. */
	RDB_OPCODE_EOF           = 255 /* End of the RDB file. */

)

type opcodeHandler func(*rio) error

var opcodeHandlerMap = map[byte]opcodeHandler{
	RDB_OPCODE_MODULE_AUX:    opcodeHandlerDefault,
	RDB_OPCODE_IDLE:          opcodeHandlerIDLE,
	RDB_OPCODE_FREQ:          opcodeHandlerFreq,
	RDB_OPCODE_AUX:           opcodeHandlerAux,
	RDB_OPCODE_RESIZEDB:      opcodeHandlerResizeDB,
	RDB_OPCODE_EXPIRETIME_MS: opcodeHandlerExpiretimeMS,
	RDB_OPCODE_EXPIRETIME:    opcodeHandlerExpiretime,
	RDB_OPCODE_SELECTDB:      opcodeHandlerSelectDB,
	RDB_OPCODE_EOF:           opcodeHandlerEOF,
}

/*------------------------------------------------------------------------------
 * 每种opcode类型的处理器实现如下
 * ---------------------------------------------------------------------------*/

func opcodeHandlerDefault(rdb *rio) error {
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerDefault] unsupported opcode!"))
	return fmt.Errorf("[opcodeHandlerDefault] no opcode handler")
}

func opcodeHandlerExpiretime(rdb *rio) error {
	/* EXPIRETIME: load an expire associated with the next key
	 * to load. Note that after loading an expire we need to
	 * load the actual type, and continue. */
	expire := rdb.rdbLoadTime() * 1000
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerExpiretime] expiretime: %v", expire))
	return nil
}

func opcodeHandlerExpiretimeMS(rdb *rio) error {
	/* EXPIRETIME_MS: milliseconds precision expire times introduced
	 * with RDB v3. Like EXPIRETIME but no with more precision. */
	expire := rdb.rdbLoadMillisecondTime(rdb.rdbver)
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerExpiretimeMS] expiretime(milliseconds): %v", expire))
	return nil
}

func opcodeHandlerFreq(rdb *rio) error {
	/* FREQ: LFU frequency. */
	buf := make([]byte, 1)
	err := rdb.rioRead(buf)
	if err != nil {
		return err
	}
	lfuFreq := int(buf[0])
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerFreq] LFU frequency: %v", lfuFreq))
	return nil
}

func opcodeHandlerIDLE(rdb *rio) error {
	/* IDLE: LRU idle time. */
	lruIDLE, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerIDLE] LRU idle time: %v", lruIDLE))
	return nil
}

func opcodeHandlerEOF(rdb *rio) error {
	/* EOF: End of file, exit the main loop. */
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerEOF] end of data."))
	return io.EOF
}

func opcodeHandlerSelectDB(rdb *rio) error {
	/* SELECTDB: Select the specified database. */
	dbid, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerSelectDB] select db: %v", dbid))
	return nil
}

func opcodeHandlerResizeDB(rdb *rio) error {
	/* RESIZEDB: Hint about the size of the keys in the currently
	 * selected data base, in order to avoid useless rehashing. */
	dbSize, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	expiresSize, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerResizeDB] db_size:%v, expire_size: %v", dbSize, expiresSize))
	return nil
}

func opcodeHandlerAux(rdb *rio) error {
	/* AUX: generic string-string fields. Use to add state to RDB
	 * which is backward compatible. Implementations of RDB loading
	 * are required to skip AUX fields they don't understand.
	 *
	 * An AUX field is composed of two strings: key and value. */
	key := rdb.rdbLoadStringObject() // read key
	val := rdb.rdbLoadStringObject() // read val
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerAux] key: %v, val: %v", key, val))
	return nil
}

func opcodeHandlerModuleAux(rdb *rio) error {
	/* Load module data that is not related to the Redis key space.
	 * Such data can be potentially be stored both before and after the
	 * RDB keys-values section. */
	moduleid, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	whenOpcode, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	when, _, err := rdb.rdbLoadLen()
	if err != nil {
		return err
	}
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[opcodeHandlerModuleAux] moduleid: %v, when_opcode: %v, when: %v", moduleid, whenOpcode, when))
	return fmt.Errorf("[opcodeHandlerModuleAux] need more check ... ")
}

/*------------------------------------------------------------------------------
 * RDB_TYPE_... 类型数据处理器实现如下
 * ---------------------------------------------------------------------------*/

func rdbTypeCommonHandler(rdb *rio, rdbType int) error {
	key := rdb.rdbGenericLoadStringObject(RDB_LOAD_SDS) // read key
	val := rdb.rdbLoadObject(int(rdbType))              // read obj
	rdbConvertPrint("green", []byte{}, fmt.Sprintf("[rdbTypeCommonHandler] key: %v, valType: %v, val: %v", key, rdbType, val))
	return nil
}
