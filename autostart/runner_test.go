package autostart

import (
  "context"
  "testing"
)

func TestRunCmd(t *testing.T) {
  ctx := context.Background()
  err := Run(ctx, RunArgs{Cmd: "/usr/bin/bash", Sudo: false, Args: []string{"-c", "echo", "hi"}})
  if err != nil {
    t.Fail()
  }
}
