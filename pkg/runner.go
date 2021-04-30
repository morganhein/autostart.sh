package autostart

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

//RunTask is main entrypoint for running and installing a task
func RunTask(ctx context.Context, config Config, task string) error {
	installer, err := loadDefaultInstallers(ctx, config)
	if err != nil {
		return err
	}
	config.Installer = *installer
	//maybe make macros etc?

	//start tracking environment variables
	vars := envVariables{}
	hydrateEnvironment(config, vars)
	return runTask(ctx, config, vars, task)
}

//RunInstall is the main entrypoint for installing a package.
func RunInstall(ctx context.Context, config Config, pkgOrTask string) error {
	installer, err := loadDefaultInstallers(ctx, config)
	if err != nil {
		return err
	}
	config.Installer = *installer
	//maybe make macros etc?

	//start tracking environment variables
	vars := envVariables{}
	hydrateEnvironment(config, vars)
	return installPackage(ctx, config, pkgOrTask)
}

//runInstall tries to first run the task, then install the package if no task is found
func runInstall(ctx context.Context, config Config, vars envVariables, pkgOrTask string) error {
	//load the task
	if _, ok := config.Tasks[pkgOrTask]; ok {
		//it's a task, awesome
		vars[CURRENT_TASK] = pkgOrTask
		return runTask(ctx, config, vars, pkgOrTask)
	}
	vars[CURRENT_PKG] = pkgOrTask
	return installPackage(ctx, config, pkgOrTask)
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
