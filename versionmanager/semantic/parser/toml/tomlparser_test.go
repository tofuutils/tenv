/*
 *
 * Copyright 2024 tofuutils authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package tomlparser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofuutils/tenv/v4/config"
)

func TestRetrieveVersion(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		tomlContent string
		expected    string
		expectError bool
		expectEmpty bool
	}{
		{
			name: "valid TOML with version string",
			tomlContent: `
version = "1.5.0"
`,
			expected:    "1.5.0",
			expectError: false,
			expectEmpty: false,
		},
		{
			name: "valid TOML with version as quoted string",
			tomlContent: `
version = "v2.0.0"
`,
			expected:    "v2.0.0",
			expectError: false,
			expectEmpty: false,
		},
		{
			name: "TOML with version as integer",
			tomlContent: `
version = 1
`,
			expected:    "",
			expectError: true,
			expectEmpty: false,
		},
		{
			name: "TOML with version as float",
			tomlContent: `
version = 1.5
`,
			expected:    "",
			expectError: true,
			expectEmpty: false,
		},
		{
			name: "TOML without version field",
			tomlContent: `
name = "test"
description = "A test configuration"
`,
			expected:    "",
			expectError: false,
			expectEmpty: true,
		},
		{
			name: "TOML with empty version",
			tomlContent: `
version = ""
`,
			expected:    "",
			expectError: false,
			expectEmpty: true,
		},
		{
			name: "TOML with version as boolean true",
			tomlContent: `
version = true
`,
			expected:    "",
			expectError: true,
			expectEmpty: false,
		},
		{
			name: "TOML with version as boolean false",
			tomlContent: `
version = false
`,
			expected:    "",
			expectError: true,
			expectEmpty: false,
		},
		{
			name: "TOML with complex version constraint",
			tomlContent: `
version = ">= 1.0.0, < 2.0.0"
`,
			expected:    ">= 1.0.0, < 2.0.0",
			expectError: false,
			expectEmpty: false,
		},
		{
			name: "TOML with semantic version",
			tomlContent: `
version = "1.2.3-alpha.1+build.123"
`,
			expected:    "1.2.3-alpha.1+build.123",
			expectError: false,
			expectEmpty: false,
		},
		{
			name: "TOML with version in nested table",
			tomlContent: `
[terraform]
version = "1.0.0"
`,
			expected:    "",
			expectError: false,
			expectEmpty: true,
		},
		{
			name: "TOML with version in array",
			tomlContent: `
version = ["1.0.0", "1.1.0"]
`,
			expected:    "", // Arrays are not strings, so this should return empty
			expectError: true,
			expectEmpty: false,
		},
		{
			name:        "empty TOML file",
			tomlContent: ``,
			expected:    "",
			expectError: false,
			expectEmpty: true,
		},
		{
			name: "TOML with comments and whitespace",
			tomlContent: `
# This is a comment
version = "1.0.0"  # Version comment
`,
			expected:    "1.0.0",
			expectError: false,
			expectEmpty: false,
		},
		{
			name: "TOML with multiple key-value pairs",
			tomlContent: `
name = "myapp"
version = "2.1.0"
description = "My application"
`,
			expected:    "2.1.0",
			expectError: false,
			expectEmpty: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary TOML file
			fileName := strings.ReplaceAll(t.Name(), "/", "_") + ".toml"
			filePath := filepath.Join(tempDir, fileName)

			err := os.WriteFile(filePath, []byte(testCase.tomlContent), 0o600)
			require.NoError(t, err)

			// Create config with mock displayer
			conf := &config.Config{
				Displayer: &mockDisplayer{},
			}

			// Test the function
			version, err := RetrieveVersion(filePath, conf)

			if testCase.expectError {
				require.Error(t, err)
				assert.Empty(t, version)
			} else {
				require.NoError(t, err)
				if testCase.expectEmpty {
					assert.Empty(t, version)
				} else {
					assert.Equal(t, testCase.expected, version)
				}
			}
		})
	}
}

func TestRetrieveVersionFileNotFound(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	// Use a non-existent file path
	nonExistentPath := filepath.Join(tempDir, "nonexistent.toml")

	conf := &config.Config{
		Displayer: &mockDisplayer{},
	}

	version, err := RetrieveVersion(nonExistentPath, conf)

	require.NoError(t, err)
	assert.Empty(t, version)
}

func TestRetrieveVersionInvalidTOML(t *testing.T) {
	t.Parallel()
	var err error
	tempDir := t.TempDir()

	// Create a file with invalid TOML syntax
	fileName := "invalid.toml"
	filePath := filepath.Join(tempDir, fileName)
	invalidTOML := `
version = "1.0.0"
   invalid syntax here [[
`

	err = os.WriteFile(filePath, []byte(invalidTOML), 0o600)
	require.NoError(t, err)

	conf := &config.Config{
		Displayer: &mockDisplayer{},
	}

	version, err := RetrieveVersion(filePath, conf)

	require.Error(t, err)
	assert.Empty(t, version)
}

func TestRetrieveVersionWithSpecialCharacters(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		tomlContent string
		expected    string
	}{
		{
			name: "version with special characters",
			tomlContent: `
version = "1.0.0-beta+build.123"
`,
			expected: "1.0.0-beta+build.123",
		},
		{
			name: "version with dots and dashes",
			tomlContent: `
version = "v1.2.3-alpha.1"
`,
			expected: "v1.2.3-alpha.1",
		},
		{
			name: "version with underscores",
			tomlContent: `
version = "1_0_0"
`,
			expected: "1_0_0",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			fileName := strings.ReplaceAll(t.Name(), "/", "_") + ".toml"
			filePath := filepath.Join(tempDir, fileName)

			err := os.WriteFile(filePath, []byte(testCase.tomlContent), 0o600)
			require.NoError(t, err)

			conf := &config.Config{
				Displayer: &mockDisplayer{},
			}

			version, err := RetrieveVersion(filePath, conf)

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, version)
		})
	}
}

// mockDisplayer is a test implementation of the displayer interface.
type mockDisplayer struct{}

func (m *mockDisplayer) Display(msg string)                                     {}
func (m *mockDisplayer) Log(level hclog.Level, msg string, args ...interface{}) {}
func (m *mockDisplayer) IsDebug() bool                                          { return false }
func (m *mockDisplayer) Flush(logMode bool)                                     {}
