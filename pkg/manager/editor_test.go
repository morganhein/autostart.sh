package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Tests that the appropriate package name and installer are chosen
func TestDeterminePackageOptions(t *testing.T) {
	c := Config{
		RunningConfig: RunningConfig{},
		Packages: map[string]Package{
			"golang": map[string]string{
				"brew": "golang",
				"apk":  "golang_apk",
			},
		},
		Installers: map[string]Installer{
			"brew": {
				Name:   "brew",
				SkipIf: nil,
				RunIf:  nil,
				Sudo:   false,
				Cmd:    "${sudo} brew install ${pkg}",
			},
			"apk": {
				Name:   "apk",
				SkipIf: nil,
				RunIf:  nil,
				Sudo:   false,
				Cmd:    "${sudo} apk add ${pkg}",
			},
		},
		Tasks: nil,
	}

	opts := determinePackageOptions("golang", c, c.Installers["brew"])
	assert.Equal(t, "golang", opts.Name)

	opts = determinePackageOptions("golang", c, c.Installers["apk"])
	assert.Equal(t, "golang_apk", opts.Name)
}
