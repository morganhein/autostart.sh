package io

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
	"golang.org/x/xerrors"
)

type Runner interface {
	Run(ctx context.Context, printOnly bool, cmdLine string) (string, error)
}

var _ Runner = (*shell)(nil)

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
	cmdLine = fmt.Sprintf("/bin/bash -c \"%v\"", cmdLine)
	cmdLine = strings.TrimSpace(cmdLine)
	args, err := shellwords.Parse(cmdLine)
	if err != nil {
		return "", xerrors.Errorf("error parsing shell words: %v", err)
	}
	var cmd *exec.Cmd
	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(cmdLine)
	}
	//runCmd the cmd
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return "", xerrors.Errorf("error running command `%v`: %v", cmd.String(), err)
	}
	return out.String(), nil
}
