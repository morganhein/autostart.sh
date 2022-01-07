package manager

import (
	"context"
	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
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
	OriginalTask   string
	ConfigLocation string
	TargetDir      string // The base directory for symlinks, defaults to ${HOME}
	SourceDir      string // The base directory to search for source files to symlink against, defaults to dir(ConfigLocation)
	Verbose        bool
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
    detect = ["which apt"]
    sudo = true
    cmd =  "${sudo} apt install -y ${pkg}"
	update = "${sudo} apt update"

[installer.brew]
    detect = ["which brew"]
    sudo = false
    cmd =  "${sudo} brew install ${pkg}"
	update = "${sudo} brew update"

[installer.apk]
    detect = ["which apk"]
    sudo = false
    cmd =  "${sudo} apk add ${pkg}"
	update = "${sudo} apk update"

[installer.dnf]
    detect = ["which dnf"]
    sudo = true
    cmd =  "${sudo} dnf install -y ${pkg}"

[installer.pacman]
    detect = ["which pacman"]
    sudo = true
    cmd =  "${sudo} pacman -S ${pkg}"

[installer.yay]
    detect = ["which yay"]
    sudo = true
    cmd =  "${sudo} yay -S ${pkg}"
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

//combineConfigs adds all values from the addition config, but keeps originals where duplicates exist
func combineConfigs(original Config, addition Config) Config {
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
func overwriteConfigs(original Config, addition Config) Config {
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

//When pre-existing values exist, the values should not be over-written.
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
	return loadDefaultInstallers(config)
}

//insure we have default installers.
//If an installer already exists where a default would be loaded, the original is kept
func loadDefaultInstallers(config Config) (Config, error) {
	defaultConfig := &Config{}
	err := toml.Unmarshal([]byte(installerDefaults), defaultConfig)
	if err != nil {
		return config, xerrors.Errorf("error unmarshalling config: %v", err)
	}
	return combineConfigs(config, *defaultConfig), nil
}

func determineBestAvailableInstaller(ctx context.Context, config Config, pkg Package, d Decider) (*Installer, error) {
	//if execution arguments have forced a specific installer to be used
	if config.ForceInstaller != "" {
		i, ok := config.Installers[config.ForceInstaller]
		if ok {
			i.Name = config.ForceInstaller
			return &i, nil
		}
		return nil, xerrors.Errorf("an installer was requested (%v), but was not found", config.ForceInstaller)
	}
	availableInstallers := make([]Installer, 0)
	for installerName, installer := range config.Installers {
		sr := d.ShouldRun(ctx, []string{}, installer.RunIf)
		if !sr {
			continue
		}
		installer.Name = installerName
		availableInstallers = append(availableInstallers, installer)
	}
	//if the package defined a required installer, check if it is available
	//TODO: This should handle a list of installer preferences, comma-separated
	if requiredInstaller, ok := pkg["prefer"]; ok {
		i, ok := config.Installers[requiredInstaller]
		if ok {
			i.Name = requiredInstaller
			return &i, nil
		}
		return nil, xerrors.Errorf("an installer was requested (%v), but was not found", requiredInstaller)
	}

	//no installer preferred, grab the first available one
	for _, installer := range availableInstallers {
		return &installer, nil
	}

	return nil, xerrors.New("unable to find a suitable installer")
}

// Finds the package <name> in the config if found, otherwise returns package with default settings matching <name>
func getPackage(config Config, name string) Package {
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
func hydrateEnvironment(config Config, env envVariables) {
	env[ORIGINAL_TASK] = config.OriginalTask
	env[CONFIG_PATH] = path.Dir(config.ConfigLocation)
	//possibly add link src and dst links here
}
