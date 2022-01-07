package manager

import (
	"testing"
)

//Tests that the appropriate package name and installer are chosen
func TestDeterminePackageOptions(t *testing.T) {
	_ = Config{
		RunningConfig: RunningConfig{},
		Packages: map[string]Package{
			"golang": map[string]string{
				"brew": "golang",
				"apk":  "golang_apk",
			},
		},
		Installers: map[string]Installer{
			"brew": {
				Name:  "brew",
				RunIf: nil,
				Sudo:  false,
				Cmd:   "${sudo} brew install ${pkg}",
			},
			"apk": {
				Name:  "apk",
				RunIf: nil,
				Sudo:  false,
				Cmd:   "${sudo} apk add ${pkg}",
			},
		},
		Tasks: nil,
	}
	panic("implement me")
	//opts := determinePackageOptions("golang", c, c.Installers["brew"])
	//assert.Equal(t, "golang", opts.Name)
	//
	//opts = determinePackageOptions("golang", c, c.Installers["apk"])
	//assert.Equal(t, "golang_apk", opts.Name)
}
