package autostart

import (
	"context"

	"github.com/morganhein/autostart.sh/pkg/io"
)

type Decider interface {
	ShouldRun(ctx context.Context, skipIf []string, runIf []string) bool
}

func NewDecider(r io.Runner) *decider {
	return &decider{r: r}
}

type decider struct {
	r io.Runner
}

func (d decider) ShouldRun(ctx context.Context, skipIf []string, runIf []string) bool {
	// compare runCmd-if
	err := d.testIf(ctx, runIf)
	if len(runIf) > 0 && err != nil {
		return false
	}
	// compare skip-if
	err = d.testIf(ctx, skipIf)
	if len(skipIf) > 0 && err == nil {
		return false
	}
	return true
}

func (d decider) testIf(ctx context.Context, ifStatements []string) error {
	for _, ifs := range ifStatements {
		//detection can never be a "dry run"
		_, err := d.r.Run(ctx, false, ifs)
		if err != nil {
			return err
		}
	}
	return nil
}
