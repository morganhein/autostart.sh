package autostart

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDefaultInstallers(t *testing.T) {
	ctx := context.Background()
	installers, err := loadDefaultInstallers(ctx, Config{})
	assert.NoError(t, err)
	assert.NotNil(t, installers)
}

func TestCombineInstallers(t *testing.T) {
	ctx := context.Background()
	c := Config{
		Installers: map[string]Installer{
			"TEST": {
				Name:   "TEST_PKG_MANAGER",
				SkipIf: nil,
				RunIf:  []string{"which ls"}, //assumed that LS exists pretty much everywhere
				Sudo:   false,
				Cmds:   nil,
			},
		},
	}
	installers, err := loadDefaultInstallers(ctx, c)
	assert.NoError(t, err)
	assert.NotNil(t, installers)
	assert.Equal(t, "TEST_PKG_MANAGER", c.Installers["TEST"].Name)
}
