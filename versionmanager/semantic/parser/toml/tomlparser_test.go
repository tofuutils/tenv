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
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/tofuutils/tenv/v4/config"
)

// mockDisplayer implements loghelper.Displayer for testing
type mockDisplayer struct{}

func (m *mockDisplayer) Display(string)                          {}
func (m *mockDisplayer) Log(hclog.Level, string, ...interface{}) {}
func (m *mockDisplayer) IsDebug() bool                           { return false }
func (m *mockDisplayer) IsTrace() bool                           { return false }
func (m *mockDisplayer) Flush(bool)                              {}

func TestRetrieveVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		content        string
		expectedResult string
		expectError    bool
	}{
		{
			name: "valid version constraint",
			content: `[tool.terraform]
			required_version = "~> 1.0.0"`,
			expectedResult: "~> 1.0.0",
		},
		{
			name: "no version constraint",
			content: `[tool.terraform]
			backend = "local"`,
			expectedResult: "",
		},
		{
			name: "invalid TOML",
			content: `[tool.terraform
			required_version = "~> 1.0.0"`,
			expectError: true,
		},
		{
			name: "non-string version constraint",
			content: `[tool.terraform]
			required_version = 1.0`,
			expectError: true,
		},
		{
			name: "empty version constraint",
			content: `[tool.terraform]
			required_version = ""`,
			expectedResult: "",
		},
		{
			name: "complex TOML with version constraint",
			content: `[tool.terraform]
			required_version = "~> 1.0.0"
			[tool.terraform.required_providers]
			aws = { source = "hashicorp/aws", version = "~> 3.0" }`,
			expectedResult: "~> 1.0.0",
		},
		{
			name: "multiple version constraints",
			content: `[tool.terraform]
			required_version = ">= 1.0.0, < 2.0.0"`,
			expectedResult: ">= 1.0.0, < 2.0.0",
		},
		{
			name: "nested tables",
			content: `[tool]
			[tool.terraform]
			required_version = "~> 1.0.0"`,
			expectedResult: "~> 1.0.0",
		},
		{
			name: "array of tables",
			content: `[[tool.terraform]]
			required_version = "~> 1.0.0"`,
			expectedResult: "~> 1.0.0",
		},
		{
			name: "inline tables",
			content: `[tool.terraform]
			required_version = "~> 1.0.0"
			backend = { type = "local" }`,
			expectedResult: "~> 1.0.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "pyproject.toml")

			// Create test file
			err := os.WriteFile(filePath, []byte(tt.content), 0600)
			if err != nil {
				t.Fatal(err)
			}

			// Create config
			conf := &config.Config{
				Displayer: &mockDisplayer{},
			}

			// Run test
			result, err := RetrieveVersion(filePath, conf)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if result != tt.expectedResult {
				t.Errorf("expected %s but got %s", tt.expectedResult, result)
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	// Setup temp directory
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "pyproject.toml")

	// Create test file
	content := `[tool.terraform]
	required_version = "~> 1.0.0"`
	err := os.WriteFile(filePath, []byte(content), 0600)
	if err != nil {
		t.Fatal(err)
	}

	// Create config
	conf := &config.Config{
		Displayer: &mockDisplayer{},
	}

	// Number of concurrent goroutines
	numGoroutines := 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Run concurrent tests
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			result, err := RetrieveVersion(filePath, conf)
			if err != nil {
				t.Error(err)
				return
			}
			if result != "~> 1.0.0" {
				t.Errorf("expected ~> 1.0.0 but got %s", result)
			}
		}()
	}

	wg.Wait()
}

func TestFileErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(string) error
		expectError bool
	}{
		{
			name: "non-existent file",
			setup: func(dir string) error {
				return nil // No setup needed, file doesn't exist
			},
			expectError: false, // Should return empty string, not error
		},
		{
			name: "unreadable file",
			setup: func(dir string) error {
				filePath := filepath.Join(dir, "pyproject.toml")
				if err := os.WriteFile(filePath, []byte(`[tool.terraform]`), 0600); err != nil {
					return err
				}
				return os.Chmod(filePath, 0000)
			},
			expectError: true,
		},
		{
			name: "directory instead of file",
			setup: func(dir string) error {
				filePath := filepath.Join(dir, "pyproject.toml")
				return os.Mkdir(filePath, 0700)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()

			// Apply setup
			if err := tt.setup(tempDir); err != nil {
				t.Fatal(err)
			}

			// Create config
			conf := &config.Config{
				Displayer: &mockDisplayer{},
			}

			// Run test
			filePath := filepath.Join(tempDir, "pyproject.toml")
			_, err := RetrieveVersion(filePath, conf)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestParserErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "empty content",
			content:     "",
			expectError: false,
		},
		{
			name: "invalid TOML syntax",
			content: `[tool.terraform
			required_version = "~> 1.0.0"`,
			expectError: true,
		},
		{
			name: "invalid version format",
			content: `[tool.terraform]
			required_version = 1.0`,
			expectError: true,
		},
		{
			name: "missing tool section",
			content: `[terraform]
			required_version = "~> 1.0.0"`,
			expectError: false,
		},
		{
			name: "missing terraform section",
			content: `[tool]
			required_version = "~> 1.0.0"`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "pyproject.toml")

			// Create test file
			err := os.WriteFile(filePath, []byte(tt.content), 0600)
			if err != nil {
				t.Fatal(err)
			}

			// Create config
			conf := &config.Config{
				Displayer: &mockDisplayer{},
			}

			// Run test
			_, err = RetrieveVersion(filePath, conf)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestFileEncodings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		content        []byte
		expectedResult string
		expectError    bool
	}{
		{
			name:           "UTF-8",
			content:        []byte(`[tool.terraform]\nrequired_version = "~> 1.0.0"`),
			expectedResult: "~> 1.0.0",
		},
		{
			name:           "UTF-8 with BOM",
			content:        append([]byte{0xEF, 0xBB, 0xBF}, []byte(`[tool.terraform]\nrequired_version = "~> 1.0.0"`)...),
			expectedResult: "~> 1.0.0",
		},
		{
			name:        "UTF-16",
			content:     append([]byte{0xFF, 0xFE}, []byte(`[tool.terraform]\nrequired_version = "~> 1.0.0"`)...),
			expectError: true,
		},
		{
			name:           "ASCII",
			content:        []byte(`[tool.terraform]\nrequired_version = "~> 1.0.0"`),
			expectedResult: "~> 1.0.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "pyproject.toml")

			// Create test file
			err := os.WriteFile(filePath, tt.content, 0600)
			if err != nil {
				t.Fatal(err)
			}

			// Create config
			conf := &config.Config{
				Displayer: &mockDisplayer{},
			}

			// Run test
			result, err := RetrieveVersion(filePath, conf)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if result != tt.expectedResult {
				t.Errorf("expected %s but got %s", tt.expectedResult, result)
			}
		})
	}
}

func TestLargeFiles(t *testing.T) {
	t.Parallel()

	// Setup temp directory
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "pyproject.toml")

	// Create a large file with version constraint
	content := make([]byte, 10*1024*1024) // 10MB
	copy(content, []byte(`[tool.terraform]\nrequired_version = "~> 1.0.0"`))

	// Create test file
	err := os.WriteFile(filePath, content, 0600)
	if err != nil {
		t.Fatal(err)
	}

	// Create config
	conf := &config.Config{
		Displayer: &mockDisplayer{},
	}

	// Run test
	result, err := RetrieveVersion(filePath, conf)
	if err != nil {
		t.Fatal(err)
	}

	if result != "~> 1.0.0" {
		t.Errorf("expected ~> 1.0.0 but got %s", result)
	}
}

func TestSymbolicLinks(t *testing.T) {
	t.Parallel()

	// Setup temp directory
	tempDir := t.TempDir()
	originalPath := filepath.Join(tempDir, "original.toml")
	linkPath := filepath.Join(tempDir, "pyproject.toml")

	// Create original file
	content := `[tool.terraform]\nrequired_version = "~> 1.0.0"`
	err := os.WriteFile(originalPath, []byte(content), 0600)
	if err != nil {
		t.Fatal(err)
	}

	// Create symbolic link
	err = os.Symlink(originalPath, linkPath)
	if err != nil {
		t.Fatal(err)
	}

	// Create config
	conf := &config.Config{
		Displayer: &mockDisplayer{},
	}

	// Run test
	result, err := RetrieveVersion(linkPath, conf)
	if err != nil {
		t.Fatal(err)
	}

	if result != "~> 1.0.0" {
		t.Errorf("expected ~> 1.0.0 but got %s", result)
	}
}

func TestMultipleFiles(t *testing.T) {
	t.Parallel()

	// Setup temp directory
	tempDir := t.TempDir()

	// Create multiple files with different version constraints
	files := []struct {
		name    string
		content string
	}{
		{
			name:    "pyproject.toml",
			content: `[tool.terraform]\nrequired_version = "~> 1.0.0"`,
		},
		{
			name:    "poetry.toml",
			content: `[tool.terraform]\nrequired_version = "~> 1.1.0"`,
		},
		{
			name:    "config.toml",
			content: `[tool.terraform]\nbackend = "local"`,
		},
	}

	// Create config
	conf := &config.Config{
		Displayer: &mockDisplayer{},
	}

	// Create and test each file
	for _, file := range files {
		filePath := filepath.Join(tempDir, file.name)
		err := os.WriteFile(filePath, []byte(file.content), 0600)
		if err != nil {
			t.Fatal(err)
		}

		result, err := RetrieveVersion(filePath, conf)
		if err != nil {
			t.Fatal(err)
		}

		expected := ""
		if file.name != "config.toml" {
			expected = "~> 1.0.0"
			if file.name == "poetry.toml" {
				expected = "~> 1.1.0"
			}
		}

		if result != expected {
			t.Errorf("for file %s, expected %s but got %s", file.name, expected, result)
		}
	}
}
