package manager

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/xerrors"
	"path"
	"strings"

	"github.com/morganhein/envy/pkg/io"
)

type Manager interface {
	// RunTask will explicitly only run a specified task, and will fail if it is not found
	RunTask(ctx context.Context, config TOMLConfig, task string) error
	//RunInstall will explicitly only run the installation of the package.
	RunInstall(ctx context.Context, config TOMLConfig, pkg string) error
}

type Task struct {
	Installers []string
	RunIf      []string
	SkipIf     []string
	Download   []Downloads
	Deps       []string
	PreCmds    []string `toml:"pre_cmd"`
	Install    []string
	PostCmds   []string `toml:"post_cmd"`
}

type Downloads []string

type Installer struct {
	Name    string
	RunIf   []string `toml:"run_if"`
	Sudo    bool
	Cmd     string
	Update  string
	Updated bool
}

type General struct {
	Installers []string
	ConfigDir  string `toml:"config_dir"`
	HomeDir    string `toml:"home_dir"`
}

type Package map[string]string

type PkgInstallOption struct {
	Name string
}

type manager struct {
	d  Decider
	r  io.Runner
	dl io.Downloader
	s  io.Filesystem
}

func New() manager {
	shell := io.NewShellRunner()
	d := NewDecider(shell)
	return manager{
		d: d,
		r: shell,
	}
}

// Start is the command line entrypoint
func Start(ctx context.Context, config RunConfig, task string) error {
	shell := io.NewShellRunner()
	d := NewDecider(shell)
	m := manager{
		d: d,
		r: shell,
	}
	config.originalTask = task
	io.PrintVerbose(config.Verbose, fmt.Sprintf("startup config: %+v", config), nil)
	return m.RunTask(ctx, config, task)
}

func (m manager) RunTask(ctx context.Context, config RunConfig, task string) error {
	//start tracking environment variables
	vars := envVariables{}
	hydrateEnvironment(config, vars)
	return m.runTaskHelper(ctx, config, vars, task)
}

func (m manager) RunInstall(ctx context.Context, config RunConfig, pkg string) error {
	//start tracking environment variables
	vars := envVariables{}
	hydrateEnvironment(config, vars)
	//this should go straight to the pkg install helper, and none of this other business
	return m.installPkgHelper(ctx, config, vars, pkg)
}

func (m manager) handleDependency(ctx context.Context, config RunConfig, vars envVariables, taskOrPkg string) error {
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

/*runTaskHelper runs, in order:
* Determines if the installers required by the task are available
* If `run_if` passes
* If `skip_if` passes
* Downloads any necessary files
* Installs any deps
* Runs the pre_cmd commands
* Installs the package
* Runs the pose_cmd commands
 */
func (m manager) runTaskHelper(ctx context.Context, config RunConfig, vars envVariables, task string) error {
	io.PrintVerbose(config.Verbose, fmt.Sprintf("starting task [%v]", task), nil)
	//load the task
	t, ok := config.TOMLConfig.Tasks[task]
	if !ok {
		return xerrors.Errorf("task '%v' not defined in config", task)
	}
	//insure the requested installer, if provided, is available
	installerAvailable := func() bool {
		if len(t.Installers) == 0 {
			return true
		}
		for _, installer := range t.Installers {
			if _, ok := config.TOMLConfig.Installers[installer]; ok {
				return true
			}
		}
		return false
	}()
	if !installerAvailable {
		return xerrors.New("none of the installers requested by the package are available")
	}

	if sr := m.d.ShouldRun(ctx, t.SkipIf, t.RunIf); !sr {
		io.PrintVerbose(config.Verbose, fmt.Sprintf("task '%v' failed skip_if or run_if check", task), nil)
		return nil
	}

	//download the files
	for _, dlReq := range t.Download {
		if len(dlReq) != 2 {
			return xerrors.New("the download command must contain two parameters, the source and the target")
		}
		_, err := m.dl.Download(ctx, dlReq[0], dlReq[1])
		if err != nil {
			return err
		}
	}

	//run the deps
	for _, dep := range t.Deps {
		if err := m.handleDependency(ctx, config, vars, dep); err != nil {
			return err
		}
	}

	//run the pre-cmds
	for _, cmd := range t.PreCmds {
		if err := m.runCmdHelper(ctx, config, vars, cmd); err != nil {
			return err
		}
	}

	//install the packages
	for _, pkg := range t.Install {
		if err := m.RunInstall(ctx, config, pkg); err != nil {
			return err
		}
	}

	//run the post-cmds
	for _, cmd := range t.PostCmds {
		if err := m.runCmdHelper(ctx, config, vars, cmd); err != nil {
			return err
		}
	}

	return nil
}

//runCmdHelper resolves any package names and installation commands to the current targets variant, and then runs it
func (m manager) runCmdHelper(ctx context.Context, config RunConfig, vars envVariables, cmdLine string) error {
	//cleanup first
	cmdLine = strings.TrimSpace(cmdLine)
	sudo := determineSudo(config, nil)
	cmdLine = injectVars(cmdLine, vars, sudo)
	io.PrintVerbose(config.Verbose, fmt.Sprintf("running command `%v`", cmdLine), nil)
	out, err := m.r.Run(ctx, config.DryRun, cmdLine)
	io.PrintVerbose(config.Verbose, out, err)
	if err != nil {
		return err
	}
	return nil
}

func (m manager) downloadHelper(ctx context.Context, dl Downloads) (string, error) {
	if len(dl) == 2 {
		return m.dl.Download(ctx, dl[0], dl[1])
	}
	return "", errors.New("incorrect syntax for a download command")
}

//TODO (@morgan): this should probably be removed? in lieu of the sync operation?
func (m manager) symlinkHelper(ctx context.Context, config RunConfig, vars envVariables, link string) error {
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

func (m manager) installPkgHelper(ctx context.Context, config RunConfig, vars envVariables, pkgName string) error {
	if len(pkgName) == 0 {
		return errors.New("unable to find the package name")
	}

	//look up the package in the config, if it exists.
	pkg := getPackage(config.TOMLConfig, pkgName)

	//determine which installer is preferred with this package
	installer, err := determineBestAvailableInstaller(ctx, config, pkg, m.d)
	if err != nil {
		return err
	}

	//determine package name in relation to the chosen installer
	newPkgName, ok := pkg[installer.Name]
	if !ok {
		newPkgName = pkgName
	}

	//run the install commands for that installer
	//do we sudo, or do we not?
	sudo := determineSudo(config, installer)
	cmdLine := installCommandVariableSubstitution(installer.Cmd, newPkgName, sudo)

	//TODO: capture output here for verbose logging
	_, err = m.r.Run(ctx, config.DryRun, cmdLine)
	if err != nil {
		return err
	}
	return nil
}
