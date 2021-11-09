package autostart

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

//editor is meant to facilitate variable substitution in command lines

func injectPackage(cmdLine, installCmd, pkg string) string {
	if len(pkg) == 0 {
		return cmdLine
	}
	if strings.HasPrefix(cmdLine, "@install") {
		return strings.Replace(installCmd, "${pkg}", pkg, 1)
	}
	return strings.Replace(cmdLine, "${pkg}", pkg, -1)
}

//injectVars first tries to replace all ${SH} style variables with the ASH configuration values,
//then with any environment variables.
func injectVars(cmdLine string, vars envVariables, sudo bool) string {
	//need to do sudo expansion first, since it's a special case
	if sudo {
		cmdLine = strings.Replace(cmdLine, "${sudo}", "sudo", -1)
		cmdLine = strings.Replace(cmdLine, "${SUDO}", "sudo", -1)
	} else {
		cmdLine = strings.Replace(cmdLine, "${sudo}", "", -1)
		cmdLine = strings.Replace(cmdLine, "${SUDO}", "", -1)
	}
	for k, v := range vars {
		cmdLine = strings.Replace(cmdLine, fmt.Sprintf("${%v}", strings.ToUpper(k)), v, -1)
		cmdLine = strings.Replace(cmdLine, fmt.Sprintf("${%v}", strings.ToLower(k)), v, -1)
	}

	//now search for any leftover requests intended to get environment variables
	//regular expressions...ewwwww you say.... But I like them!
	reg := regexp.MustCompile(`\${(\w+)}`)
	matches := reg.FindAllStringSubmatch(cmdLine, -1)
	if matches == nil {
		return cmdLine
	}
	for _, match := range matches {
		//try to get the environment variable defined here
		v := os.Getenv(match[1])
		if v == "" {
			//TODO (@morgan): possibly warning here that a leftover variable expression did not get expanded
			continue
		}
		cmdLine = strings.Replace(cmdLine, match[0], v, -1)
	}
	return cmdLine
}

func clean(input string) string {
	input = strings.TrimSpace(input)
	return strings.ToLower(input)
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

// Chooses what actual package name to use, and which installer
//TODO (@morgan): this needs to support a preferred hierarchy of installers?
func determinePackageOptions(pkgName string, config Config, installer Installer) PkgInstallOption {
	res := PkgInstallOption{
		Name: pkgName,
	}
	p, ok := config.Packages[pkgName]
	if !ok {
		return res
	}
	if opt, ok := p[installer.Name]; ok {
		res.Name = opt
	}
	return res
}

func getPackageName(config Config, unparsedPkg string) string {
	//determine if line contains any pkg placeholders denoted by ^pkg
	//only a single replacement happens at once
	reg := regexp.MustCompile(`\^(\w+)`)
	matches := reg.FindAllStringSubmatch(unparsedPkg, 1)
	if matches != nil {
		//do package gathering and replacement here
		for _, match := range matches {
			opt := determinePackageOptions(match[1], config, config.Installer)
			return opt.Name
		}
	}
	//no '^pkg' was found, so try to return 'pkg'
	return unparsedPkg
}
