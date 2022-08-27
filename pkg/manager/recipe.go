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
	General       General              `toml:"general"`
	Packages      map[string]Package   `toml:"pkg"`
	InstallerDefs map[string]Installer `toml:"installer"`
	Tasks         map[string]Task      `toml:"task"`
}

// The General section of a TOML config
type General struct {
	AllowedInstallers []string
	ConfigDir         string `toml:"config_dir"`
	HomeDir           string `toml:"home_dir"`
}

// A task as define in a TOML config
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

// An installer definition from a TOML config
type Installer struct {
	Name    string
	RunIf   []string `toml:"run_if"`
	Sudo    bool
	Cmd     string
	Update  string
	Updated bool
}

// A package alias as defined in a TOML config
// It translates a common name like "vim" to the
// package name for the specific installer.
type Package map[string]string

type PkgInstallOption struct {
	Name string
}

func ResolveRecipe(fs io.Filesystem, configLocation string) (*Recipe, error) {
	recipes, err := loadAllRecipes(fs, configLocation)
	if err != nil {
		return nil, err
	}
	var r Recipe
	// for each config loaded, compose them
	for _, c := range recipes {
		r = overwriteRecipe(r, c)
	}
	return &r, nil
}

func loadAllRecipes(fs io.Filesystem, configLocation string) ([]Recipe, error) {
	var recipes []Recipe
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	locations := []string{
		"/usr/share/envy/default.toml",
		"$HOME/.envy/default.toml",
		"$HOME/.config/envy/default.toml",
		"$HOME/.envy/config.toml",
		"$HOME/.config/envy/config.toml",
		configLocation,
	}
	for _, loc := range locations {
		c, err := loadRecipeFromFS(fs, strings.Replace(loc, "$HOME", home, 1))
		if err == nil {
			recipes = append(recipes, *c)
		}
	}
	if len(locations) == 0 {
		return nil, xerrors.New("No configuration files found, none could be loaded.")
	}
	return recipes, nil
}

func loadRecipeFromFS(fs io.Filesystem, location string) (*Recipe, error) {
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
// TODO: clean this up. Why is it always returning a package? how is this used? should we signal back when we don't have a package?
func getPackage(config Recipe, name string) Package {
	for pkgName, pkg := range config.Packages {
		if name == pkgName {
			return pkg
		}
	}
	return Package{}
}

// overwriteRecipe adds all values from the addition recipe, and over-writes
// the original where duplicates exist
func overwriteRecipe(original Recipe, addition Recipe) Recipe {
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

// combineRecipe adds all values from the addition config, but keeps originals where duplicates exist
func combineRecipe(original Recipe, addition Recipe) Recipe {
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
