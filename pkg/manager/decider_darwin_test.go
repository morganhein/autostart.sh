//go:build darwin
// +build darwin

package manager

//These tests are tightly coupled to my mac environment. It might be good to make a docker container as a test
//environment and test there
func TestDeciderMac(t *testing.T) {
	d := NewDecider(io.NewShellRunner())

	t.Run("empty should run_if and skip_if runs", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), nil, nil)
		assert.True(t, s)
	})

	t.Run("passing run_if runs", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), nil, []string{"which brew"})
		assert.True(t, s)
	})

	t.Run("a failing skip_if prohibits running", func(t *testing.T) {
		s := d.ShouldRun(context.Background(), []string{"which brew"}, nil)
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
		s := d.ShouldRun(context.Background(), []string{"which apk"}, []string{"which brew"})
		assert.True(t, s)
	})
}
