//go:build ubuntu
// +build ubuntu

package tests

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type UbuntuSuite struct {
	suite.Suite
}

func TestUbuntuSuite(t *testing.T) {
	suite.Run(t, new(UbuntuSuite))
}

//Puts the config file in /usr/var/shoelace/default.toml and defaults are loaded
func (u *UbuntuSuite) TestLoadConfigFromUsrDefault() {
	defaultLocation := "/usr/share/shoelace/default.toml"
	err := os.Mkdir("/usr/share/shoelace", os.ModeDir)
	u.NoError(err)
	_, err = copy("../configs/default.toml", defaultLocation)
	u.NoError(err)
	e, err := exists(defaultLocation)
	u.NoError(err)
	u.True(e)
}
