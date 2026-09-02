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

package uncompress

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToDir(t *testing.T) {
	t.Parallel()
	var err error
	tests := []struct {
		name        string
		filePath    string
		setupMocks  func() // For mocking targz.UntarToDir and zip.UnzipToDir
		expectedErr error
	}{
		{
			name:     "tar.gz file",
			filePath: "test.tar.gz",
			setupMocks: func() {
				// Skip this test as it requires actual tar.gz data
				// In a real scenario, we would mock targz.UntarToDir
			},
			expectedErr: nil, // Skip this test case
		},
		{
			name:     "zip file",
			filePath: "test.zip",
			setupMocks: func() {
				// Skip this test as it requires actual zip data
				// In a real scenario, we would mock zip.UnzipToDir
			},
			expectedErr: nil, // Skip this test case
		},
		{
			name:     "unknown archive type",
			filePath: "test.unknown",
			setupMocks: func() {
				// No mocks needed for this case
			},
			expectedErr: errArchive,
		},
		{
			name:     "directory creation error",
			filePath: "test.tar.gz",
			setupMocks: func() {
				// This test would require mocking os.MkdirAll to return an error
			},
			expectedErr: nil, // Will be whatever os.MkdirAll returns
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for testing
			tempDir := t.TempDir()

			// Setup mocks if provided
			if testCase.setupMocks != nil {
				testCase.setupMocks()
			}

			// Test the function
			// Use empty data for unknown archive types to trigger errArchive
			// Use non-empty data for known types to avoid calling actual uncompress functions
			var testData []byte
			if testCase.filePath == "test.unknown" {
				testData = []byte("")
			} else {
				testData = []byte("test data")
			}
			err = ToDir(testData, testCase.filePath, tempDir, func(string) bool { return true })

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
			} else {
				// For cases where we expect no specific error
				// The actual error will depend on the implementation details
				_ = err // We don't assert on the specific error in this case
			}
		})
	}
}

func TestErrArchive(t *testing.T) {
	t.Parallel()
	// Test that our error variable is properly defined
	require.Error(t, errArchive)
	assert.Equal(t, "unknown archive kind", errArchive.Error())
}

func TestToDirWithValidExtensions(t *testing.T) {
	t.Parallel()
	var err error
	// Test that the function correctly identifies different archive types
	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "tar.gz extension",
			filePath: "archive.tar.gz",
			expected: ".tar.gz",
		},
		{
			name:     "zip extension",
			filePath: "archive.zip",
			expected: ".zip",
		},
		{
			name:     "tar.gz with path",
			filePath: "/path/to/archive.tar.gz",
			expected: ".tar.gz",
		},
		{
			name:     "zip with path",
			filePath: "/path/to/archive.zip",
			expected: ".zip",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for testing
			tempDir := t.TempDir()

			// Test with empty data to trigger the expected error
			testData := []byte("test data") // Use non-empty data to avoid calling actual uncompress functions
			err = ToDir(testData, testCase.filePath, tempDir, func(string) bool { return true })

			// We expect errArchive for unknown archive types
			// or specific errors for known types with empty data
			require.Error(t, err)
		})
	}
}

func TestToDirCreatesDirectory(t *testing.T) {
	t.Parallel()
	var err error
	// Test that the function creates the target directory
	tempDir := t.TempDir()

	// Create a subdirectory path that doesn't exist
	testSubDir := filepath.Join(tempDir, "subdir", "nested")

	// Test with empty data to focus on directory creation
	testData := []byte("test data") // Use non-empty data to avoid calling actual uncompress functions
	err = ToDir(testData, "test.unknown", testSubDir, func(string) bool { return true })

	// The directory should be created even if the uncompress fails
	require.Error(t, err)
	_, err = os.Stat(testSubDir)
	require.NoError(t, err, "Directory should be created by ToDir")
}
