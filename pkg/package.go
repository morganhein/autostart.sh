package autostart

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type Package map[string]string

type PkgInstallOption struct {
	Name string
}

func resolveCmdLine(cmdLine string, config Config, vars envVariables) string {
	//determine if line contains any pkgs denoted by #{name}
	reg := regexp.MustCompile(`#{(\w+)}`)
	matches := reg.FindAllStringSubmatch(cmdLine, -1)
	if matches != nil {
		//do package gathering and replacement here
		for _, match := range matches {
			opt := determinePackageOptions(match[1], config, config.Installer)
			cmdLine = strings.Replace(cmdLine, match[0], opt.Name, -1)
		}
	}
	cmdLine = injectMacros(cmdLine, config)
	//do we sudo, or do we not?
	sudo := determineSudo(config, config.Installer)
	if vars == nil {
		vars = envVariables{}
	}
	cmdLine = injectVars(cmdLine, vars, sudo)
	return cmdLine
}

func determineSudo(config Config, installer Installer) bool {
	if clean(config.Sudo) == "true" {
		return true
	}
	if clean(config.Sudo) == "false" {
		return false
	}
	return installer.Sudo
}

func determinePackageOptions(pkgName string, config Config, installer Installer) PkgInstallOption {
	res := PkgInstallOption{
		Name: pkgName,
	}
	p, ok := config.Packages[pkgName]
	if !ok {
		return res
	}
	if def, ok := p["default"]; ok {
		res.Name = def
	}
	if opt, ok := p[installer.Name]; ok {
		res.Name = opt
	}
	return res
}

func installPackage(ctx context.Context, config Config, cmdLine string) error {
	cmdLine = resolveCmdLine(cmdLine, config, nil)
	out, err := runCmd(ctx, config.DryRun, cmdLine)
	if err != nil {
		return err
	}
	if config.Verbose {
		fmt.Println(out)
	}
	return nil
}
