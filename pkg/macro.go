package autostart

import (
	"fmt"
	"strings"
)

type Macro struct {
	Insert string //what the macro expands out to/replaces
	Deps   []string
	Prefix []string
	Pre    []string
	Post   []string
}

func injectMacros(cmdLine string, config Config) string {
	install := "@install"
	cmdLine = strings.Replace(cmdLine, install, config.Installer.Cmd, -1)
	for k, v := range config.Macros {
		macro := fmt.Sprintf("@%v", k)
		//detect appropriate version of macro, depending on installer/target etc
		if strings.Contains(cmdLine, macro) {
			//is there a installer specific version of this macro?
			if im, ok := config.Macros[fmt.Sprintf("macro__%v", config.Installer.Name)]; ok {
				cmdLine = strings.Replace(cmdLine, macro, im.Insert, -1)
				continue
			}
			//otherwise default to basic macro
			cmdLine = strings.Replace(cmdLine, macro, v.Insert, -1)
		}
	}
	return cmdLine
}
