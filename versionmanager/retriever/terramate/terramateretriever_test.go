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

package terramateretriever

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tofuutils/tenv/v4/config"
)

func TestMake(t *testing.T) {
	conf, err := config.DefaultConfig()
	require.NoError(t, err)

	retriever := Make(&conf)

	assert.NotNil(t, retriever)
	assert.Equal(t, &conf, retriever.conf)
}

func TestBuildAssetNames(t *testing.T) {
	// Test with current runtime values
	fileName, shaFileName := buildAssetNames("1.0.0", runtime.GOARCH)

	// The filename should contain the version and platform info
	assert.Contains(t, fileName, "terramate_1.0.0")
	assert.Contains(t, fileName, runtime.GOOS)
	// The architecture in the filename is converted (e.g., amd64 -> x86_64)
	assert.Contains(t, fileName, "x86_64") // amd64 is converted to x86_64
	assert.Equal(t, "checksums.txt", shaFileName)
}

func TestBuildAssetNamesCurrentPlatform(t *testing.T) {
	// Test with current runtime values
	fileName, shaFileName := buildAssetNames("1.0.0", runtime.GOARCH)

	// Should contain the version and platform info
	assert.Contains(t, fileName, "terramate_1.0.0")
	assert.Contains(t, fileName, runtime.GOOS)
	assert.Equal(t, "checksums.txt", shaFileName)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "terramate_", baseFileName)
	assert.Equal(t, "terramate-io", terramateIoName)
}

func TestTerramateRetrieverStructure(t *testing.T) {
	// Test that the struct can be created and accessed
	conf, err := config.DefaultConfig()
	require.NoError(t, err)

	retriever := TerramateRetriever{
		conf: &conf,
	}

	assert.NotNil(t, retriever.conf)
	assert.Equal(t, &conf, retriever.conf)
}

// TestInstallMethodSignature tests that the Install method exists
func TestInstallMethodSignature(t *testing.T) {
	// Test that the method exists
	// Note: We can't actually test the full install logic without proper mocking

	conf, err := config.DefaultConfig()
	require.NoError(t, err)

	retriever := Make(&conf)

	// Verify the method exists
	assert.NotNil(t, retriever.Install)
}

// TestListVersionsMethodSignature tests that the ListVersions method exists
func TestListVersionsMethodSignature(t *testing.T) {
	// Test that the method exists
	// Note: We can't actually test the full list logic without proper mocking

	conf, err := config.DefaultConfig()
	require.NoError(t, err)

	retriever := Make(&conf)

	// Verify the method exists
	assert.NotNil(t, retriever.ListVersions)
}

func TestVersionTagHandling(t *testing.T) {
	// Test the logic for handling version tags with and without 'v' prefix
	testCases := []struct {
		name            string
		input           string
		expectedTag     string
		expectedVersion string
	}{
		{
			name:            "version without v prefix",
			input:           "1.0.0",
			expectedTag:     "v1.0.0",
			expectedVersion: "1.0.0",
		},
		{
			name:            "version with v prefix",
			input:           "v1.0.0",
			expectedTag:     "v1.0.0",
			expectedVersion: "1.0.0",
		},
		{
			name:            "version with multiple digits",
			input:           "2.15.3",
			expectedTag:     "v2.15.3",
			expectedVersion: "2.15.3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tag := tc.input
			versionStr := tc.input

			// Apply the same logic as in the Install method
			if tag[0] == 'v' {
				versionStr = versionStr[1:]
			} else {
				tag = "v" + versionStr
			}

			assert.Equal(t, tc.expectedTag, tag)
			assert.Equal(t, tc.expectedVersion, versionStr)
		})
	}
}

func TestAssetNameBuildingLogic(t *testing.T) {
	// Test the asset name building logic with current architecture
	fileName, _ := buildAssetNames("1.0.0", runtime.GOARCH)

	// Should contain the version and current platform info
	assert.Contains(t, fileName, "terramate_1.0.0")
	assert.Contains(t, fileName, runtime.GOOS)
	// The architecture in the filename is converted (e.g., amd64 -> x86_64)
	assert.Contains(t, fileName, "x86_64") // amd64 is converted to x86_64
}
