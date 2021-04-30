package autostart

import (
	"context"
	"regexp"
	"strings"
)

type Package map[string]string

type PkgInstallOption struct {
	Name string
}

func resolveCmdLine(cmdLine string, config Config, vars envVariables, installer Installer, opt PkgInstallOption) string {
	//determine if line contains any pkgs denoted by #{name}
	reg := regexp.MustCompile(`#{(\w+)}`)
	matches := reg.FindAllStringSubmatch(cmdLine, -1)
	if matches != nil {
		//do package gathering and replacement here
		for _, match := range matches {
			cmdLine = strings.Replace(cmdLine, match[1], opt.Name, -1)
		}
	}
	//do we sudo, or do we not?
	sudo := determineSudo(config, installer)
	cmdLine = injectVars(cmdLine, vars, sudo)
	cmdLine = injectMacros(cmdLine, config)
	return cmdLine
}

func determineSudo(config Config, installer Installer) bool {
	if clean(config.Sudo) == "true" {
		return true
	}
	if clean(config.Sudo) == "false" {
		return false
	}
	return true
}

//installPackage tries to find a pre-defined package in the config, and runs it if found.
//Otherwise it tries to install the package through the highest priority available package manager
func installPackage(ctx context.Context, config Config, pkg string) error {

	//map of environment variables and their values to replace inplace in the command
	//check if sudo, and set shell variable as well
	//run the cmds through a parser that replaces variables with values
	//for k, cmd := range installer.Cmds {
	//	installer.Cmds[k] =
	//}
	//installer.Cmds
}

func determineInstallerOptions(config Config, installerName string, pkg string) PkgInstallOption {
	//get the package options for the specified package, if available, otherwise just use a default
	returnOptions := PkgInstallOption{
		Name: pkg,
	}
	definedPkg, ok := config.Packages[pkg]
	//if we don't have package information for this pkg request, then just try to install it using
	//the requested package manager
	if !ok {
		return returnOptions
	}
	opts, ok := definedPkg.InstallOpts[installerName]
	var addOpts *PkgInstallOption
	if ok {
		addOpts = &opts
	}
	//combine defaults with package defaults with specific package installer options
	return combineOpts(returnOptions, definedPkg.Default, addOpts)
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
	return original
}
