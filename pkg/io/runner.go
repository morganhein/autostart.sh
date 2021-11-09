package io

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(ctx context.Context, printOnly bool, cmdLine string) (string, error)
}

func NewShellRunner() *shell {
	return &shell{}
}

type shell struct{}

//TODO (@morgan): this should spawn the cmd execution in a goroutine,
// and check if context gets cancelled.. if it does, stop the cmd and return
func (s shell) Run(ctx context.Context, printOnly bool, cmdLine string) (string, error) {
	if printOnly {
		fmt.Println(cmdLine)
		return "", nil
	}
	args := strings.Split(cmdLine, " ")
	var cmd *exec.Cmd
	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(cmdLine)
	}
	//runCmd the cmd
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
