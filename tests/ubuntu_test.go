//go:build ubuntu
// +build ubuntu

package tests

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//Puts the config file in /usr/var/shoelace/default.toml and defaults are loaded
func TestLoadConfigFromUsrDefault(t *testing.T) {
	defaultLocation := "/usr/share/shoelace/default.toml"
	err := os.Mkdir("/usr/share/shoelace", os.ModeDir)
	assert.NoError(t, err)
	_, err = copy("../configs/default.toml", defaultLocation)
	assert.NoError(t, err)
	e, err := exists(defaultLocation)
	assert.NoError(t, err)
	assert.True(t, e)
}

func TestDoNewTest(t *testing.T) {

}
