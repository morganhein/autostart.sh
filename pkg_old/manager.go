package autostart

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/morganhein/autostart.sh/pkg_old/io"
	"golang.org/x/xerrors"
)

type Manager interface {
	// RunTask will explicitly only run a specified task, and will fail if it is not found
	RunTask(ctx context.Context, config Config, task string) error
	//RunInstall will explicitly only run the installation of the package.
	RunInstall(ctx context.Context, config Config, pkg string) error
	//Symlink creates the necessary symlinks as requested
	Symlink(ctx context.Context, config Config) error
}

type Task struct {
	RunIf  []string
	SkipIf []string
	Deps   []string
	Cmds   []string
	Links  []string
}

type Installer struct {
	Name    string
	SkipIf  []string `toml:"skip_if"`
	RunIf   []string `toml:"run_if"`
	Sudo    bool
	Cmd     string
	Update  string
	Updated bool
}

type Package map[string]string

type PkgInstallOption struct {
	Name string
}

type manager struct {
	d  Decider
	r  io.Runner
	dl io.Downloader
	s  io.Symlinker
}

// Start is the command line entrypoint
func Start(ctx context.Context, config Config, task string) error {
	shell := io.NewShellRunner()
	d := NewDecider(shell)
	m := manager{
		d: d,
		r: shell,
	}
	io.PrintVerbose(config.Verbose, fmt.Sprintf("startup config: %+v", config), nil)
	return m.RunTask(ctx, config, task)
}

func (m manager) RunTask(ctx context.Context, config Config, task string) error {
	config, err := insureDefaults(config)
	if err != nil {
		return err
	}
	config, err = loadDefaultInstallers(config)
	if err != nil {
		return err
	}
	installer, err := detectInstaller(ctx, config, m.d)
	if err != nil {
		return err
	}
	config.Installer = *installer
	//maybe make macros etc?

	//start tracking environment variables
	vars := envVariables{}
	hydrateEnvironment(config, vars)
	return m.runTaskHelper(ctx, config, vars, task)
}

func (m manager) RunInstall(ctx context.Context, config Config, pkg string) error {
	config, err := insureDefaults(config)
	if err != nil {
		return err
	}
	config, err = loadDefaultInstallers(config)
	if err != nil {
		return err
	}
	installer, err := detectInstaller(ctx, config, m.d)
	if err != nil {
		return err
	}
	config.Installer = *installer
	//start tracking environment variables
	vars := envVariables{}
	hydrateEnvironment(config, vars)
	cmdLine := fmt.Sprintf("@install %v", pkg)
	return m.runCmdHelper(ctx, config, vars, cmdLine)
}

func (m manager) handleDependency(ctx context.Context, config Config, vars envVariables, taskOrPkg string) error {
	if len(taskOrPkg) == 0 {
		return xerrors.New("task or package is empty")
	}
	switch taskOrPkg[0] {
	case '^':
		cmdLine := fmt.Sprintf("@install %v", taskOrPkg)
		return m.runCmdHelper(ctx, config, vars, cmdLine)
	case '#':
		return m.runTaskHelper(ctx, config, vars, taskOrPkg[1:])
	}
	//default is just a plain package name
	cmdLine := fmt.Sprintf("@install %v", taskOrPkg)
	return m.runCmdHelper(ctx, config, vars, cmdLine)
}

func (m manager) runTaskHelper(ctx context.Context, config Config, vars envVariables, task string) error {
	io.PrintVerbose(config.Verbose, fmt.Sprintf("starting task [%v]", task), nil)
	//load the task
	t, ok := config.Tasks[task]
	if !ok {
		return xerrors.Errorf("task '%v' not defined in config", task)
	}
	if sr := m.d.ShouldRun(ctx, t.SkipIf, t.RunIf); !sr {
		io.PrintVerbose(config.Verbose, fmt.Sprintf("task '%v' failed skip_if or run_if check", task), nil)
		return nil
	}

	//run the deps
	for _, dep := range t.Deps {
		if err := m.handleDependency(ctx, config, vars, dep); err != nil {
			return err
		}
	}

	//symlinks first, so that we can create links before installers do
	for _, link := range t.Links {
		if err := m.symlinkHelper(ctx, config, vars, link); err != nil {
			return err
		}
	}

	//copy env vars, b/c from here on out it's destructive
	//run the cmds
	for _, cmd := range t.Cmds {
		if err := m.runCmdHelper(ctx, config, vars, cmd); err != nil {
			return err
		}
	}

	return nil
}

//runCmdHelper resolves any package names and installation commands to the current targets variant, and then runs it
func (m manager) runCmdHelper(ctx context.Context, config Config, vars envVariables, cmdLine string) error {
	//cleanup first
	cmdLine = strings.TrimSpace(cmdLine)
	if strings.HasPrefix(cmdLine, "@download") {
		out, err := m.downloadHelper(ctx, cmdLine)
		io.PrintVerbose(config.Verbose, out, err)
		return err
	}
	if strings.HasPrefix(cmdLine, "@install") {
		out, err := m.installHelper(ctx, config, vars, cmdLine)
		io.PrintVerbose(config.Verbose, out, err)
		return err
	}
	sudo := determineSudo(config, config.Installer)
	cmdLine = injectVars(cmdLine, vars, sudo)
	io.PrintVerbose(config.Verbose, fmt.Sprintf("running command `%v`", cmdLine), nil)
	out, err := m.r.Run(ctx, config.DryRun, cmdLine)
	io.PrintVerbose(config.Verbose, out, err)
	if err != nil {
		return err
	}
	return nil
}

//TODO (@morgan): needs to inject vars as necessary
func (m manager) downloadHelper(ctx context.Context, cmdLine string) (string, error) {
	cmdLine = strings.TrimPrefix(cmdLine, "@download")
	parts := strings.Split(cmdLine, " ")
	if len(parts) == 2 {
		return m.dl.Download(ctx, parts[0], parts[1])
	}
	return "", errors.New("incorrect syntax for a download command, it must be of form `@download http://source.com file://target_location")
}

func (m manager) installHelper(ctx context.Context, config Config, vars envVariables, cmdLine string) (string, error) {
	//get package name
	pkgName := getPackageName(config, strings.TrimPrefix(cmdLine, "@install "))
	io.PrintVerbose(config.Verbose, fmt.Sprintf("installing `%v`", pkgName), nil)
	if len(pkgName) == 0 {
		return "", errors.New("unable to find the package name")
	}
	cmdLine = injectPackage(cmdLine, config.Installer.Cmd, pkgName)
	//do we sudo, or do we not?
	sudo := determineSudo(config, config.Installer)
	if vars == nil {
		vars = envVariables{}
	}
	cmdLine = strings.TrimSpace(injectVars(cmdLine, vars, sudo))
	//if the update command hasn't happened for this installer, and the command is not empty
	var updateResult string
	var err error
	if !config.Installer.Updated && config.Installer.Update != "" {
		updateLine := strings.TrimSpace(injectVars(config.Installer.Update, vars, sudo))
		updateResult, err = m.r.Run(ctx, config.DryRun, updateLine)
		io.PrintVerbose(config.Verbose, updateResult, err)
		if err != nil {
			return "", err
		}

	}
	cmdResult, err := m.r.Run(ctx, config.DryRun, cmdLine)
	if len(updateResult) > 0 {
		cmdResult = updateResult + "\n" + cmdResult
	}
	return cmdResult, err
}

func (m manager) symlinkHelper(ctx context.Context, config Config, vars envVariables, link string) error {
	io.PrintVerbose(config.Verbose, fmt.Sprintf("creating symlink `%v`", link), nil)
	parts := strings.Split(link, " ")
	if len(parts) > 2 {
		return xerrors.New("unexpected symlink format, which is `from [to]`")
	}
	from := path.Join(config.SourceDir, parts[0])
	to := path.Join(config.TargetDir, parts[0])
	if len(parts) == 2 {
		to = path.Join(config.TargetDir, parts[1])
	}

	if config.DryRun {
		fmt.Printf("symlinking from %v to %v\n", from, to)
		return nil
	}

	out, err := func() (string, error) {
		fmt.Println(link)
		return "", nil
	}()
	io.PrintVerbose(config.Verbose, out, err)
	return err
}
