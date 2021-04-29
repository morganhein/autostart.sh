package autostart

import (
	"context"
	"fmt"
)

//Package is meant to represent both shell scripts and installable items from package managers
type Package struct {
	Description string
	Default     string                      //default name for this package
	InstallOpts map[string]PkgInstallOption `toml:"installers"`
	RestrictTo  []string                    //restricts installing this package to the specified set of installers, this
	//prevents "falling through" to try installing using the preferred package manager, and instead restricts to this method
}

type PkgInstallOption struct {
	Name string
	Deps []string
	URL  string //url of file to download
}

//InstallPackage tries to find a pre-defined package in the config, and runs it if so,
//otherwise it tries to install the package through the highest priority available package manager
func InstallPackage(ctx context.Context, config Config, pkg string) error {
	installer, opts, err := determineInstallerOptions(ctx, config, pkg)
	if err != nil {
		return err
	}
	//install deps
	for _, dep := range opts.Deps {
		err = InstallPackage(ctx, config, dep)
		if err != nil {
			return err
		}
	}
	//unless every command is run through the shell, then we need to grab all the environment variables
	//and inject them into the command somehow. No history should be available here unless we run through the shell
	//so there are tradoffs
	//map of environment variables to set in the shell for this command and their values
	//TODO (@morgan): need to prefix variables so we can avoid collisions and are easily identifiable
	setEnv := envVariables{}
	if opts.URL != "" {
		//download the URL/file to send to install command
		filename, err := Download(ctx, opts.URL, config.TmpDir)
		if err != nil {
			return err
		}
		setEnv.add("file", filename)
		//need to set variables in shell for filename location
	}
	//check if sudo, and set shell variable as well
	//run the cmds through a parser that replaces variables with values
	for k, cmd := range installer.Cmds {
		installer.Cmds[k] =
	}
	installer.Cmds
}

//determine installer and options
func determineInstallerOptions(ctx context.Context, config Config, pkg string) (*Installer, *PkgInstallOption, error) {
	//get the package options for the specified package, if available, otherwise just use a default
	returnOptions := PkgInstallOption{
		Name: pkg,
	}
	returnInstaller := config.Installers[config.InstallerPriority[0]]
	definedPkg, ok := config.Packages[pkg]
	//if we don't have package information for this pkg request, then just try to install it using
	//the highest priority package manager
	if !ok {
		return &returnInstaller, &returnOptions, nil
	}
	//no restrictions! easy peasy, see if there are any pkgOptions for our highest priority installer
	if len(definedPkg.RestrictTo) == 0 {
		opts, ok := definedPkg.InstallOpts[config.InstallerPriority[0]]
		var addOpts *PkgInstallOption
		if ok {
			addOpts = &opts
		}
		//combine defaults with package defaults with specific package installer options
		returnOptions = combineOpts(returnOptions, definedPkg.Default, addOpts)
		return &returnInstaller, &returnOptions, nil
	}
	//there are installer restrictions, let's try to fulfill them
	//check that the installer is enabled
	installerName := findOverlap(definedPkg.RestrictTo, config.InstallerPriority)
	if installerName == "" {
		//installers not enabled for this package
		return nil, nil, fmt.Errorf("package %v requires one of the following installers, which are not enabled: %v", pkg, definedPkg.RestrictTo)
	}
	returnInstaller = config.Installers[installerName]
	//now combine options for the found installer
	opts, ok := definedPkg.InstallOpts[installerName]
	var addOpts *PkgInstallOption
	if ok {
		addOpts = &opts
	}
	returnOptions = combineOpts(returnOptions, definedPkg.Default, addOpts)
	return &returnInstaller, &returnOptions, nil
}

func combineOpts(original PkgInstallOption, defaultName string, new *PkgInstallOption) PkgInstallOption {
	if defaultName != "" {
		original.Name = defaultName
	}
	if new == nil {
		return original
	}
	if new.Name != "" {
		original.Name = new.Name
	}
	original.Deps = new.Deps
	original.URL = new.URL
	return original
}

//findOverlap finds the first needle that exists in haystack
func findOverlap(needles []string, haystack []string) string {
	for _, needle := range needles {
		for _, hay := range haystack {
			if needle == hay {
				return needle
			}
		}
	}
	return ""
}

