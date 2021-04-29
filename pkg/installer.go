package autostart

import (
	"context"

	"github.com/BurntSushi/toml"
)

const installerDefaults = `
[installer.apt]
    run_if = ["which apt"]
    sudo = true
    cmds = ["${sudo} apt install -y ${pkg}"]

[installer.brew]
    run_if = ["which brew"]
    sudo = false
    cmds = ["${sudo} brew install ${pkg}"]

[installer.apk]
    run_if = ["which apk"]
    sudo = false
    cmds = ["${sudo} apk add ${pkg}"]

[installer.dnf]
    run_if = ["which dnf"]
    sudo = true
    cmds = ["${sudo} dnf install -y ${pkg}"]

[installer.pacman]
    run_if = ["which pacman"]
    sudo = true
    cmds = ["${sudo} pacman -Syu ${pkg}"]

[installer.yay]
    run_if = ["which yay"]
    sudo = true
    cmds = ["${sudo} yay -Syu ${pkg}"]

[installer.sh]
    run_if = ["which sh"]
    sudo = true
    cmds = ["${sudo} sh ${pkg}"]

[installer.bash]
    run_if = ["which bash"]
    sudo = true
    cmds = ["${sudo} bash ${pkg}"]
`

type Installer struct {
	SkipIf []string `toml:"skip_if"`
	RunIf  []string `toml:"run_if"`
	Sudo   bool
	Cmds   []string
}

func loadDefaultInstallers(ctx context.Context, config Config) (map[string]Installer, error) {
	defaultConfig := &Config{}
	err := toml.Unmarshal([]byte(installerDefaults), defaultConfig)
	if err != nil {
		return nil, err
	}
	c := combineConfigs(*defaultConfig, config)
	//for each installer configured, determine if it's available
	return detectInstallers(ctx, c)
}

func detectInstallers(ctx context.Context, config Config) (map[string]Installer, error) {
	installers := map[string]Installer{}
	for k, v := range config.Installers {
		// compare runCmd-if
		err := testIf(ctx, v.RunIf)
		if len(v.RunIf) > 0 && err != nil {
			continue
		}
		// compare skip-if
		err = testIf(ctx, v.SkipIf)
		if len(v.SkipIf) > 0 && err == nil {
			continue
		}
		installers[k] = v
	}
	return installers, nil
}
