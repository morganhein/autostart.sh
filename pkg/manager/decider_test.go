package manager

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/morganhein/envy/pkg/io"
	"github.com/stretchr/testify/assert"
)

func TestDeciderMock(t *testing.T) {
	r := &io.RunnerMock{}
	d := NewDecider(r)

	t.Run("empty should run_if and skip_if runs", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), nil, nil)
		assert.True(t, s)
	})

	t.Run("passing run_if runs", func(t *testing.T) {
		r.RunFunc = func(ctx context.Context, printOnly bool, cmdLine string) (string, error) {
			assert.Equal(t, "which brew", cmdLine)
			return "/usr/local/brew", nil
		}
		s := d.ShouldRun(context.Background(), nil, []string{"which brew"})
		assert.True(t, s)
	})

	t.Run("a passing skip_if prohibits running", func(t *testing.T) {
		r.RunFunc = func(ctx context.Context, printOnly bool, cmdLine string) (string, error) {
			assert.Equal(t, "which brew", cmdLine)
			return "/usr/local/brew", nil
		}
		s := d.ShouldRun(context.Background(), []string{"which brew"}, nil)
		assert.False(t, s)
	})

	t.Run("a failing skip_if runs", func(t *testing.T) {
		r.RunFunc = func(ctx context.Context, printOnly bool, cmdLine string) (string, error) {
			assert.Equal(t, "which apk", cmdLine)
			return "", errors.New("command exited with a non-zero exit code")
		}
		s := d.ShouldRun(context.Background(), []string{"which apk"}, nil)
		assert.True(t, s)
	})

	t.Run("passing run_if and failing skip_if prohibits running", func(t *testing.T) {
		r.RunFunc = func(ctx context.Context, printOnly bool, cmdLine string) (string, error) {
			if strings.Contains(cmdLine, "brew") {
				return "/usr/local/brew", errors.New("command exited with a non-zero exit code")
			}
			return "/usr/local/gvm", nil
		}
		s := d.ShouldRun(context.Background(), []string{"which gvm"}, []string{"which brew"})
		assert.False(t, s)
	})

	t.Run("passing run_if and passing skip_if runs", func(t *testing.T) {
		r.RunFunc = func(ctx context.Context, printOnly bool, cmdLine string) (string, error) {
			if strings.Contains(cmdLine, "brew") {
				return "/usr/local/brew", nil
			}
			return "/usr/local/apt", errors.New("command exited with a non-zero exit code")
		}
		s := d.ShouldRun(context.Background(), []string{"which apk"}, []string{"which brew"})
		assert.True(t, s)
	})
}
