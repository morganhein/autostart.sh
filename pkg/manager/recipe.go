package manager

import (
	"errors"
	"golang.org/x/xerrors"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/morganhein/envy/pkg/io"
)

type Recipe struct {
	// ForceInstaller will try to force the specified installer
	ForceInstaller string
	// Sudo will force using sudo when performing commands
	Sudo          string
	General       General              `toml:"general"`
	Packages      map[string]Package   `toml:"pkg"`
	InstallerDefs map[string]Installer `toml:"installer"`
	Tasks         map[string]Task      `toml:"task"`
}

func ResolveConfig(fs io.Filesystem, configLocation string) (*Recipe, error) {
	cfgs, err := loadAllConfigs(fs, configLocation)
	if err != nil {
		return nil, err
	}
	var r Recipe
	// for each config loaded, compose them
	for _, c := range cfgs {
		r = overwriteConfigs(r, c)
	}
	return &r, nil
}

func loadAllConfigs(fs io.Filesystem, configLocation string) ([]Recipe, error) {
	var configs []Recipe
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	locations := []string{
		"/usr/share/envy/default.toml",
		"$HOME/.envy/config.toml", "$HOME/.envy/default.toml",
		"$HOME/.config/envy/config.toml", "$HOME/.config/envy/default.toml",
		configLocation,
	}
	for _, loc := range locations {
		c, err := loadPackageConfigHelper(fs, strings.Replace(loc, "$HOME", home, 1))
		if err == nil {
			configs = append(configs, *c)
		}
	}
	if len(locations) == 0 {
		return nil, xerrors.New("No configuration files found, none could be loaded.")
	}
	return configs, nil
}

func loadPackageConfigHelper(fs io.Filesystem, location string) (*Recipe, error) {
	if location == "" {
		return nil, errors.New("config location is empty")
	}
	f, err := fs.ReadFile(location)
	if err != nil {
		return nil, err
	}
	k := &Recipe{}
	_, err = toml.Decode(string(f), k)
	if err != nil {
		return nil, err
	}
	return k, nil
}

// Finds the package <name> in the config if found, otherwise returns package with default settings matching <name>
func getPackage(config Recipe, name string) Package {
	for pkgName, pkg := range config.Packages {
		if name == pkgName {
			return pkg
		}
	}
	return Package{}
}

// overwriteConfigs adds all values from the addition config, and over-writes
// the original where duplicates exist
func overwriteConfigs(original Recipe, addition Recipe) Recipe {
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
func combineConfigs(original Recipe, addition Recipe) Recipe {
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
