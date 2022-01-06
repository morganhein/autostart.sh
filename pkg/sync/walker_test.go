package sync

import (
	"github.com/morganhein/autostart.sh/pkg/io"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"testing"
)

func TestSourceToTargetHelperMissingSource(t *testing.T) {
	fs := &io.FilesystemMock{
		StatFunc: func(name string) (fs.FileInfo, error) {
			return nil, os.ErrNotExist
		},
	}
	w := walker{
		fs:         fs,
		log:        io.NewLogger(),
		baseSource: "/home/tester/.repo",
		baseTarget: "/home/tester/",
		issues:     []Mismatch{},
	}
	err := w.sourceToTargetHelper("/home/tester/.repo/missing_file")
	assert.NoError(t, err)
	assert.Len(t, w.issues, 1)
	assert.Equal(t, MissingFromTarget, w.issues[0].Issue)
}

func TestSourceToTargetHelperAlreadyLinked(t *testing.T) {
	fs := &io.FilesystemMock{
		StatFunc: func(name string) (fs.FileInfo, error) {
			return nil, nil
		},
		IsSymlinkToFunc: func(from string, to string) (bool, error) {
			return true, nil
		},
	}
	w := walker{
		fs:         fs,
		log:        io.NewLogger(),
		baseSource: "/home/tester/.repo",
		baseTarget: "/home/tester/",
		issues:     []Mismatch{},
	}
	err := w.sourceToTargetHelper("/home/tester/.repo/already_symlinked")
	assert.NoError(t, err)
	assert.Len(t, w.issues, 0)
}

func TestSourceToTargetHelperFileCollision(t *testing.T) {
	fs := &io.FilesystemMock{
		StatFunc: func(name string) (fs.FileInfo, error) {
			return nil, nil
		},
		IsSymlinkToFunc: func(from string, to string) (bool, error) {
			return false, nil
		},
	}
	w := walker{
		fs:         fs,
		log:        io.NewLogger(),
		baseSource: "/home/tester/.repo",
		baseTarget: "/home/tester/",
		issues:     []Mismatch{},
	}
	err := w.sourceToTargetHelper("/home/tester/.repo/missing_file")
	assert.NoError(t, err)
	assert.Len(t, w.issues, 1)
	assert.Equal(t, FileCollision, w.issues[0].Issue)
}
