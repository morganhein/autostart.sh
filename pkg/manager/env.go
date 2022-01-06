package manager

import (
	"context"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	RunningConfig
	Packages   map[string]Package   `toml:"pkg"`
	Installers map[string]Installer `toml:"installer"`
	Tasks      map[string]Task      `toml:"task"`
}

//What is unique about this config vs the other ones??
type RunningConfig struct {
	TmpDir         string `toml:"temp_dir"`
	Sudo           string `toml:"-"` //a string so we can verify if it's set or not
	Task           string
	ConfigLocation string
	TargetDir      string // The base directory for symlinks, defaults to ${HOME}
	SourceDir      string // The base directory to search for source files to symlink against, defaults to dir(ConfigLocation)
	Verbose        bool
	Installer      Installer
	DryRun         bool
	ForceInstaller string // will force the specified installer without detection
}

const (
	ORIGINAL_TASK = "ORIGINAL_TASK"
	CURRENT_TASK  = "CURRENT_TASK"
	CURRENT_PKG   = "CURRENT_PKG"
	SUDO          = "SUDO"
	CONFIG_PATH   = "CONFIG_PATH"
	TARGET_PATH   = "TARGET_PATH"
	SOURCE_PATH   = "SOURCE_PATH"

	installerDefaults = `
[installer.apt]
    run_if = ["which apt", "which apt-get"]
    sudo = true
    cmd =  "${sudo} apt install -y ${pkg}"
	update = "${sudo} apt update"

[installer.brew]
    run_if = ["which brew"]
    sudo = false
    cmd =  "${sudo} brew install ${pkg}"
	update = "${sudo} brew update"

[installer.apk]
    run_if = ["which apk"]
    sudo = false
    cmd =  "${sudo} apk add ${pkg}"
	update = "${sudo} apk update"

[installer.dnf]
    run_if = ["which dnf"]
    sudo = true
    cmd =  "${sudo} dnf install -y ${pkg}"

[installer.pacman]
    run_if = ["which pacman"]
    sudo = true
    cmd =  "${sudo} pacman -Syu ${pkg}"

[installer.yay]
    run_if = ["which yay"]
    sudo = true
    cmd =  "${sudo} yay -Syu ${pkg}"
`
)

// Config Handling

func ParsePackageConfig(config string) (*Config, error) {
	k := &Config{}
	_, err := toml.Decode(config, &k)
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

func LoadPackageConfig(ctx context.Context, location string) (*Config, error) {
	f, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	return ParsePackageConfig(string(f))
}

//combineConfigs adds all values from the addition config, and over-writes
//the original where duplicates exist
func combineConfigs(original Config, addition Config) Config {
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

func insureDefaults(config Config) (Config, error) {
	if config.SourceDir == "" {
		if config.ConfigLocation == "" {
			return config, xerrors.New("cannot determine source directory, since SourceDir and ConfigLocation are unset")
		}
		config.SourceDir = path.Dir(config.ConfigLocation)
	}
	if config.TargetDir == "" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return config, xerrors.Errorf("error retrieving home directory: %v", err)
		}
		config.TargetDir = dirname
	}
	return config, nil
}

func loadDefaultInstallers(config Config) (Config, error) {
	defaultConfig := &Config{}
	err := toml.Unmarshal([]byte(installerDefaults), defaultConfig)
	if err != nil {
		return config, xerrors.Errorf("error unmarshalling config: %v", err)
	}
	return combineConfigs(config, *defaultConfig), nil
}

func detectInstaller(ctx context.Context, config Config, d Decider) (*Installer, error) {
	if config.ForceInstaller != "" {
		i, ok := config.Installers[config.ForceInstaller]
		if ok {
			i.Name = config.ForceInstaller
			return &i, nil
		}
		return nil, xerrors.Errorf("an installer was requested (%v), but was not found", config.ForceInstaller)
	}
	for k, v := range config.Installers {
		sr := d.ShouldRun(ctx, v.SkipIf, v.RunIf)
		if !sr {
			continue
		}
		v.Name = k
		return &v, nil
	}
	return nil, xerrors.New("unable to find a suitable installer")
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
func hydrateEnvironment(config Config, env envVariables) {
	env[ORIGINAL_TASK] = config.Task
	env[CONFIG_PATH] = path.Dir(config.ConfigLocation)
	//possibly add link src and dst links here
}
