package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"github.com/axgle/mahonia"
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


func ConvertToString(src string, srcCode string, tagCode string) string {
    srcCoder := mahonia.NewDecoder(srcCode)
    srcResult := srcCoder.ConvertString(src)
    tagCoder := mahonia.NewDecoder(tagCode)
    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
    result := string(cdata)
    return result
}
