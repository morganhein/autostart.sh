package autostart

import (
	"context"
	"errors"

	"github.com/BurntSushi/toml"
)

const installerDefaults = `
[installer.apt]
    run_if = ["which apt", "which apt-get"]
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
	Name   string
	SkipIf []string `toml:"skip_if"`
	RunIf  []string `toml:"run_if"`
	Sudo   bool
	Cmds   []string
}

func loadDefaultInstallers(ctx context.Context, config Config) (*Installer, error) {
	defaultConfig := &Config{}
	err := toml.Unmarshal([]byte(installerDefaults), defaultConfig)
	if err != nil {
		return nil, err
	}
	c := combineConfigs(*defaultConfig, config)
	//for each installer configured, determine if it's available
	return detectInstaller(ctx, c)
}

func detectInstaller(ctx context.Context, config Config) (*Installer, error) {
	for k, v := range config.Installers {
		sr := shouldRun(ctx, v.SkipIf, v.RunIf)
		if !sr {
			continue
		}
		v.Name = k
		return &v, nil
	}
	return nil, errors.New("unable to find a suitable installer")
}
