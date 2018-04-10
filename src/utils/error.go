package utils

import (
	"img-crawler/src/log"
	"runtime"
)

func CheckError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Errorf("[%s:%d] {%s}", file, line, err)
		panic(err)
	}
}
