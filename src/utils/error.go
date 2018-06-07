package utils

import (
	"fmt"
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

type CustomError struct {
	Code int
	Msg  string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Msg)
}
