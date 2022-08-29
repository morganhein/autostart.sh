package manager

import (
	"github.com/morganhein/envy/pkg/io"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveRecipe(t *testing.T) {
	r, err := ResolveRecipe(io.NewFilesystem(), "../../configs/examples/package.toml")
	assert.NoError(t, err)
	assert.NotEmpty(t, r.Shells["asdf"])
}
