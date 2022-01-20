package io

import (
	"fmt"
	"runtime"
)

func PrintVerbose(verbose bool, result string, err error) {
	if !verbose {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	if len(result) > 0 {
		fmt.Printf("called from: %v:%v, result: %v\n", file, line, result)
	}
	if err != nil {
		fmt.Printf("called from: %v:%v, error encountered: %v\n", file, line, err)
	}
}
