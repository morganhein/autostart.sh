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
