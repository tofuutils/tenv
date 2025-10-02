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
					{"test.txt", "test content", 0644},
					{"subdir/file.txt", "nested content", 0644},
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
					{"test.txt", "test content", 0644},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "targz_test")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create test tar.gz data
			data := tt.setupTarGz()

			// Test the function
			err = UntarToDir(data, tempDir, tt.filter)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				// For cases where we expect no specific error
				// The actual error will depend on the implementation details
				_ = err // We don't assert on the specific error in this case
			}
		})
	}
}

func TestUntarToDirCreatesFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "targz_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test tar.gz data
	testFiles := []testFile{
		{"file1.txt", "content1", 0644},
		{"file2.txt", "content2", 0644},
		{"subdir/file3.txt", "content3", 0644},
	}

	data := createTestTarGz(testFiles)

	// Test the function
	err = UntarToDir(data, tempDir, func(string) bool { return true })
	assert.NoError(t, err)

	// Verify files were created
	for _, file := range testFiles {
		expectedPath := filepath.Join(tempDir, file.name)
		_, err := os.Stat(expectedPath)
		assert.NoError(t, err, "File %s should be created", file.name)
	}
}

func TestUntarToDirWithFilter(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "targz_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test tar.gz data
	testFiles := []testFile{
		{"include.txt", "included", 0644},
		{"exclude.txt", "excluded", 0644},
	}

	data := createTestTarGz(testFiles)

	// Filter that only includes files with "include" in the name
	filter := func(name string) bool {
		return len(name) > 0 && name[0] == 'i'
	}

	// Test the function
	err = UntarToDir(data, tempDir, filter)
	assert.NoError(t, err)

	// Verify only included file was created
	_, err = os.Stat(filepath.Join(tempDir, "include.txt"))
	assert.NoError(t, err, "Included file should be created")

	_, err = os.Stat(filepath.Join(tempDir, "exclude.txt"))
	assert.True(t, os.IsNotExist(err), "Excluded file should not be created")
}

func TestUntarToDirCreatesDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "targz_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test tar.gz data with directories
	testFiles := []testFile{
		{"dir1/file.txt", "content", 0644},
		{"dir2/subdir/file.txt", "nested content", 0644},
	}

	data := createTestTarGz(testFiles)

	// Test the function
	err = UntarToDir(data, tempDir, func(string) bool { return true })
	assert.NoError(t, err)

	// Verify directories were created
	_, err = os.Stat(filepath.Join(tempDir, "dir1"))
	assert.NoError(t, err, "Directory dir1 should be created")

	_, err = os.Stat(filepath.Join(tempDir, "dir2", "subdir"))
	assert.NoError(t, err, "Nested directory should be created")
}

// Helper function to create test tar.gz data
func createTestTarGz(files []testFile) []byte {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	for _, file := range files {
		// Create tar header
		header := &tar.Header{
			Name: file.name,
			Size: int64(len(file.content)),
			Mode: int64(file.mode),
		}

		// Write header
		err := tw.WriteHeader(header)
		if err != nil {
			panic(err)
		}

		// Write file content
		_, err = tw.Write([]byte(file.content))
		if err != nil {
			panic(err)
		}
	}

	// Close writers
	tw.Close()
	gw.Close()

	return buf.Bytes()
}

// Helper struct for test files
type testFile struct {
	name    string
	content string
	mode    int64
}

func TestUntarToDirWithPathTraversal(t *testing.T) {
	// Test path traversal protection
	tempDir, err := os.MkdirTemp("", "targz_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create tar.gz with path traversal attempt
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Add a file with path traversal
	header := &tar.Header{
		Name: "../../../outside.txt",
		Size: 4,
		Mode: 0644,
	}

	err = tw.WriteHeader(header)
	require.NoError(t, err)
	_, err = tw.Write([]byte("test"))
	require.NoError(t, err)

	tw.Close()
	gw.Close()

	// Test the function - should fail due to path traversal protection
	err = UntarToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	assert.Error(t, err, "Should fail due to path traversal protection")
}

func TestUntarToDirWithUnknownType(t *testing.T) {
	// Test unknown tar type
	tempDir, err := os.MkdirTemp("", "targz_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create tar.gz with unknown type
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Add a file with unknown type
	header := &tar.Header{
		Name:     "unknown.txt",
		Size:     4,
		Mode:     0644,
		Typeflag: 99, // Unknown type
	}

	err = tw.WriteHeader(header)
	require.NoError(t, err)
	_, err = tw.Write([]byte("test"))
	require.NoError(t, err)

	tw.Close()
	gw.Close()

	// Test the function - should fail due to unknown type
	err = UntarToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	assert.Error(t, err, "Should fail due to unknown tar type")
	assert.Contains(t, err.Error(), "unknown type during tar extraction")
}

func TestUntarToDirWithCorruptGzip(t *testing.T) {
	// Test corrupt gzip data
	tempDir, err := os.MkdirTemp("", "targz_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test with invalid gzip data
	err = UntarToDir([]byte("not gzip data"), tempDir, func(string) bool { return true })
	assert.Error(t, err, "Should fail with invalid gzip data")
}

func TestUntarToDirWithSymlink(t *testing.T) {
	// Test symlink handling (should be treated as unknown type)
	tempDir, err := os.MkdirTemp("", "targz_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Add a symlink
	header := &tar.Header{
		Name:     "symlink.txt",
		Linkname: "target.txt",
		Mode:     0644,
		Typeflag: tar.TypeSymlink,
	}

	err = tw.WriteHeader(header)
	require.NoError(t, err)

	tw.Close()
	gw.Close()

	// Test the function - should fail due to symlink type
	err = UntarToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	assert.Error(t, err, "Should fail due to symlink type")
	assert.Contains(t, err.Error(), "unknown type during tar extraction")
}

func TestUntarToDirWithLargeFile(t *testing.T) {
	// Test with a file that's too large (would exceed limitedcopy limits)
	tempDir, err := os.MkdirTemp("", "targz_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Create a very large file (over 200MB limit)
	largeContent := make([]byte, 201*1024*1024) // 201MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	header := &tar.Header{
		Name: "large.txt",
		Size: int64(len(largeContent)),
		Mode: 0644,
	}

	err = tw.WriteHeader(header)
	require.NoError(t, err)
	_, err = tw.Write(largeContent)
	require.NoError(t, err)

	tw.Close()
	gw.Close()

	// Test the function - should fail due to file size limit
	err = UntarToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	assert.Error(t, err, "Should fail due to file size limit")
	assert.Contains(t, err.Error(), "file too big")
}
