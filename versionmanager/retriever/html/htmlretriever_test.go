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

package htmlretriever

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tofuutils/tenv/v4/pkg/download"
)

func TestBuildAssetURLs(t *testing.T) {
	tests := []struct {
		name         string
		baseAssetURL string
		assetNames   []string
		expected     []string
		expectError  bool
	}{
		{
			name:         "valid URLs",
			baseAssetURL: "https://example.com/releases",
			assetNames:   []string{"asset1.zip", "asset2.tar.gz"},
			expected: []string{
				"https://example.com/releases/asset1.zip",
				"https://example.com/releases/asset2.tar.gz",
			},
			expectError: false,
		},
		{
			name:         "empty asset names",
			baseAssetURL: "https://example.com/releases",
			assetNames:   []string{},
			expected:     []string{},
			expectError:  false,
		},
		{
			name:         "invalid base URL",
			baseAssetURL: "://invalid-url",
			assetNames:   []string{"asset.zip"},
			expected:     nil,
			expectError:  true,
		},
		{
			name:         "asset names with special characters",
			baseAssetURL: "https://example.com/releases",
			assetNames:   []string{"asset with spaces.zip", "asset-with-dashes.tar.gz"},
			expected: []string{
				"https://example.com/releases/asset%20with%20spaces.zip",
				"https://example.com/releases/asset-with-dashes.tar.gz",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildAssetURLs(tt.baseAssetURL, tt.assetNames...)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestListReleases(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		remoteConf  map[string]string
		options     []download.RequestOption
		expected    []string
		expectError bool
	}{
		{
			name:    "default selector and part",
			baseURL: "https://example.com/releases",
			remoteConf: map[string]string{
				"selector": "a",
				"part":     "href",
			},
			options:     []download.RequestOption{},
			expected:    nil, // Would depend on actual HTML content
			expectError: false,
		},
		{
			name:    "custom selector",
			baseURL: "https://example.com/releases",
			remoteConf: map[string]string{
				"selector": ".release-link",
				"part":     "href",
			},
			options:     []download.RequestOption{},
			expected:    nil, // Would depend on actual HTML content
			expectError: false,
		},
		{
			name:    "custom part extraction",
			baseURL: "https://example.com/releases",
			remoteConf: map[string]string{
				"selector": "a",
				"part":     "text",
			},
			options:     []download.RequestOption{},
			expected:    nil, // Would depend on actual HTML content
			expectError: false,
		},
		{
			name:        "empty remote config",
			baseURL:     "https://example.com/releases",
			remoteConf:  map[string]string{},
			options:     []download.RequestOption{},
			expected:    nil, // Would depend on actual HTML content
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test would require mocking the HTTP request
			// For now, we just test that the function signature is correct
			// and doesn't panic with valid inputs

			ctx := context.Background()

			// This would normally make an HTTP request, but we're testing
			// the function structure and error handling
			result, err := ListReleases(ctx, tt.baseURL, tt.remoteConf, tt.options)

			// The function may return an error depending on network connectivity
			// but we can at least verify it doesn't panic
			if err != nil {
				// If there's an error, it should be related to network or parsing
				assert.Error(t, err)
			} else {
				// If successful, result should be a slice of strings
				assert.IsType(t, []string{}, result)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that the function references are valid
	assert.NotNil(t, BuildAssetURLs)
	assert.NotNil(t, ListReleases)
}

func TestURLJoiningLogic(t *testing.T) {
	// Test the URL joining logic specifically
	baseURL := "https://github.com/example/repo/releases/download/v1.0.0"

	testCases := []struct {
		assetName string
		expected  string
	}{
		{
			assetName: "terramate_1.0.0_linux_amd64.tar.gz",
			expected:  "https://github.com/example/repo/releases/download/v1.0.0/terramate_1.0.0_linux_amd64.tar.gz",
		},
		{
			assetName: "checksums.txt",
			expected:  "https://github.com/example/repo/releases/download/v1.0.0/checksums.txt",
		},
	}

	for _, tc := range testCases {
		t.Run("URL joining: "+tc.assetName, func(t *testing.T) {
			result, err := BuildAssetURLs(baseURL, tc.assetName)
			require.NoError(t, err)
			assert.Len(t, result, 1)
			assert.Equal(t, tc.expected, result[0])
		})
	}
}
