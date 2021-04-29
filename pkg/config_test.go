package autostart

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAndReadConfig(t *testing.T) {
	strConfig, err := LoadPackageConfig("../packages.toml")
	assert.NoError(t, err)
	config, err := ParsePackageConfig(strConfig)
	assert.NoError(t, err)
	assert.NotNil(t, config)
}
