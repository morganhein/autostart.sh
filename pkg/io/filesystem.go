package io

import (
	"errors"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"

	"github.com/morganhein/shoelace/pkg/oops"
)

type Filesystem interface {
	//Detects if the file already links to the appropriate location, which becomes a NOP.
	//If a file exists at the target location, it gets moved to the backup location
	CreateSymlink(from, to, backup string) error
	Stat(name string) (os.FileInfo, error)
	// IsSymlinkTo detects if the file at `from` symlinks to `to`
	IsSymlinkTo(from, to string) (bool, error)
	//Move(from, to string) error
}

func NewFilesystem() *filesystem {
	return &filesystem{}
}

type filesystem struct{}

func (f filesystem) CreateSymlink(from, to, backup string) error {
	//detect if file is symlink to correct location already
	alreadyGood := func() bool {
		stat, err := os.Lstat(to)
		if err != nil {
			return false
		}
		if !(stat.Mode()&os.ModeSymlink != 0) {
			return false
		}
		ogFile, err := os.Readlink(stat.Name())
		return from == ogFile
	}()

	if alreadyGood {
		return nil
	}

	//try moving the file, ignore fileNotFound error
	err := os.Rename(to, backup)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return xerrors.Errorf("error moving old file from %v to %v: %v", to, backup, err)
	}

	//make new symlink
	return os.Symlink(from, to)
}

func (f filesystem) Move(from, to string) error {
	panic("implement me")
}

func (f filesystem) IsSymlinkTo(from, to string) (bool, error) {
	stat, err := os.Lstat(from)
	if err != nil {
		return false, oops.Log(err)
	}
	if stat.Mode()&os.ModeSymlink == 0 {
		//not a symlink
		return false, nil
	}
	ogFile, err := filepath.EvalSymlinks(from)
	if err != nil {
		return false, oops.Log(err)
	}
	return to == ogFile, nil
}

func (f filesystem) Symlink(from, to string) error {
	return os.Symlink(from, to)
}

func (f filesystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
