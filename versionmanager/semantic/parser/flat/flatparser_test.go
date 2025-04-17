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

package flatparser

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

// mockDisplayer implements loghelper.Displayer for testing
type mockDisplayer struct{}

func (m *mockDisplayer) Display(string)                          {}
func (m *mockDisplayer) Log(hclog.Level, string, ...interface{}) {}
func (m *mockDisplayer) IsDebug() bool                           { return false }
func (m *mockDisplayer) IsTrace() bool                           { return false }
func (m *mockDisplayer) Flush(bool)                              {}

func TestRetrieve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		content        string
		displayMsg     func(loghelper.Displayer, string, string) string
		expectedResult string
		expectError    bool
	}{
		{
			name:           "valid version with NoMsg",
			content:        "1.0.0",
			displayMsg:     NoMsg,
			expectedResult: "1.0.0",
		},
		{
			name:           "valid version with DisplayDetectionInfo",
			content:        "1.0.0",
			displayMsg:     types.DisplayDetectionInfo,
			expectedResult: "1.0.0",
		},
		{
			name:           "version with spaces",
			content:        "    1.0.0    ",
			displayMsg:     NoMsg,
			expectedResult: "1.0.0",
		},
		{
			name:           "version with tabs",
			content:        "\t1.0.0\t",
			displayMsg:     NoMsg,
			expectedResult: "1.0.0",
		},
		{
			name:           "empty file",
			content:        "",
			displayMsg:     NoMsg,
			expectedResult: "",
		},
		{
			name:           "whitespace only",
			content:        "    \t    ",
			displayMsg:     NoMsg,
			expectedResult: "",
		},
		{
			name:           "newlines only",
			content:        "\n\n\n",
			displayMsg:     NoMsg,
			expectedResult: "",
		},
		{
			name:           "version with newlines",
			content:        "\n1.0.0\n",
			displayMsg:     NoMsg,
			expectedResult: "1.0.0",
		},
		{
			name:           "version with comments",
			content:        "1.0.0 # comment",
			displayMsg:     NoMsg,
			expectedResult: "1.0.0 # comment",
		},
		{
			name:           "version with multiple lines",
			content:        "1.0.0\n2.0.0",
			displayMsg:     NoMsg,
			expectedResult: "1.0.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "version.txt")

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
			result, err := Retrieve(filePath, conf, tt.displayMsg)

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

func TestRetrieveVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		content        string
		expectedResult string
		expectError    bool
	}{
		{
			name:           "valid version",
			content:        "1.0.0",
			expectedResult: "1.0.0",
		},
		{
			name:           "version with spaces",
			content:        "    1.0.0    ",
			expectedResult: "1.0.0",
		},
		{
			name:           "version with tabs",
			content:        "\t1.0.0\t",
			expectedResult: "1.0.0",
		},
		{
			name:           "empty file",
			content:        "",
			expectedResult: "",
		},
		{
			name:           "whitespace only",
			content:        "    \t    ",
			expectedResult: "",
		},
		{
			name:           "newlines only",
			content:        "\n\n\n",
			expectedResult: "",
		},
		{
			name:           "version with newlines",
			content:        "\n1.0.0\n",
			expectedResult: "1.0.0",
		},
		{
			name:           "version with comments",
			content:        "1.0.0 # comment",
			expectedResult: "1.0.0 # comment",
		},
		{
			name:           "version with multiple lines",
			content:        "1.0.0\n2.0.0",
			expectedResult: "1.0.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "version.txt")

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
	filePath := filepath.Join(tempDir, "version.txt")

	// Create test file
	content := "1.0.0"
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
			// Test both Retrieve and RetrieveVersion
			funcs := []struct {
				name string
				fn   func(string, *config.Config) (string, error)
			}{
				{"RetrieveVersion", RetrieveVersion},
			}

			for _, f := range funcs {
				result, err := f.fn(filePath, conf)
				if err != nil {
					t.Error(err)
					return
				}
				if result != "1.0.0" {
					t.Errorf("for %s, expected 1.0.0 but got %s", f.name, result)
				}
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
				filePath := filepath.Join(dir, "version.txt")
				if err := os.WriteFile(filePath, []byte("1.0.0"), 0600); err != nil {
					return err
				}
				return os.Chmod(filePath, 0000)
			},
			expectError: true,
		},
		{
			name: "directory instead of file",
			setup: func(dir string) error {
				filePath := filepath.Join(dir, "version.txt")
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
			filePath := filepath.Join(tempDir, "version.txt")
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
			content:        []byte("1.0.0"),
			expectedResult: "1.0.0",
		},
		{
			name:           "UTF-8 with BOM",
			content:        append([]byte{0xEF, 0xBB, 0xBF}, []byte("1.0.0")...),
			expectedResult: "1.0.0",
		},
		{
			name:        "UTF-16",
			content:     append([]byte{0xFF, 0xFE}, []byte("1.0.0")...),
			expectError: true,
		},
		{
			name:           "ASCII",
			content:        []byte("1.0.0"),
			expectedResult: "1.0.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "version.txt")

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
	filePath := filepath.Join(tempDir, "version.txt")

	// Create a large file with version constraint
	content := make([]byte, 10*1024*1024) // 10MB
	copy(content, []byte("1.0.0"))

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

	if result != "1.0.0" {
		t.Errorf("expected 1.0.0 but got %s", result)
	}
}

func TestSymbolicLinks(t *testing.T) {
	t.Parallel()

	// Setup temp directory
	tempDir := t.TempDir()
	originalPath := filepath.Join(tempDir, "original.txt")
	linkPath := filepath.Join(tempDir, "version.txt")

	// Create original file
	content := "1.0.0"
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

	if result != "1.0.0" {
		t.Errorf("expected 1.0.0 but got %s", result)
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
			name:    "version.txt",
			content: "1.0.0",
		},
		{
			name:    "other.txt",
			content: "1.1.0",
		},
		{
			name:    "config.txt",
			content: "2.0.0",
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

		if result != file.content {
			t.Errorf("for file %s, expected %s but got %s", file.name, file.content, result)
		}
	}
}
