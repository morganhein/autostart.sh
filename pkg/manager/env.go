package manager

import (
	"context"
	"path"

	"golang.org/x/xerrors"
)

type Operation string

const (
	SYNC    Operation = "sync"
	INSTALL Operation = "install"
	TASK    Operation = "task"
)

const (
	ORIGINAL_TASK = "ORIGINAL_TASK"
	CURRENT_TASK  = "CURRENT_TASK"
	CURRENT_PKG   = "CURRENT_PKG"
	SUDO          = "SUDO"
	CONFIG_PATH   = "CONFIG_PATH"
	TARGET_PATH   = "TARGET_PATH"
	SOURCE_PATH   = "SOURCE_PATH"
)

// determineBestAvailableInstaller determines installer based on following precedence:
// 1. Installer specified by command line
// 2. Package has a preferred installer method
// 3. First available installer that is supported by the pkg
func determineBestAvailableInstaller(ctx context.Context, config RunConfig, pkg Package, d Decider) (*Installer, error) {
	availableInstallers := determineAvailableInstallers(ctx, config.Recipe.InstallerDefs, d)
	//if execution arguments have forced a specific installer to be used
	if config.ForceInstaller != "" {
		i, ok := config.Recipe.InstallerDefs[config.ForceInstaller]
		if ok && isAvailableInstaller(i, availableInstallers) {
			i.Name = config.ForceInstaller
			return &i, nil
		}
		return nil, xerrors.Errorf("an installer was requested (%v), but was not found", config.ForceInstaller)
	}
	// if preferred installer is available, use it
	if requiredInstaller, ok := pkg["prefer"]; ok {
		i, ok := config.Recipe.InstallerDefs[requiredInstaller]
		if ok {
			i.Name = requiredInstaller
			return &i, nil
		}
		return nil, xerrors.Errorf("an installer was requested (%v), but was not found", requiredInstaller)
	}
	if len(config.Recipe.General.InstallerPreferences) > 0 {
		for _, v := range config.Recipe.General.InstallerPreferences {
			for _, availableInstaller := range availableInstallers {
				if v == availableInstaller.Name {
					return &availableInstaller, nil
				}
			}
		}
		return nil, xerrors.Errorf("preferred installer(fs) are not available (%+v)", config.Recipe.General.InstallerPreferences)
	}

	//no installer preferred, grab the first available one
	for _, installer := range availableInstallers {
		return &installer, nil
	}

	return nil, xerrors.New("unable to find a suitable installer")
}

func determineAvailableInstallers(ctx context.Context, definedInstallers map[string]Installer, d Decider) []Installer {
	var availableInstallers []Installer
	for installerName, installer := range definedInstallers {
		sr := d.ShouldRun(ctx, []string{}, installer.RunIf)
		if !sr {
			continue
		}
		installer.Name = installerName
		availableInstallers = append(availableInstallers, installer)
	}
	return availableInstallers
}

func isAvailableInstaller(needle Installer, haystack []Installer) bool {
	for _, v := range haystack {
		if v.Name == needle.Name {
			return true
		}
	}
	return false
}

// Environment Variables

type envVariables map[string]string

func (e envVariables) copy() envVariables {
	//TODO (@morgan): I think this can be copied more efficiently
	newEnv := envVariables{}
	for k, v := range e {
		newEnv[k] = v
	}
	return newEnv
}

// set default environment variables
func hydrateEnvironment(config RunConfig, env envVariables) {
	env[ORIGINAL_TASK] = config.originalTask
	env[CONFIG_PATH] = path.Dir(config.RecipeLocation)
	//possibly add link src and dst links here
}
