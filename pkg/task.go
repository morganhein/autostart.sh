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

//runTask runs a specific task. It does not try to install the task as regular
//package if the task is not found.
func runTask(ctx context.Context, config Config, vars envVariables, task string) error {
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
		if err := startPkgOrTask(ctx, config, vars, dep); err != nil {
			return err
		}
	}
	//copy env vars, b/c from here on out it's destructive
	//run the cmds
	for _, cmd := range t.Cmds {
		err := installPackage(ctx, config, cmd)
		if err != nil {
			return err
		}
	}
	//TODO (@morgan): run the links/linking
	return nil
}
