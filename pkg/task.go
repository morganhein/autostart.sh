package autostart

import (
	"context"
	"fmt"
)

type Task struct {
	RunIf  []string
	SkipIf []string
	Deps   []string
	Cmds   []string
	Links  []string
}

//RunTask runs a specific task. It does not try to install the task as regular
//package if the task is not found.
func RunTask(ctx context.Context, config Config, task string) error {
	//load the task
	t, ok := config.Tasks[task]
	if !ok {
		return fmt.Errorf("task '%v' not defined in config", task)
	}
	if sr := shouldRun(ctx, t.SkipIf, t.RunIf); !sr {
		return fmt.Errorf("task '%v' failed skip_if or run_if check", task)
	}

	//run the deps
	for _, dep := range t.Deps {
		// try to run that task
		if err := RunInstall(ctx, config, dep); err != nil {
			return err
		}
	}

	//run the cmds

	//run the links
	return nil
}
