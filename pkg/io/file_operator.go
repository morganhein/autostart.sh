package io

import (
	"errors"
	"os"
)

type Symlinker interface {
	//Detects if the file already links to the appropriate location, which becomes a NOP.
	//If a file exists at the target location, it gets moved to the backup location
	CreateSymlink(from, to, backup string) error
}

var _ Symlinker = (*symlinker)(nil)

func NewFileOperator() *symlinker {
	return &symlinker{}
}

type symlinker struct{}

func (f symlinker) CreateSymlink(from, to, backup string) error {
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
		return err
	}

	//make new symlink
	return os.Symlink(from, to)
}
