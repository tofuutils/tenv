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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

func TestNoMsg(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "simple version",
			value:    "1.0.0",
			expected: "1.0.0",
		},
		{
			name:     "version with spaces",
			value:    "  1.0.0  ",
			expected: "  1.0.0  ",
		},
		{
			name:     "empty string",
			value:    "",
			expected: "",
		},
		{
			name:     "version with v prefix",
			value:    "v2.1.0",
			expected: "v2.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NoMsg(nil, tt.value, "")
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRetrieve(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "flatparser_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		fileContent string
		fileName    string
		displayMsg  func(loghelper.Displayer, string, string) string
		expected    string
		expectError bool
	}{
		{
			name:        "valid version file",
			fileContent: "1.0.0",
			fileName:    "version.txt",
			displayMsg:  NoMsg,
			expected:    "1.0.0",
			expectError: false,
		},
		{
			name:        "version file with spaces",
			fileContent: "  1.5.0  ",
			fileName:    "version_with_spaces.txt",
			displayMsg:  NoMsg,
			expected:    "1.5.0",
			expectError: false,
		},
		{
			name:        "empty file",
			fileContent: "",
			fileName:    "empty.txt",
			displayMsg:  NoMsg,
			expected:    "",
			expectError: false,
		},
		{
			name:        "file with only whitespace",
			fileContent: "   \n\t  \n",
			fileName:    "whitespace.txt",
			displayMsg:  NoMsg,
			expected:    "",
			expectError: false,
		},
		{
			name:        "version with v prefix",
			fileContent: "v2.0.0",
			fileName:    "version_with_v.txt",
			displayMsg:  NoMsg,
			expected:    "v2.0.0",
			expectError: false,
		},
		{
			name:        "non-existent file",
			fileContent: "",
			fileName:    "nonexistent.txt",
			displayMsg:  NoMsg,
			expected:    "",
			expectError: false,
		},
		{
			name:        "with display message",
			fileContent: "1.2.3",
			fileName:    "version_display.txt",
			displayMsg:  types.DisplayDetectionInfo,
			expected:    "1.2.3",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the test file if it has content
			if tt.fileContent != "" || tt.name == "empty file" || tt.name == "file with only whitespace" {
				filePath := filepath.Join(tempDir, tt.fileName)
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				require.NoError(t, err)
			}

			// Create a config for testing
			conf, err := config.DefaultConfig()
			require.NoError(t, err)
			conf.InitDisplayer(false)

			// Test the Retrieve function
			filePath := filepath.Join(tempDir, tt.fileName)
			result, err := Retrieve(filePath, &conf, tt.displayMsg)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestRetrieveVersion(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "flatparser_version_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		fileContent string
		fileName    string
		expected    string
		expectError bool
	}{
		{
			name:        "valid version file",
			fileContent: "1.0.0",
			fileName:    "version.txt",
			expected:    "1.0.0",
			expectError: false,
		},
		{
			name:        "version file with spaces",
			fileContent: "  1.5.0  ",
			fileName:    "version_with_spaces.txt",
			expected:    "1.5.0",
			expectError: false,
		},
		{
			name:        "empty file",
			fileContent: "",
			fileName:    "empty.txt",
			expected:    "",
			expectError: false,
		},
		{
			name:        "version with v prefix",
			fileContent: "v2.0.0",
			fileName:    "version_with_v.txt",
			expected:    "v2.0.0",
			expectError: false,
		},
		{
			name:        "non-existent file",
			fileContent: "",
			fileName:    "nonexistent.txt",
			expected:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the test file if it has content
			if tt.fileContent != "" || tt.name == "empty file" {
				filePath := filepath.Join(tempDir, tt.fileName)
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				require.NoError(t, err)
			}

			// Create a config for testing
			conf, err := config.DefaultConfig()
			require.NoError(t, err)
			conf.InitDisplayer(false)

			// Test the RetrieveVersion function
			filePath := filepath.Join(tempDir, tt.fileName)
			result, err := RetrieveVersion(filePath, &conf)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that the functions exist and are callable
	assert.NotNil(t, NoMsg)
	assert.NotNil(t, Retrieve)
	assert.NotNil(t, RetrieveVersion)
}

func TestFileOperations(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "flatparser_fileops_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test with a file that has various edge cases
	testFile := filepath.Join(tempDir, "test_version.txt")
	testContent := "  v1.3.0-alpha  \n"

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Create a config for testing
	conf, err := config.DefaultConfig()
	require.NoError(t, err)
	conf.InitDisplayer(false)

	// Test Retrieve with NoMsg
	result, err := Retrieve(testFile, &conf, NoMsg)
	assert.NoError(t, err)
	assert.Equal(t, "v1.3.0-alpha", result)

	// Test RetrieveVersion (uses types.DisplayDetectionInfo internally)
	result, err = RetrieveVersion(testFile, &conf)
	assert.NoError(t, err)
	assert.Equal(t, "v1.3.0-alpha", result)
}

func TestErrorHandling(t *testing.T) {
	// Create a config for testing
	conf, err := config.DefaultConfig()
	require.NoError(t, err)
	conf.InitDisplayer(false)

	// Test with non-existent file
	result, err := Retrieve("/non/existent/file", &conf, NoMsg)
	assert.NoError(t, err)
	assert.Equal(t, "", result)

	// Test RetrieveVersion with non-existent file
	result, err = RetrieveVersion("/non/existent/file", &conf)
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}
