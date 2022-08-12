//+build sync

package sync

import (
	"testing"

	"github.com/karrick/godirwalk"
	"github.com/morganhein/envy/pkg/io"
	"github.com/stretchr/testify/assert"
)

func TestGoWalkerSourceToTarget(t *testing.T) {
	w := &walker{
		fs:         io.NewFilesystem(),
		baseSource: "../tests/sync_tests_folders/source",
		baseTarget: "../tests/sync_tests_folders/target",
		issues:     []Mismatch{},
		ignores:    []string{".3T", ".Trash", ".azure", ".cache", ".cups", ".dlv", ".docker", ".eclipse", ".gvm", ".iterm2", ".kube", ".local", ".matrix", ".node-gyp", ".npm", ".oh-my-zsh", ".pgadmin", ".ssh", ".tabnine", ".tldrc", ".vnc", ".vscode", "Applications", "Desktop", "Documents", "Downloads", "Library", "Movies", "Music", "OneDrive", "Pictures", "Projects", "Public", "athens-storage", "dump", "go", "tmp"},
		log:        io.NewLogger(),
	}

	err := godirwalk.Walk("../tests/sync_tests_folders/source", &godirwalk.Options{
		Callback: w.GoWalkerSourceToTarget},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, w.issues)
	assert.Len(t, w.issues, 2)
	assert.Equal(t, w.issues[0].Issue, FileCollision)
	assert.Equal(t, w.issues[1].Issue, MissingFromTarget)
}

func TestGoWalkerTargetToSource(t *testing.T) {
	w := &walker{
		fs:         io.NewFilesystem(),
		baseSource: "../tests/sync_tests_folders/source",
		baseTarget: "../tests/sync_tests_folders/target",
		issues:     []Mismatch{},
	}

	err := godirwalk.Walk("../tests/sync_tests_folders/target", &godirwalk.Options{
		Callback: w.GoWalkerTargetToSource},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, w.issues)
	assert.Len(t, w.issues, 2)
	assert.Equal(t, w.issues[0].Issue, FileCollision)
	assert.Equal(t, w.issues[1].Issue, MissingFromSource)
}

func TestOnDisk_TEMPORARY(t *testing.T) {
	w := &walker{
		fs:         io.NewFilesystem(),
		baseSource: "/Users/morgan/.matrix",
		baseTarget: "/Users/morgan/",
		issues:     []Mismatch{},
		ignores:    []string{"/Users/morgan/.3T", "/Users/morgan/.Trash", "/Users/morgan/.azure", "/Users/morgan/.cache", "/Users/morgan/.cups", "/Users/morgan/.dlv", "/Users/morgan/.docker", "/Users/morgan/.eclipse", "/Users/morgan/.gvm", "/Users/morgan/.iterm2", "/Users/morgan/.kube", "/Users/morgan/.local", "/Users/morgan/.matrix", "/Users/morgan/.node-gyp", "/Users/morgan/.npm", "/Users/morgan/.oh-my-zsh", "/Users/morgan/.pgadmin", "/Users/morgan/.ssh", "/Users/morgan/.tabnine", "/Users/morgan/.tldrc", "/Users/morgan/.vnc", "/Users/morgan/.vscode", "/Users/morgan/Applications", "/Users/morgan/Desktop", "/Users/morgan/Documents", "/Users/morgan/Downloads", "/Users/morgan/Library", "/Users/morgan/Movies", "/Users/morgan/Music", "/Users/morgan/OneDrive", "/Users/morgan/Pictures", "/Users/morgan/Projects", "/Users/morgan/Public", "/Users/morgan/athens-storage", "/Users/morgan/dump", "/Users/morgan/go", "/Users/morgan/tmp"},
		log:        io.NewLogger(),
		linkDirs:   true,
	}

	err := godirwalk.Walk("/Users/morgan", &godirwalk.Options{
		Callback: w.GoWalkerSourceToTarget},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, w.issues)
	for _, issue := range w.issues {
		t.Logf("issue discovered: %v: %v", issue.Issue, issue.To)
	}
}

func TestOnDisk_TEMPORARY2(t *testing.T) {
	w := &walker{
		fs:         io.NewFilesystem(),
		baseSource: "/Users/morgan/.matrix",
		baseTarget: "/Users/morgan/",
		issues:     []Mismatch{},
		ignores:    []string{"/Users/morgan/.3T", "/Users/morgan/.Trash", "/Users/morgan/.azure", "/Users/morgan/.cache", "/Users/morgan/.cups", "/Users/morgan/.dlv", "/Users/morgan/.docker", "/Users/morgan/.eclipse", "/Users/morgan/.gvm", "/Users/morgan/.iterm2", "/Users/morgan/.kube", "/Users/morgan/.local", "/Users/morgan/.matrix", "/Users/morgan/.node-gyp", "/Users/morgan/.npm", "/Users/morgan/.oh-my-zsh", "/Users/morgan/.pgadmin", "/Users/morgan/.ssh", "/Users/morgan/.tabnine", "/Users/morgan/.tldrc", "/Users/morgan/.vnc", "/Users/morgan/.vscode", "/Users/morgan/Applications", "/Users/morgan/Desktop", "/Users/morgan/Documents", "/Users/morgan/Downloads", "/Users/morgan/Library", "/Users/morgan/Movies", "/Users/morgan/Music", "/Users/morgan/OneDrive", "/Users/morgan/Pictures", "/Users/morgan/Projects", "/Users/morgan/Public", "/Users/morgan/athens-storage", "/Users/morgan/dump", "/Users/morgan/go", "/Users/morgan/tmp"},
		log:        io.NewLogger(),
		linkDirs:   true,
	}

	err := godirwalk.Walk("/Users/morgan/.matrix", &godirwalk.Options{
		Callback: w.GoWalkerTargetToSource},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, w.issues)
	for _, issue := range w.issues {
		t.Logf("issue discovered: %v: %v", issue.Issue, issue.To)
	}
}
