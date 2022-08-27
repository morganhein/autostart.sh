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

func PrintVerboseF(verbose bool, format string, args ...interface{}) {
	if !verbose {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	if len(format) > 0 {
		fmt.Printf("called from: %v:%v, result: %v\n", file, line, fmt.Sprintf(format, args...))
	}
}

func PrintVerboseError(verbose bool, err error) {
	if !verbose {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		fmt.Printf("called from: %v:%v, error encountered: %v\n", file, line, err)
	}
}
