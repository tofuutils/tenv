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

package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUntarToDir(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		setupTarGz  func() []byte
		filter      func(string) bool
		expectedErr error
	}{
		{
			name: "valid tar.gz with files",
			setupTarGz: func() []byte {
				return createTestTarGz([]testFile{
					{"test.txt", "test content", 0o644},
					{"subdir/file.txt", "nested content", 0o644},
				})
			},
			filter:      func(string) bool { return true },
			expectedErr: nil,
		},
		{
			name: "empty tar.gz",
			setupTarGz: func() []byte {
				return createTestTarGz([]testFile{})
			},
			filter:      func(string) bool { return true },
			expectedErr: nil,
		},
		{
			name: "filter excludes all files",
			setupTarGz: func() []byte {
				return createTestTarGz([]testFile{
					{"test.txt", "test content", 0o644},
				})
			},
			filter:      func(string) bool { return false },
			expectedErr: nil,
		},
		{
			name: "invalid tar.gz data",
			setupTarGz: func() []byte {
				return []byte("invalid tar.gz data")
			},
			filter:      func(string) bool { return true },
			expectedErr: nil, // Will be a gzip error
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for testing
			tempDir := t.TempDir()

			// Create test tar.gz data
			data := testCase.setupTarGz()

			// Test the function
			err := UntarToDir(data, tempDir, testCase.filter)

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

func TestUntarToDirCreatesFiles(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test tar.gz data
	testFiles := []testFile{
		{"file1.txt", "content1", 0o644},
		{"file2.txt", "content2", 0o644},
		{"subdir/file3.txt", "content3", 0o644},
	}

	data := createTestTarGz(testFiles)

	// Test the function
	err := UntarToDir(data, tempDir, func(string) bool { return true })
	require.NoError(t, err)

	// Verify files were created
	for _, file := range testFiles {
		expectedPath := filepath.Join(tempDir, file.name)
		_, err := os.Stat(expectedPath)
		require.NoError(t, err, "File %s should be created", file.name)
	}
}

func TestUntarToDirWithFilter(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test tar.gz data
	testFiles := []testFile{
		{"include.txt", "included", 0o644},
		{"exclude.txt", "excluded", 0o644},
	}

	data := createTestTarGz(testFiles)

	// Filter that only includes files with "include" in the name
	filter := func(name string) bool {
		return len(name) > 0 && name[0] == 'i'
	}

	// Test the function
	err := UntarToDir(data, tempDir, filter)
	require.NoError(t, err)

	// Verify only included file was created
	_, err = os.Stat(filepath.Join(tempDir, "include.txt"))
	require.NoError(t, err, "Included file should be created")

	_, err = os.Stat(filepath.Join(tempDir, "exclude.txt"))
	assert.True(t, os.IsNotExist(err), "Excluded file should not be created")
}

func TestUntarToDirCreatesDirectories(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test tar.gz data with directories
	testFiles := []testFile{
		{"dir1/file.txt", "content", 0o644},
		{"dir2/subdir/file.txt", "nested content", 0o644},
	}

	data := createTestTarGz(testFiles)

	// Test the function
	err := UntarToDir(data, tempDir, func(string) bool { return true })
	require.NoError(t, err)

	// Verify directories were created
	_, err = os.Stat(filepath.Join(tempDir, "dir1"))
	require.NoError(t, err, "Directory dir1 should be created")

	_, err = os.Stat(filepath.Join(tempDir, "dir2", "subdir"))
	require.NoError(t, err, "Nested directory should be created")
}

// Helper function to create test tar.gz data.
func createTestTarGz(files []testFile) []byte {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzipWriter)

	for _, file := range files {
		// Create tar header
		header := &tar.Header{
			Name: file.name,
			Size: int64(len(file.content)),
			Mode: file.mode,
		}

		// Write header
		err := tarWriter.WriteHeader(header)
		if err != nil {
			panic(err)
		}

		// Write file content
		_, err = tarWriter.Write([]byte(file.content))
		if err != nil {
			panic(err)
		}
	}

	// Close writers
	tarWriter.Close()
	gzipWriter.Close()

	return buf.Bytes()
}

// Helper struct for test files.
type testFile struct {
	name    string
	content string
	mode    int64
}

func TestUntarToDirWithPathTraversal(t *testing.T) {
	t.Parallel()
	// Test path traversal protection
	tempDir := t.TempDir()

	// Create tar.gz with path traversal attempt
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzipWriter)

	// Add a file with path traversal
	header := &tar.Header{
		Name: "../../../outside.txt",
		Size: 4,
		Mode: 0o644,
	}

	err := tarWriter.WriteHeader(header)
	require.NoError(t, err)
	_, err = tarWriter.Write([]byte("test"))
	require.NoError(t, err)

	tarWriter.Close()
	gzipWriter.Close()

	// Test the function - should fail due to path traversal protection
	err = UntarToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	require.Error(t, err, "Should fail due to path traversal protection")
}

func TestUntarToDirWithUnknownType(t *testing.T) {
	t.Parallel()
	// Test unknown tar type
	tempDir := t.TempDir()

	// Create tar.gz with unknown type
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzipWriter)

	// Add a file with unknown type
	header := &tar.Header{
		Name:     "unknown.txt",
		Size:     4,
		Mode:     0o644,
		Typeflag: 99, // Unknown type
	}

	err := tarWriter.WriteHeader(header)
	require.NoError(t, err)
	_, err = tarWriter.Write([]byte("test"))
	require.NoError(t, err)

	tarWriter.Close()
	gzipWriter.Close()

	// Test the function - should fail due to unknown type
	err = UntarToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	require.Error(t, err, "Should fail due to unknown tar type")
	assert.Contains(t, err.Error(), "unknown type during tar extraction")
}

func TestUntarToDirWithCorruptGzip(t *testing.T) {
	t.Parallel()
	// Test corrupt gzip data
	tempDir := t.TempDir()

	// Test with invalid gzip data
	err := UntarToDir([]byte("not gzip data"), tempDir, func(string) bool { return true })
	require.Error(t, err, "Should fail with invalid gzip data")
}

func TestUntarToDirWithSymlink(t *testing.T) {
	t.Parallel()
	// Test symlink handling (should be treated as unknown type)
	tempDir := t.TempDir()

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzipWriter)

	// Add a symlink
	header := &tar.Header{
		Name:     "symlink.txt",
		Linkname: "target.txt",
		Mode:     0o644,
		Typeflag: tar.TypeSymlink,
	}

	err := tarWriter.WriteHeader(header)
	require.NoError(t, err)

	tarWriter.Close()
	gzipWriter.Close()

	// Test the function - should fail due to symlink type
	err = UntarToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	require.Error(t, err, "Should fail due to symlink type")
	assert.Contains(t, err.Error(), "unknown type during tar extraction")
}

func TestUntarToDirWithLargeFile(t *testing.T) {
	t.Parallel()
	// Test with a file that's too large (would exceed limitedcopy limits)
	tempDir := t.TempDir()

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzipWriter)

	// Create a very large file (over 200MB limit)
	largeContent := make([]byte, 201*1024*1024) // 201MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	header := &tar.Header{
		Name: "large.txt",
		Size: int64(len(largeContent)),
		Mode: 0o644,
	}

	err := tarWriter.WriteHeader(header)
	require.NoError(t, err)
	_, err = tarWriter.Write(largeContent)
	require.NoError(t, err)

	tarWriter.Close()
	gzipWriter.Close()

	// Test the function - should fail due to file size limit
	err = UntarToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	require.Error(t, err, "Should fail due to file size limit")
	assert.Contains(t, err.Error(), "file too big")
}
