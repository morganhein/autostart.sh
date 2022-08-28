package sync

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/karrick/godirwalk"
	"github.com/morganhein/envy/pkg/io"
	"golang.org/x/xerrors"
)

// Sync files from the dotfile cache into the home/root targets. Conflicts will be resolved by prompting the user.

/*
Terms:
	$home_source is the path where actual dotfiles/configs exist on disk, and are symlinked from
	$home_target is the location to symlink into for $home files
	$root_source is the same as $home_source but for / files
	$root_target same as $home_target, but for / files
    $config_path is the full path and filename for the config file
	$ignores are the list of files, fullpath or pattern, to ignore

For every inner union of files in $home_source and $home_target, make sure it is a link from source to the target.
	May have to resolve collisions
For every outer left union of $home_source and $home_target, make a symlink from source to target
For every outer right union of $home_source and $home_target, ask to move file to source and symlink or ignore (this means we need to keep state?)

For now we won't into subfolders, instead the entire folder will just be moved and linked

Folder structure in source will be .repo/
										/home
										/extra
*/

type SyncConfig struct {
	Source  string
	Target  string
	DryRun  bool
	Ignores []string
	syncer  Syncer
	term    io.Terminal
}

func Sync(config SyncConfig) error {
	config.syncer = NewSyncer(io.NewFilesystem())
	config.term = io.NewTerminal()
	ctx := context.Background()
	return syncHelper(ctx, config)
}

func syncHelper(ctx context.Context, config SyncConfig) error {
	//if we don't already have an ignores (or a config isn't set, or something), then do a first pass to skip certain folders
	dirs, err := config.syncer.GatherDirs(ctx, config.Target)
	if err != nil {
		return err
	}

	selectedDirs := []string{}
	//now ask which ones we want to keep
	prompt := &survey.MultiSelect{
		Message:  "Select which directories to keep in sync",
		Options:  dirs,
		PageSize: 15,
	}
	err = config.term.AskOne(prompt, &selectedDirs)
	if err != nil {
		return err
	}

	ignoredDirs := determineLeftOuterUnion(dirs, selectedDirs)
	for i := 0; i < len(ignoredDirs); i++ {
		ignoredDirs[i] = filepath.Join(config.Target, ignoredDirs[i])
	}
	config.term.Infof("Starting, scanning base directories and %d directories", len(dirs)-len(ignoredDirs))
	config.term.Infof("Ignoring %+v", ignoredDirs)

	mismatches, err := config.syncer.GatherMissingSymlinks(ctx, ignoredDirs, config.Source, config.Target)
	if err != nil {
		return err
	}
	config.term.Infof("%+v\n", mismatches)

	for _, f := range mismatches {
		//if the mismatch is "missing from target", just symlink it
		switch f.Issue {
		case MissingFromTarget:
			err = config.syncer.CreateSymlink(f.From, f.To, "backup")
			if err != nil {
				return xerrors.Errorf("error creating symlink: %v", err)
			}
		case MissingFromSource:
			//copy to source and symlink, which is done first since "backup" is defined
			err = config.syncer.CreateSymlink(f.From, f.To, "backup")
			if err != nil {
				return xerrors.Errorf("error creating symlink: %v", err)
			}
		case FileCollision:
			//decide what to do?
		}
	}

	return nil
}

type Syncer interface {
	GatherDirs(ctx context.Context, target string) ([]string, error)
	GatherMissingSymlinks(ctx context.Context, ignoredDirs []string, source, target string) ([]Mismatch, error)
	CreateSymlink(from, to, backup string) error
}

func NewSyncer(fs io.Filesystem) *syncer {
	return &syncer{fs: fs}
}

type syncer struct {
	fs io.Filesystem
}

func (s syncer) CreateSymlink(from, to, backup string) error {
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
	err = os.Symlink(from, to)
	if err != nil {
		return xerrors.Errorf("error symlinking file: %v", err)
	}
	return nil
}

func (s syncer) GatherDirs(ctx context.Context, target string) ([]string, error) {
	var dirs []string
	fileInfo, err := ioutil.ReadDir(target)
	if err != nil {
		return nil, err
	}
	for _, f := range fileInfo {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	return dirs, nil
}

// GatherMissingSymlinks looks and for all the files missing, and creates a collection of mismatched files
func (s syncer) GatherMissingSymlinks(ctx context.Context, ignores []string, source, target string) ([]Mismatch, error) {
	issues := make([]Mismatch, 0)
	w := walker{
		fs:         s.fs,
		baseSource: source,
		baseTarget: target,
		issues:     issues,
		ignores:    ignores,
		log:        io.NewLogger(), //TODO (@morgan): this logger should be injected
	}

	err := godirwalk.Walk(source, &godirwalk.Options{
		Callback: w.GoWalkerSourceToTarget,
		ErrorCallback: func(s string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})

	if err != nil {
		return nil, err
	}

	w.baseSource = target
	w.baseTarget = source
	err = godirwalk.Walk(target, &godirwalk.Options{
		Callback: w.GoWalkerTargetToSource},
	)
	if err != nil {
		return nil, err
	}
	return w.issues, nil
}

// Prompts the user on what course of action to perform on a file conflict:
// 1. Rename target and move to a backup, symlink source to target
// 2. Ignore file and ignore symlink from now on.
// 3. Append file contents, move file, and symlink?
func (s syncer) ResolveFileConflict(ctx context.Context, from, to string) error {
	return nil
}
