package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"img-crawler/src/log"
	"io/ioutil"
)

func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		log.Warn("[ParseGzip] NewReader error: %v, maybe data is ungzip", err)
		return nil, err
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			log.Warn("[ParseGzip]  ioutil.ReadAll error: %v", err)
			return nil, err
		}
		return undatas, nil
	}
}
