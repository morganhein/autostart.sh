package manager

import (
	"testing"
)

func TestInstallCommandVariableSubstitution(t *testing.T) {
	config := Config{
		Installers: map[string]Installer{
			"brew": {
				Name:  "brew",
				RunIf: nil,
				Sudo:  true,
				Cmd:   "${sudo} brew install ${pkg}",
			},
			"apk": {
				Name:  "apk",
				RunIf: nil,
				Sudo:  false,
				Cmd:   "${sudo} apk add ${pkg}",
			},
		},
		Tasks: nil,
	}
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "brew",
			expected: "sudo brew install pkg",
		},
		{
			name:     "apk",
			expected: "apk add pkg",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := installCommandVariableSubstitution(config.Installers[test.name].Cmd, "pkg", config.Installers[test.name].Sudo)
			if result != test.expected {
				t.Errorf("expected the result to be `%v`, but received `%v`", test.expected, result)
			}
		})
	}
}
