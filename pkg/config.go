package autostart

import (
	"context"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type RunConfig struct {
}

type Config struct {
	RunningConfig
	Packages   map[string]Package   `toml:"pkg"`
	Installers map[string]Installer `toml:"installer"`
	Tasks      map[string]Task      `toml:"task"`
	Macros     map[string]Macro     `toml:"macro"`
}

type RunningConfig struct {
	TmpDir         string `toml:"temp_dir"`
	Sudo           string `toml:"-"` //a string so we can verify if it's set or not
	Task           string
	ConfigLocation string
	Verbose        bool
	Installer      Installer
	DryRun         bool
}

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
		k.Packages = map[string]Package{}
	}
	if k.Tasks == nil {
		k.Tasks = map[string]Task{}
	}
	if k.Macros == nil {
		k.Macros = map[string]Macro{}
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
	if original.Macros == nil {
		original.Macros = map[string]Macro{}
	}
	for macroName, macro := range addition.Macros {
		original.Macros[macroName] = macro
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
