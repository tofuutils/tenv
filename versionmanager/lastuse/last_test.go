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
	"github.com/tofuutils/tenv/v4/config/envname"
	configutils "github.com/tofuutils/tenv/v4/config/utils"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestRead(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tenv-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, fileName)
			err := os.WriteFile(testFile, []byte(tt.fileContent), 0644)
			require.NoError(t, err)

			// Test Read function
			result := Read(tempDir, conf)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test with non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "non-existent")
		result := Read(nonExistentDir, conf)
		assert.Equal(t, time.Time{}, result)
	})
}

func TestWriteNow(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tenv-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock config with inert displayer and empty getenv function
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
		Getenv:    config.EmptyGetenv,
	}

	// Test WriteNow function
	WriteNow(tempDir, mockConfig)

	// Check if file was written
	testFile := filepath.Join(tempDir, fileName)
	_, err = os.Stat(testFile)
	assert.NoError(t, err, "File should be written")

	// Verify the file contains a valid date
	content, readErr := os.ReadFile(testFile)
	assert.NoError(t, readErr)
	_, parseErr := time.Parse(time.DateOnly, string(content))
	assert.NoError(t, parseErr, "File should contain valid date")

	// Test that function can be called multiple times without error
	WriteNow(tempDir, mockConfig)
	assert.NoError(t, err, "Second call to WriteNow should not error")
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "last-use.txt", fileName)
	assert.Equal(t, "Unable to retrieve TENV_SKIP_LAST_USE environment variable", skipLastUseErrMsg)
}

// createMockGetenvFunc creates a mock GetenvFunc for testing
func createMockGetenvFunc(envValue string) configutils.GetenvFunc {
	return func(key string) string {
		if key != envname.TenvSkipLastUse {
			return ""
		}

		return envValue
	}
}
