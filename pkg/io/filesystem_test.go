package io

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSymlink(t *testing.T) {
	fs := NewFilesystem()
	err := fs.Symlink("README.md", "/tmp/README.md")
	assert.NoError(t, err)
	isLink, err := fs.IsSymlinkTo("README.md", "/tmp/README.md")
	assert.NoError(t, err)
	assert.True(t, isLink)
}
