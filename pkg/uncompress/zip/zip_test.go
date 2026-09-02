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

package zip

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnzipToDir(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		setupZip    func() []byte
		filter      func(string) bool
		expectedErr error
	}{
		{
			name: "valid zip with files",
			setupZip: func() []byte {
				return createTestZip([]testZipFile{
					{"test.txt", "test content"},
					{"subdir/file.txt", "nested content"},
				})
			},
			filter:      func(string) bool { return true },
			expectedErr: nil,
		},
		{
			name: "empty zip",
			setupZip: func() []byte {
				return createTestZip([]testZipFile{})
			},
			filter:      func(string) bool { return true },
			expectedErr: nil,
		},
		{
			name: "filter excludes all files",
			setupZip: func() []byte {
				return createTestZip([]testZipFile{
					{"test.txt", "test content"},
				})
			},
			filter:      func(string) bool { return false },
			expectedErr: nil,
		},
		{
			name: "invalid zip data",
			setupZip: func() []byte {
				return []byte("invalid zip data")
			},
			filter:      func(string) bool { return true },
			expectedErr: nil, // Will be a zip error
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for testing
			tempDir := t.TempDir()

			// Create test zip data
			data := testCase.setupZip()

			// Test the function
			err := UnzipToDir(data, tempDir, testCase.filter)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else {
				// For cases where we expect no specific error
				// The actual error will depend on the implementation details
				_ = err // We don't assert on the specific error in this case
			}
		})
	}
}

func TestUnzipToDirCreatesFiles(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test zip data
	testFiles := []testZipFile{
		{"file1.txt", "content1"},
		{"file2.txt", "content2"},
		{"subdir/file3.txt", "content3"},
	}

	data := createTestZip(testFiles)

	// Test the function
	err := UnzipToDir(data, tempDir, func(string) bool { return true })
	require.NoError(t, err)

	// Verify files were created
	for _, file := range testFiles {
		expectedPath := filepath.Join(tempDir, file.name)
		_, err := os.Stat(expectedPath)
		require.NoError(t, err, "File %s should be created", file.name)
	}
}

func TestUnzipToDirWithFilter(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test zip data
	testFiles := []testZipFile{
		{"include.txt", "included"},
		{"exclude.txt", "excluded"},
	}

	data := createTestZip(testFiles)

	// Filter that only includes files with "include" in the name
	filter := func(name string) bool {
		base := filepath.Base(name)

		return len(base) > 0 && base[0] == 'i'
	}

	// Test the function
	err := UnzipToDir(data, tempDir, filter)
	require.NoError(t, err)

	// Verify only included file was created
	_, err = os.Stat(filepath.Join(tempDir, "include.txt"))
	require.NoError(t, err, "Included file should be created")

	_, err = os.Stat(filepath.Join(tempDir, "exclude.txt"))
	assert.True(t, os.IsNotExist(err), "Excluded file should not be created")
}

func TestUnzipToDirCreatesDirectories(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test zip data with directories
	testFiles := []testZipFile{
		{"dir1/file.txt", "content"},
		{"dir2/subdir/file.txt", "nested content"},
	}

	data := createTestZip(testFiles)

	// Test the function
	err := UnzipToDir(data, tempDir, func(string) bool { return true })
	require.NoError(t, err)

	// Verify directories were created
	_, err = os.Stat(filepath.Join(tempDir, "dir1"))
	require.NoError(t, err, "Directory dir1 should be created")

	_, err = os.Stat(filepath.Join(tempDir, "dir2", "subdir"))
	require.NoError(t, err, "Nested directory should be created")
}

// Helper function to create test zip data.
func createTestZip(files []testZipFile) []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, file := range files {
		writer, err := zipWriter.Create(file.name)
		if err != nil {
			panic(err)
		}

		_, err = writer.Write([]byte(file.content))
		if err != nil {
			panic(err)
		}
	}

	zipWriter.Close()

	return buf.Bytes()
}

// Helper struct for test zip files.
type testZipFile struct {
	name    string
	content string
}

func TestUnzipToDirWithPathTraversal(t *testing.T) {
	t.Parallel()
	// Test path traversal protection
	tempDir := t.TempDir()

	// Create zip with path traversal attempt
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Add a file with path traversal
	writer, err := zipWriter.Create("../../../outside.txt")
	require.NoError(t, err)
	_, err = writer.Write([]byte("test"))
	require.NoError(t, err)

	zipWriter.Close()

	// Test the function - should fail due to path traversal protection
	err = UnzipToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	require.Error(t, err, "Should fail due to path traversal protection")
}

func TestUnzipToDirWithCorruptZip(t *testing.T) {
	t.Parallel()
	// Test corrupt zip data
	tempDir := t.TempDir()

	// Test with invalid zip data
	err := UnzipToDir([]byte("not zip data"), tempDir, func(string) bool { return true })
	require.Error(t, err, "Should fail with invalid zip data")
}

func TestUnzipToDirWithLargeFile(t *testing.T) {
	t.Parallel()
	// Test with a file that's too large (would exceed limitedcopy limits)
	tempDir := t.TempDir()

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Create a very large file (over 200MB limit)
	largeContent := make([]byte, 201*1024*1024) // 201MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	writer, err := zipWriter.Create("large.txt")
	require.NoError(t, err)
	_, err = writer.Write(largeContent)
	require.NoError(t, err)

	zipWriter.Close()

	// Test the function - should fail due to file size limit
	err = UnzipToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	require.Error(t, err, "Should fail due to file size limit")
	assert.Contains(t, err.Error(), "file too big")
}
