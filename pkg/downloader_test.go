package autostart

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloader(t *testing.T) {
	location, err := Download(context.Background(), "http://localhost:8090/shell", "/tmp")
	assert.NoError(t, err)
	fmt.Println(location)
}
