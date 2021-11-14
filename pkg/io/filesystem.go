package io

import (
	"os"
	"path/filepath"

	"github.com/morganhein/autostart.sh/pkg/T"
)

type Filesystem interface {
	Symlink(from, to string) error
	Stat(name string) (os.FileInfo, error)
	// IsSymlinkTo detects if the file at `from` symlinks to `to`
	IsSymlinkTo(from, to string) (bool, error)
}

func NewFilesystem() *filesystem {
	return &filesystem{}
}

type filesystem struct{}

func (f filesystem) IsSymlinkTo(from, to string) (bool, error) {
	stat, err := os.Lstat(from)
	if err != nil {
		return false, T.Log(err)
	}
	if stat.Mode()&os.ModeSymlink == 0 {
		//not a symlink
		return false, nil
	}
	ogFile, err := filepath.EvalSymlinks(from)
	if err != nil {
		return false, T.Log(err)
	}
	return to == ogFile, nil
}

func (f filesystem) Symlink(from, to string) error {
	return os.Symlink(from, to)
}

func (f filesystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
