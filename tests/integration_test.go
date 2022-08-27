//go:build integrated
// +build integrated

package tests

import (
	"github.com/morganhein/envy/pkg/io"
	"github.com/morganhein/envy/pkg/manager"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

// Puts the config file in /usr/var/envy/default.toml and defaults are loaded
func TestLoadConfigFromUsrDefault(t *testing.T) {
	defaultLocation := "/usr/share/envy/default.toml"
	err := os.Mkdir("/usr/share/envy", os.ModeDir)
	assert.NoError(t, err)
	_, err = copy("../configs/default.toml", defaultLocation)
	assert.NoError(t, err)
	e, err := exists(defaultLocation)
	assert.NoError(t, err)
	assert.True(t, e)

	r, err := manager.ResolveRecipe(io.NewFilesystem(), "")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Contains(t, r.InstallerDefs, "apt")
}

//
////Puts config in $HOME/.config/envy/default.toml and defaults are loaded
//func TestLoadDefaultConfigFromHomeConfig(t *testing.T) {
//	homeDir, err := os.UserHomeDir()
//	assert.NoError(t, err)
//	homeConfigLocation := fmt.Sprintf("%v/.config/envy/ubuntu.toml", homeDir)
//	err = os.MkdirAll(fmt.Sprintf("%v/.config/envy/", homeDir), os.ModeDir)
//	assert.NoError(t, err)
//	_, err = copy("../configs/default.toml", homeConfigLocation)
//	assert.NoError(t, err)
//	e, err := exists(homeConfigLocation)
//	assert.NoError(t, err)
//	assert.True(t, e)
//
//	cmd.Execute()
//}

func TestWhich(t *testing.T) {
	r, err := io.CreateShell()
	assert.NoError(t, err)
	ctx, cancel := newCtx(10 * time.Second)
	//assert we get a known positive
	exists, out, err := r.Which(ctx, "ls")
	cancel()
	assert.NoError(t, err, out)
	assert.True(t, true)

	//assert we get a known negative
	exists, out, err = r.Which(ctx, "monkey-pox-and-covid-suck")
	cancel()
	assert.Error(t, err, out)
	assert.False(t, exists)
}

func TestInstallCommandInstallsPackage(t *testing.T) {
	sh, err := io.CreateShell()
	assert.NoError(t, err)
	ctx, cancel := newCtx(10 * time.Second)
	//assert vim doesn't already exist
	exists, out, err := sh.Which(ctx, "vim")
	cancel()
	assert.Error(t, err, out)
	assert.False(t, exists)

	//install it
	ctx, cancel = newCtx(10 * time.Second)
	mgr := manager.New(io.NewFilesystem(), sh)
	appConfig := manager.RunConfig{
		RecipeLocation: "../configs/default.toml",
		Operation:      manager.INSTALL,
		Sudo:           "false",
		Verbose:        false,
	}
	err = mgr.Start(ctx, appConfig, "vim")
	cancel()
	assert.NoError(t, err)

	//assert vim exists
	exists, out, err = sh.Which(ctx, "vim")
	cancel()
	assert.NoError(t, err)
	assert.True(t, exists, out)
}

func TestTaskInstallsPackageCorrectly(t *testing.T) {
	//copy default installers first
	defaultLocation := "/usr/share/envy/default.toml"
	_, err := copy("../configs/default.toml", defaultLocation)
	assert.NoError(t, err)

	// make shell
	sh, err := io.CreateShell()
	assert.NoError(t, err)
	ctx, cancel := newCtx(10 * time.Second)

	//assert vim doesn't already exist
	exists, out, err := sh.Which(ctx, "vim")
	cancel()
	assert.Error(t, err, out)
	assert.False(t, exists)

	//install it
	ctx, cancel = newCtx(10 * time.Second)
	mgr := manager.New(io.NewFilesystem(), sh)
	appConfig := manager.RunConfig{
		RecipeLocation: "configs/simple_task.toml",
		Operation:      manager.TASK,
		Sudo:           "false",
		Verbose:        false,
	}
	err = mgr.Start(ctx, appConfig, "vim")
	cancel()
	assert.NoError(t, err)

	//assert vim exists
	exists, out, err = sh.Which(ctx, "vim")
	cancel()
	assert.NoError(t, err)
	assert.True(t, exists, out)
}
