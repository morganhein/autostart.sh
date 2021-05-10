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
    cmd =  "${sudo} apt install -y ${pkg}"

[installer.brew]
    run_if = ["which brew"]
    sudo = false
    cmd =  "${sudo} brew install ${pkg}"

[installer.apk]
    run_if = ["which apk"]
    sudo = false
    cmd =  "${sudo} apk add ${pkg}"

[installer.dnf]
    run_if = ["which dnf"]
    sudo = true
    cmd =  "${sudo} dnf install -y ${pkg}"

[installer.pacman]
    run_if = ["which pacman"]
    sudo = true
    cmd =  "${sudo} pacman -Syu ${pkg}"

[installer.yay]
    run_if = ["which yay"]
    sudo = true
    cmd =  "${sudo} yay -Syu ${pkg}"

[installer.sh]
    run_if = ["which sh"]
    sudo = true
    cmd =  "${sudo} sh ${pkg}"

[installer.bash]
    run_if = ["which bash"]
    sudo = true
    cmd =  "${sudo} bash ${pkg}"
`

type Installer struct {
	Name   string
	SkipIf []string `toml:"skip_if"`
	RunIf  []string `toml:"run_if"`
	Sudo   bool
	Cmd    string
}

func loadDefaultInstallers(config Config) (Config, error) {
	defaultConfig := &Config{}
	err := toml.Unmarshal([]byte(installerDefaults), defaultConfig)
	if err != nil {
		return config, err
	}
	return combineConfigs(config, *defaultConfig), nil
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
