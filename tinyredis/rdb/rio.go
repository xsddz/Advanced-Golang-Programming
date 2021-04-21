package rdb

import "os"

type rio struct {
	fp *os.File
}

func rioRead(r *rio, buf []byte) error {
	_, e := r.fp.Read(buf)
	if e != nil {
		return e
	}
	return nil
}
