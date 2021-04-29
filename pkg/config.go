package autostart

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Packages   map[string]Package   `toml:"pkg"`
	Installers map[string]Installer `toml:"installer"`
	Tasks      map[string]Task      `toml:"task"`
	Macros     map[string]Macro     `toml:"macro"`
	//TODO (@morgan): if no priority is given, then all available installers should be used, with sh/bash being last
	InstallerPriority []string `toml:"-"` //ordered list of requested installers, from cli args
	Sudo              string   `toml:"-"` //a string so we can verify if it's set or not
	TmpDir            string   `toml:"temp_dir"`
}

type envVariables map[string]string

func (e envVariables) add(key, value string) {
	//TODO (@morgan): need to prefix variables so we can avoid collisions and are easily identifiable
	e["ash-"+key] = value
}

func ParsePackageConfig(config string) (*Config, error) {
	k := &Config{}
	_, err := toml.Decode(config, &k)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func LoadPackageConfig(location string) (string, error) {
	f, err := ioutil.ReadFile(location)
	if err != nil {
		return "", err
	}
	return string(f), nil
}

//combineConfigs adds all values from the addition config, and over-writes
//the original where duplicates exist
func combineConfigs(original Config, addition Config) Config {
	for pkgName, pkg := range addition.Packages {
		original.Packages[pkgName] = pkg
	}
	for macroName, macro := range addition.Macros {
		original.Macros[macroName] = macro
	}
	for installerName, installer := range addition.Installers {
		original.Installers[installerName] = installer
	}
	for taskName, task := range addition.Tasks {
		original.Tasks[taskName] = task
	}
	return original
}
