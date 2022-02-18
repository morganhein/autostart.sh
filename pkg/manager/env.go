package manager

import (
	"context"
	"errors"
	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//Loading handled by custom loader
type TOMLConfig struct {
	General    General              `toml:"general"`
	Packages   map[string]Package   `toml:"pkg"`
	Installers map[string]Installer `toml:"installer"`
	Tasks      map[string]Task      `toml:"task"`
}

type Operation string

const (
	SYNC    Operation = "sync"
	INSTALL Operation = "install"
	TASK    Operation = "task"
)

type RunConfig struct {
	Operation      Operation
	Sudo           string
	ConfigLocation string
	//TargetDir is the base directory for symlinks, defaults to ${HOME}
	TargetDir string
	//SourceDir is the base directory to search for source files to symlink against, defaults to dir(ConfigLocation)
	SourceDir string
	Verbose   bool
	DryRun    bool
	//ForceInstaller will force the specified installer without detection
	ForceInstaller string
	TOMLConfig     TOMLConfig
	originalTask   string
}

const (
	ORIGINAL_TASK = "ORIGINAL_TASK"
	CURRENT_TASK  = "CURRENT_TASK"
	CURRENT_PKG   = "CURRENT_PKG"
	SUDO          = "SUDO"
	CONFIG_PATH   = "CONFIG_PATH"
	TARGET_PATH   = "TARGET_PATH"
	SOURCE_PATH   = "SOURCE_PATH"
)

func LoadFileConfig(runConfig RunConfig) (*TOMLConfig, error) {
	k, err := LoadPackageConfig(runConfig)
	if err != nil {
		return nil, err
	}
	if k.Packages == nil {
		k.Packages = map[string]Package{}
	}
	if k.Installers == nil {
		k.Installers = map[string]Installer{}
	}
	if k.Tasks == nil {
		k.Tasks = map[string]Task{}
	}
	return k, nil
}

func LoadPackageConfig(runConfig RunConfig) (*TOMLConfig, error) {
	c, err := loadPackageConfigHelper(runConfig.ConfigLocation)
	if err == nil {
		return c, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	locations := []string{
		"$HOME/.config/shoelace/config.toml", "$HOME/.config/shoelace/default.toml",
		"$HOME/.shoelace/config.toml", "$HOME/.shoelace/default.toml",
		"/usr/share/shoelace/default.toml",
	}
	for _, loc := range locations {
		c, err := loadPackageConfigHelper(strings.Replace(loc, "$HOME", home, 1))
		if err == nil {
			return c, err
		}
	}
	return nil, errors.New("could not find a config file to load")
}

func loadPackageConfigHelper(location string) (*TOMLConfig, error) {
	if location == "" {
		return nil, errors.New("config location is empty")
	}
	f, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	k := &TOMLConfig{}
	_, err = toml.Decode(string(f), k)
	if err != nil {
		return nil, err
	}
	return k, nil
}

//combineConfigs adds all values from the addition config, but keeps originals where duplicates exist
func combineConfigs(original TOMLConfig, addition TOMLConfig) TOMLConfig {
	if original.Packages == nil {
		original.Packages = map[string]Package{}
	}
	for pkgName, pkg := range addition.Packages {
		if _, alreadyExists := original.Packages[pkgName]; !alreadyExists {
			original.Packages[pkgName] = pkg
		}
	}
	if original.Installers == nil {
		original.Installers = map[string]Installer{}
	}
	for installerName, installer := range addition.Installers {
		if _, alreadyExists := original.Installers[installerName]; !alreadyExists {
			original.Installers[installerName] = installer
		}
	}
	if original.Tasks == nil {
		original.Tasks = map[string]Task{}
	}
	for taskName, task := range addition.Tasks {
		if _, alreadyExists := original.Tasks[taskName]; !alreadyExists {
			original.Tasks[taskName] = task
		}
	}
	return original
}

//overwriteConfigs adds all values from the addition config, and over-writes
//the original where duplicates exist
func overwriteConfigs(original TOMLConfig, addition TOMLConfig) TOMLConfig {
	if original.Packages == nil {
		original.Packages = map[string]Package{}
	}
	for pkgName, pkg := range addition.Packages {
		original.Packages[pkgName] = pkg
	}
	if original.Installers == nil {
		original.Installers = map[string]Installer{}
	}
	for installerName, installer := range addition.Installers {
		original.Installers[installerName] = installer
	}
	if original.Tasks == nil {
		original.Tasks = map[string]Task{}
	}
	for taskName, task := range addition.Tasks {
		original.Tasks[taskName] = task
	}
	return original
}

func determineBestAvailableInstaller(ctx context.Context, config RunConfig, pkg Package, d Decider) (*Installer, error) {
	//if execution arguments have forced a specific installer to be used
	if config.ForceInstaller != "" {
		i, ok := config.TOMLConfig.Installers[config.ForceInstaller]
		if ok {
			i.Name = config.ForceInstaller
			return &i, nil
		}
		return nil, xerrors.Errorf("an installer was requested (%v), but was not found", config.ForceInstaller)
	}
	availableInstallers := make([]Installer, 0)
	for installerName, installer := range config.TOMLConfig.Installers {
		sr := d.ShouldRun(ctx, []string{}, installer.RunIf)
		if !sr {
			continue
		}
		installer.Name = installerName
		availableInstallers = append(availableInstallers, installer)
	}
	if requiredInstaller, ok := pkg["prefer"]; ok {
		i, ok := config.TOMLConfig.Installers[requiredInstaller]
		if ok {
			i.Name = requiredInstaller
			return &i, nil
		}
		return nil, xerrors.Errorf("an installer was requested (%v), but was not found", requiredInstaller)
	}
	if len(config.TOMLConfig.General.Installers) > 0 {
		for _, v := range config.TOMLConfig.General.Installers {
			for _, availableInstaller := range availableInstallers {
				if v == availableInstaller.Name {
					return &availableInstaller, nil
				}
			}
		}
		return nil, xerrors.Errorf("preferred installer(s) are not available (%+v)", config.TOMLConfig.General.Installers)
	}

	//no installer preferred, grab the first available one
	for _, installer := range availableInstallers {
		return &installer, nil
	}

	return nil, xerrors.New("unable to find a suitable installer")
}

// Finds the package <name> in the config if found, otherwise returns package with default settings matching <name>
func getPackage(config TOMLConfig, name string) Package {
	for pkgName, pkg := range config.Packages {
		if name == pkgName {
			return pkg
		}
	}
	return Package{}
}

// Environment Variables

type envVariables map[string]string

func (e envVariables) copy() envVariables {
	//TODO (@morgan): I think this can be copied more efficiently
	newEnv := envVariables{}
	for k, v := range e {
		newEnv[k] = v
	}
	return newEnv
}

//set default environment variables
func hydrateEnvironment(config RunConfig, env envVariables) {
	env[ORIGINAL_TASK] = config.originalTask
	env[CONFIG_PATH] = path.Dir(config.ConfigLocation)
	//possibly add link src and dst links here
}
