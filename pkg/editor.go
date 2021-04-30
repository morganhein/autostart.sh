package autostart

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

//editor is used for translating configuration strings with replacement variables into strings that contain
//the desired value

type envVariables map[string]string

const (
	ORIGINAL_TASK = "ORIGINAL_TASK"
	CURRENT_TASK  = "CURRENT_TASK"
	CURRENT_PKG   = "CURRENT_PKG"
	SUDO          = "SUDO"
	CONFIG_PATH   = "CONFIG_PATH"
)

func (e envVariables) copy() envVariables {
	//TODO (@morgan): I think this can be copied more efficiently
	newEnv := envVariables{}
	for k, v := range e {
		newEnv[k] = v
	}
	return newEnv
}

//set default environment variables
func hydrateEnvironment(config Config, env envVariables) {
	env[ORIGINAL_TASK] = config.Task
	env[CONFIG_PATH] = config.ConfigLocation
	env[ORIGINAL_TASK] = config.Task
	//possibly add link src and dst links here
}

//injectVars first tries to replace all ${SH} style variables with the ASH configuration values,
//then with any environment variables.
func injectVars(cmdLine string, vars envVariables, sudo bool) string {
	//need to do sudo expansion first, since it's a special case
	if sudo {
		cmdLine = strings.Replace(cmdLine, "${sudo}", "sudo", -1)
		cmdLine = strings.Replace(cmdLine, "${SUDO}", "sudo", -1)
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
