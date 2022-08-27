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

type Shell interface {
	Run(ctx context.Context, printOnly bool, cmdLine string) (string, error)
	// Runs `which` on the target shell, to determine if a program exists or not.
	Which(ctx context.Context, search string) (bool, string, error)
}

var _ Shell = (*shell)(nil)

func NewShell() *shell {
	return &shell{}
}

type shell struct{}

func (s shell) Which(ctx context.Context, search string) (bool, string, error) {
	//TODO (@morgan): this exact location is not good! needs to handle other locations and OS
	cmdLine := fmt.Sprintf("/bin/bash -c \"which %v\"", search)
	cmdLine = strings.TrimSpace(cmdLine)
	args, err := shellwords.Parse(cmdLine)
	if err != nil {
		return false, "", xerrors.Errorf("error parsing shell words: %v", err)
	}
	var cmd *exec.Cmd
	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(cmdLine)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return false, out.String(), xerrors.Errorf("error running `which` for `%v` due to: %v", search, err)
	}
	return true, out.String(), nil
}

// TODO (@morgan): this should spawn the cmd execution in a goroutine,
// and check if context gets cancelled.. if it does, stop the cmd and return
// it returns the stdio/stderr combined as a single string
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
		return "", xerrors.Errorf("error running command `%v` due to: %v", cmd.String(), err)
	}
	return out.String(), nil
}
