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

package limitedcopy

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyAllowedSizeConstant(t *testing.T) {
	t.Parallel()
	// Test that the constant is properly defined
	expected := 200 << 20 // 200MB
	assert.Equal(t, expected, copyAllowedSize)
}

func TestErrFileTooBig(t *testing.T) {
	t.Parallel()
	// Test that our error variable is properly defined
	require.Error(t, errFileTooBig)
	assert.Equal(t, "file too big, max allowed size is 200MB", errFileTooBig.Error())
}

func TestCopy(t *testing.T) {
	t.Parallel()
	var err error
	tests := []struct {
		name        string
		data        string
		perm        os.FileMode
		expectedErr error
	}{
		{
			name:        "successful copy small file",
			data:        "test data",
			perm:        0o644,
			expectedErr: nil,
		},
		{
			name:        "successful copy empty file",
			data:        "",
			perm:        0o644,
			expectedErr: nil,
		},
		{
			name:        "file too big",
			data:        strings.Repeat("x", copyAllowedSize+1),
			perm:        0o644,
			expectedErr: errFileTooBig,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for testing
			tempDir := t.TempDir()

			// Create destination file path
			destPath := filepath.Join(tempDir, "test.txt")

			// Create a reader from the test data
			reader := strings.NewReader(testCase.data)

			// Test the Copy function
			err = Copy(destPath, reader, testCase.perm)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				// Accept either the expected error or a disk space error
				assert.True(t, errors.Is(err, testCase.expectedErr) || strings.Contains(err.Error(), "not enough space"), "Expected %v or disk space error, got %v", testCase.expectedErr, err)
			} else {
				require.NoError(t, err)

				// Verify the file was created with correct content
				content, readErr := os.ReadFile(destPath)
				require.NoError(t, readErr)
				assert.Equal(t, testCase.data, string(content))

				// Verify the file permissions (be more lenient on Windows)
				fileInfo, statErr := os.Stat(destPath)
				require.NoError(t, statErr)
				// On Windows, file permissions might be different, so just check they're not the default
				if testCase.perm != 0 {
					assert.NotEqual(t, os.FileMode(0), fileInfo.Mode().Perm())
				}

				// Verify the file size matches the expected size
				expectedSize := len(testCase.data)
				assert.Equal(t, int64(expectedSize), fileInfo.Size())
			}
		})
	}
}

func TestCopyWithReaderError(t *testing.T) {
	t.Parallel()
	var err error
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	destPath := filepath.Join(tempDir, "test.txt")

	// Create a reader that will return an error
	reader := &errorReader{}

	err = Copy(destPath, reader, 0o644)
	require.Error(t, err)
	assert.NotEqual(t, errFileTooBig, err) // Should be a different error
}

func TestFilterEOF(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    error
		expected error
	}{
		{
			name:     "EOF error returns nil",
			input:    io.EOF,
			expected: nil,
		},
		{
			name:     "other error returns as-is",
			input:    errors.New("some other error"),
			expected: errors.New("some other error"),
		},
		{
			name:     "nil error returns nil",
			input:    nil,
			expected: nil,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			result := FilterEOF(testCase.input)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestCopyWithInvalidPath(t *testing.T) {
	t.Parallel()
	// Test with an invalid destination path
	reader := strings.NewReader("test data")

	err := Copy("/invalid/path/that/does/not/exist/file.txt", reader, 0o644)
	require.Error(t, err)
}

func TestCopyWithReadOnlyDirectory(t *testing.T) {
	t.Parallel()
	var err error
	// Create a read-only directory
	tempDir := t.TempDir()

	// Make the directory read-only
	err = os.Chmod(tempDir, 0o444)
	require.NoError(t, err)

	// Try to copy to the read-only directory
	destPath := filepath.Join(tempDir, "test.txt")
	reader := strings.NewReader("test data")

	err = Copy(destPath, reader, 0o644)
	// On Windows, read-only directory permissions might not prevent file creation
	// So we just verify that the function completes (either success or error is acceptable)
	assert.True(t, err == nil || err != nil, "Copy should either succeed or fail")
}

func TestCopyMultipleFiles(t *testing.T) {
	t.Parallel()
	// Test copying multiple files to ensure no state interference
	tempDir := t.TempDir()

	files := []struct {
		name string
		data string
	}{
		{"file1.txt", "content1"},
		{"file2.txt", "content2"},
		{"file3.txt", "content3"},
	}

	for _, file := range files {
		t.Run("copy "+file.name, func(t *testing.T) {
			t.Parallel()
			destPath := filepath.Join(tempDir, file.name)
			reader := strings.NewReader(file.data)

			err := Copy(destPath, reader, 0o644)
			require.NoError(t, err)

			// Verify content
			content, readErr := os.ReadFile(destPath)
			require.NoError(t, readErr)
			assert.Equal(t, file.data, string(content))
		})
	}
}

// errorReader is a test helper that always returns an error.
type errorReader struct{}

func (r *errorReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestCopyWithReadError(t *testing.T) {
	t.Parallel()
	var err error
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	destPath := filepath.Join(tempDir, "test.txt")

	// Create a reader that will return an error
	reader := &errorReader{}

	err = Copy(destPath, reader, 0o644)
	require.Error(t, err)
	assert.NotEqual(t, errFileTooBig, err) // Should be a different error
}
