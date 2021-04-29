package autostart

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

//Run is the main entrypoint. It prunes and prepares the configuration
func Run(ctx context.Context, config Config, task string) error {
	installers, err := loadDefaultInstallers(ctx, config)
	if err != nil {
		return err
	}
	config.Installers = installers
	//maybe make macros etc?
	return RunTask(ctx, config, task)
}

//RunInstall entrypoint for running either a task or installing a package.
func RunInstall(ctx context.Context, config Config, pkgOrTask string) error {
	//load the task
	if _, ok := config.Tasks[pkgOrTask]; ok {
		//it's a task, awesome
		return RunTask(ctx, config, pkgOrTask)
	}
	return InstallPackage(ctx, config, pkgOrTask)
}

//TODO (@morgan): this should spawn the cmd execution in a goroutine,
// and check if context gets cancelled.. if it does, stop the cmd and return
func runCmd(ctx context.Context, cmdLine string) (string, error) {
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

func shouldRun(ctx context.Context, skipIf []string, runIf []string) bool {
	// compare runCmd-if
	err := testIf(ctx, runIf)
	if len(runIf) > 0 && err != nil {
		return false
	}
	// compare skip-if
	err = testIf(ctx, skipIf)
	if len(skipIf) > 0 && err == nil {
		return false
	}
	return true
}

func testIf(ctx context.Context, ifStatements []string) error {
	for _, ifs := range ifStatements {
		_, err := runCmd(ctx, ifs)
		if err != nil {
			return err
		}
	}
	return nil
}
