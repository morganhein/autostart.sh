package manager

import (
	"errors"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/morganhein/envy/pkg/io"
)

type TOMLConfig struct {
	General       General              `toml:"general"`
	Packages      map[string]Package   `toml:"pkg"`
	InstallerDefs map[string]Installer `toml:"installer"`
	Tasks         map[string]Task      `toml:"task"`
}

func LoadFileConfig(fs io.Filesystem, configLocation string) (*TOMLConfig, error) {
	k, err := LoadPackageConfig(fs, configLocation)
	if err != nil {
		return nil, err
	}
	if k.Packages == nil {
		k.Packages = map[string]Package{}
	}
	if k.InstallerDefs == nil {
		k.InstallerDefs = map[string]Installer{}
	}
	if k.Tasks == nil {
		k.Tasks = map[string]Task{}
	}
	return k, nil
}

func LoadPackageConfig(fs io.Filesystem, configLocation string) (*TOMLConfig, error) {
	c, err := loadPackageConfigHelper(fs, configLocation)
	if err == nil {
		return c, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	locations := []string{
		"$HOME/.config/envy/config.toml", "$HOME/.config/envy/default.toml",
		"$HOME/.envy/config.toml", "$HOME/.envy/default.toml",
		"/usr/share/envy/default.toml",
	}
	for _, loc := range locations {
		c, err := loadPackageConfigHelper(fs, strings.Replace(loc, "$HOME", home, 1))
		if err == nil {
			return c, err
		}
	}
	return nil, errors.New("could not find a config file to load")
}

func loadPackageConfigHelper(fs io.Filesystem, location string) (*TOMLConfig, error) {
	if location == "" {
		return nil, errors.New("config location is empty")
	}
	f, err := fs.ReadFile(location)
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

// Finds the package <name> in the config if found, otherwise returns package with default settings matching <name>
func getPackage(config TOMLConfig, name string) Package {
	for pkgName, pkg := range config.Packages {
		if name == pkgName {
			return pkg
		}
	}
	return Package{}
}

// overwriteConfigs adds all values from the addition config, and over-writes
// the original where duplicates exist
func overwriteConfigs(original TOMLConfig, addition TOMLConfig) TOMLConfig {
	if original.Packages == nil {
		original.Packages = map[string]Package{}
	}
	for pkgName, pkg := range addition.Packages {
		original.Packages[pkgName] = pkg
	}
	if original.InstallerDefs == nil {
		original.InstallerDefs = map[string]Installer{}
	}
	for installerName, installer := range addition.InstallerDefs {
		original.InstallerDefs[installerName] = installer
	}
	if original.Tasks == nil {
		original.Tasks = map[string]Task{}
	}
	for taskName, task := range addition.Tasks {
		original.Tasks[taskName] = task
	}
	return original
}

// combineConfigs adds all values from the addition config, but keeps originals where duplicates exist
func combineConfigs(original TOMLConfig, addition TOMLConfig) TOMLConfig {
	if original.Packages == nil {
		original.Packages = map[string]Package{}
	}
	for pkgName, pkg := range addition.Packages {
		if _, alreadyExists := original.Packages[pkgName]; !alreadyExists {
			original.Packages[pkgName] = pkg
		}
	}
	if original.InstallerDefs == nil {
		original.InstallerDefs = map[string]Installer{}
	}
	for installerName, installer := range addition.InstallerDefs {
		if _, alreadyExists := original.InstallerDefs[installerName]; !alreadyExists {
			original.InstallerDefs[installerName] = installer
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
