package oops

import (
	"fmt"
	"log"
	"runtime"
)

func Log(err error) error {
	return errorF(err)
}

func ErrorF(any interface{}, a ...interface{}) error {
	return errorF(any, a...)
}

func errorF(any interface{}, a ...interface{}) error {
	if any == nil {
		return nil
	}
	var err error

	switch any.(type) {
	case string:
		err = fmt.Errorf(any.(string), a...)
	case error:
		err = fmt.Errorf(any.(error).Error(), a...)
	default:
		err = fmt.Errorf("%v", err)
	}

	_, fn, line, _ := runtime.Caller(2)

	log.Printf("ERROR: [%s:%d] %v \n", fn, line, err)
	return err
}
