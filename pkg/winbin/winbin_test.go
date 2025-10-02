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

package winbin

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArchiveFormat(t *testing.T) {
	// Test constants
	assert.Equal(t, ".exe", suffix)
	assert.Equal(t, "windows", osName)
	assert.Equal(t, ".zip", zipSuffix)
	assert.Equal(t, ".tar.gz", tarGzSuffix)

	// Test actual runtime behavior
	result := GetArchiveFormat()

	// On Windows, should return .zip
	// On other systems, should return .tar.gz
	if runtime.GOOS == "windows" {
		assert.Equal(t, ".zip", result)
	} else {
		assert.Equal(t, ".tar.gz", result)
	}

	// Test that function is deterministic
	result2 := GetArchiveFormat()
	assert.Equal(t, result, result2)
}

func TestGetBinaryName(t *testing.T) {
	tests := []struct {
		name     string
		execName string
	}{
		{
			name:     "terraform binary",
			execName: "terraform",
		},
		{
			name:     "tofu binary",
			execName: "tofu",
		},
		{
			name:     "empty exec name",
			execName: "",
		},
		{
			name:     "already has exe suffix",
			execName: "terraform.exe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBinaryName(tt.execName)

			// On Windows, should add .exe suffix
			// On other systems, should return original name
			if runtime.GOOS == "windows" {
				assert.True(t, strings.HasSuffix(result, ".exe"))
			} else {
				assert.Equal(t, tt.execName, result)
			}

			// Test that function is deterministic
			result2 := GetBinaryName(tt.execName)
			assert.Equal(t, result, result2)
		})
	}
}

func TestWriteSuffixTo(t *testing.T) {
	var result strings.Builder

	// Test actual runtime behavior
	n, err := WriteSuffixTo(&result)

	// Should not return an error
	assert.NoError(t, err)

	// On Windows, should write .exe suffix
	// On other systems, should write nothing
	if runtime.GOOS == "windows" {
		assert.Equal(t, 4, n)
		assert.Equal(t, ".exe", result.String())
	} else {
		assert.Equal(t, 0, n)
		assert.Equal(t, "", result.String())
	}

	// Test that function is deterministic
	var result2 strings.Builder
	n2, err2 := WriteSuffixTo(&result2)
	assert.NoError(t, err2)
	assert.Equal(t, n, n2)
	assert.Equal(t, result.String(), result2.String())
}

// TestConstants verifies that all constants are properly defined
func TestConstants(t *testing.T) {
	assert.Equal(t, ".exe", suffix)
	assert.Equal(t, "windows", osName)
	assert.Equal(t, ".zip", zipSuffix)
	assert.Equal(t, ".tar.gz", tarGzSuffix)
}

// TestCrossPlatformBehavior tests the behavior across different platforms
// This is a conceptual test since runtime.GOOS can't be easily mocked
func TestCrossPlatformBehavior(t *testing.T) {
	// Test that functions behave consistently
	archiveFormat := GetArchiveFormat()
	binaryName := GetBinaryName("test")
	var suffixWriter strings.Builder
	suffixLength, _ := WriteSuffixTo(&suffixWriter)

	// Verify that the functions don't panic and return reasonable values
	assert.NotEmpty(t, archiveFormat)
	assert.NotNil(t, binaryName)
	assert.True(t, suffixLength >= 0)

	// Test multiple calls return consistent results
	archiveFormat2 := GetArchiveFormat()
	assert.Equal(t, archiveFormat, archiveFormat2)
}
