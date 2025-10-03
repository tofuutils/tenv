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

package lastuse

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestRead(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a mock config with inert displayer
	conf := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	tests := []struct {
		name        string
		fileContent string
		expected    time.Time
	}{
		{
			name:        "valid date",
			fileContent: "2024-01-15",
			expected:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "another valid date",
			fileContent: "2023-12-31",
			expected:    time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "invalid date format",
			fileContent: "invalid-date",
			expected:    time.Time{}, // Should return zero time
		},
		{
			name:        "empty file",
			fileContent: "",
			expected:    time.Time{}, // Should return zero time
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create unique test file for this subtest
			testFile := filepath.Join(tempDir, testCase.name, fileName)
			err := os.MkdirAll(filepath.Dir(testFile), 0o755)
			require.NoError(t, err)
			err = os.WriteFile(testFile, []byte(testCase.fileContent), 0o600)
			require.NoError(t, err)

			// Test Read function
			result := Read(filepath.Dir(testFile), conf)
			assert.Equal(t, testCase.expected, result)
		})
	}

	// Test with non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		t.Parallel()
		nonExistentDir := filepath.Join(tempDir, "non-existent")
		result := Read(nonExistentDir, conf)
		assert.Equal(t, time.Time{}, result)
	})
}

func TestWriteNow(t *testing.T) {
	t.Parallel()
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a mock config with inert displayer and empty getenv function
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
		Getenv:    config.EmptyGetenv,
	}

	// Test WriteNow function
	WriteNow(tempDir, mockConfig)

	// Check if file was written
	testFile := filepath.Join(tempDir, fileName)
	_, err := os.Stat(testFile)
	require.NoError(t, err, "File should be written")

	// Verify the file contains a valid date
	content, readErr := os.ReadFile(testFile)
	require.NoError(t, readErr)
	_, parseErr := time.Parse(time.DateOnly, string(content))
	require.NoError(t, parseErr, "File should contain valid date")

	// Test that function can be called multiple times without error
	WriteNow(tempDir, mockConfig)

	// Verify the file still exists after second call
	_, err = os.Stat(testFile)
	assert.NoError(t, err, "File should still exist after second call")
}

func TestConstants(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "last-use.txt", fileName)
	assert.Equal(t, "Unable to retrieve TENV_SKIP_LAST_USE environment variable", skipLastUseErrMsg)
}
