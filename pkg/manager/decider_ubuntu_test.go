//go:build ubuntu
// +build ubuntu

package manager

import (
	"context"
	"github.com/morganhein/shoelace/pkg/io"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeciderUbuntu(t *testing.T) {
	d := NewDecider(io.NewShellRunner())

	t.Run("empty should run_if and skip_if runs", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), nil, nil)
		assert.True(t, s)
	})

	t.Run("passing run_if runs", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), nil, []string{"which apt"})
		assert.True(t, s)
	})

	t.Run("a failing skip_if prohibits running", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), []string{"which apt"}, nil)
		assert.False(t, s)
	})

	t.Run("passing skip_if runs", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), []string{"which apk"}, nil)
		assert.True(t, s)
	})

	t.Run("passing run_if and failing skip_if prohibits running", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), []string{"which gvm"}, []string{"which brew"})
		assert.False(t, s)
	})

	t.Run("passing run_if and passing skip_if runs", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), []string{"which apk"}, []string{"which apt"})
		assert.True(t, s)
	})
}
