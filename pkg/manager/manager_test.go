package manager

 /*
func TestInstallPackageWithoutOverride(t *testing.T) {
	runner := io.RunnerMock{
		RunFunc: func(ctx context.Context, printOnly bool, cmdLine string) (string, error) {
			if strings.Contains(cmdLine, "which apt") {
				return "/bin/apt", nil
			}
			if strings.Contains(cmdLine, "which") {
				return "", errors.New("not found")
			}
			assert.Equal(t, "sudo apt install -y vim", cmdLine)
			return cmdLine, nil
		},
	}
	m := manager{
		d:  NewDecider(&runner),
		r:  &runner,
		dl: nil,
		s:  nil,
	}
	config, err := insureDefaults(TOMLConfig{
		RunningConfig: RunConfig{
			ConfigLocation: "/tmp/any/location",
			DryRun:         true,
		},
	})
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = m.installPkgHelper(ctx, config, envVariables{}, "vim")
	assert.NoError(t, err)
}

func TestInstallPackageWithOverride(t *testing.T) {
	runner := io.RunnerMock{
		RunFunc: func(ctx context.Context, printOnly bool, cmdLine string) (string, error) {
			if strings.Contains(cmdLine, "which apt") {
				return "/bin/apt", nil
			}
			if strings.Contains(cmdLine, "which") {
				return "", errors.New("not found")
			}
			assert.Equal(t, "sudo apt install -y neovim", cmdLine)
			return cmdLine, nil
		},
	}
	m := manager{
		d:  NewDecider(&runner),
		r:  &runner,
		dl: nil,
		s:  nil,
	}
	config, err := insureDefaults(TOMLConfig{
		Packages: map[string]Package{
			"vim": {
				"apt": "neovim",
			},
		},
		RunningConfig: RunConfig{
			ConfigLocation: "/tmp/any/location",
			DryRun:         true,
		},
	})
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = m.installPkgHelper(ctx, config, envVariables{}, "vim")
	assert.NoError(t, err)
}
*/