package io

import "fmt"

func PrintVerbose(verbose bool, result string, err error) {
	if !verbose {
		return
	}
	if len(result) > 0 {
		fmt.Printf("result: %v\n", result)
	}
	if err != nil {
		fmt.Printf("error encountered: %v\n", err)
	}
}
