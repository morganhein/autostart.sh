package io

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloader(t *testing.T) {
	d := downloader{}
	location, err := d.Download(context.Background(), "http://localhost:8090/shell", "/tmp")
	assert.NoError(t, err)
	fmt.Println(location)
}
