package pkg

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/karrick/godirwalk"
	"github.com/morganhein/autostart.sh/pkg/T"
	"github.com/morganhein/autostart.sh/pkg/io"
)

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
}

func Sync(config SyncConfig) error {
	syncr := NewSyncer(io.NewFilesystem())
	log := io.NewLogger()
	ctx := context.Background()

	//if we don't already have an ignores (or a config isn't set, or something), then do a first pass to skip certain folders
	dirs, err := syncr.GatherDirs(ctx, config.Target)
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
	err = survey.AskOne(prompt, &selectedDirs)
	if err != nil {
		return T.Log(err)
	}

	ignoredDirs := determineLeftOuterUnion(dirs, selectedDirs)
	for i := 0; i < len(ignoredDirs); i++ {
		ignoredDirs[i] = filepath.Join(config.Target, ignoredDirs[i])
	}
	log.Infof("Starting, scanning base directories and %d directories", len(dirs)-len(ignoredDirs))
	log.Infof("Ignoring %+v", ignoredDirs)

	mismatches, err := syncr.GatherMissingSymlinks(ctx, ignoredDirs, config.Source, config.Target)
	if err != nil {
		return T.Log(err)
	}
	log.Infof("%+v\n", mismatches)
	return nil
}

type Syncer interface {
	GatherMissingSymlinks(ctx context.Context, source, target string) ([]Mismatch, error)
	// ConfigureIgnores scans a target and determines what will be ignored, only scanning one directory deep
	ConfigureIgnores(ctx context.Context, target string) ([]string, error)
}

func NewSyncer(fs io.Filesystem) *syncer {
	return &syncer{fs: fs}
}

type syncer struct {
	fs io.Filesystem
}

func (s syncer) GatherDirs(ctx context.Context, target string) ([]string, error) {
	var dirs []string
	fileInfo, err := ioutil.ReadDir(target)
	if err != nil {
		return nil, T.Log(err)
	}
	for _, f := range fileInfo {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	return dirs, nil
}

//GatherMissingSymlinks looks and for all the files missing, and creates a collection of mismatched files
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
			_ = T.ErrorF("unable to scan specified folder '%v', skipping", s)
			return godirwalk.SkipNode
		},
	})

	if err != nil {
		return nil, T.Log(err)
	}

	w.baseSource = target
	w.baseTarget = source
	err = godirwalk.Walk(target, &godirwalk.Options{
		Callback: w.GoWalkerTargetToSource},
	)
	if err != nil {
		return nil, T.Log(err)
	}
	return w.issues, nil
}
