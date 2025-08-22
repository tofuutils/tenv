package versionmanager

import (
	"testing"

	"github.com/tofuutils/tenv/v4/config/envname"
)

func TestEnvPrefix_Version(t *testing.T) {
	tests := []struct {
		name     string
		prefix   EnvPrefix
		expected string
	}{
		{
			name:     "empty prefix",
			prefix:   "",
			expected: envname.VersionSuffix,
		},
		{
			name:     "non-empty prefix",
			prefix:   "TF_",
			expected: "TF_" + envname.VersionSuffix,
		},
		{
			name:     "underscore prefix",
			prefix:   "TF_",
			expected: "TF_" + envname.VersionSuffix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.prefix.Version()
			if result != tt.expected {
				t.Errorf("Version() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEnvPrefix_constraint(t *testing.T) {
	tests := []struct {
		name     string
		prefix   EnvPrefix
		expected string
	}{
		{
			name:     "empty prefix",
			prefix:   "",
			expected: envname.DefaultConstraintSuffix,
		},
		{
			name:     "non-empty prefix",
			prefix:   "TF_",
			expected: "TF_" + envname.DefaultConstraintSuffix,
		},
		{
			name:     "underscore prefix",
			prefix:   "TF_",
			expected: "TF_" + envname.DefaultConstraintSuffix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.prefix.constraint()
			if result != tt.expected {
				t.Errorf("constraint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEnvPrefix_defaultVersion(t *testing.T) {
	tests := []struct {
		name     string
		prefix   EnvPrefix
		expected string
	}{
		{
			name:     "empty prefix",
			prefix:   "",
			expected: envname.DefaultVersionSuffix,
		},
		{
			name:     "non-empty prefix",
			prefix:   "TF_",
			expected: "TF_" + envname.DefaultVersionSuffix,
		},
		{
			name:     "underscore prefix",
			prefix:   "TF_",
			expected: "TF_" + envname.DefaultVersionSuffix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.prefix.defaultVersion()
			if result != tt.expected {
				t.Errorf("defaultVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}
