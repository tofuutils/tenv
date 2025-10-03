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

package apimsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "AssetsName constant",
			constant: AssetsName,
			expected: "assets",
		},
		{
			name:     "MsgFetchAllReleases constant",
			constant: MsgFetchAllReleases,
			expected: "Fetching all releases information from ",
		},
		{
			name:     "MsgFetchRelease constant",
			constant: MsgFetchRelease,
			expected: "Fetching release information from ",
		},
		{
			name:     "MsgSearch constant",
			constant: MsgSearch,
			expected: "Search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

func TestErrorVariables(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrAsset error",
			err:      ErrAsset,
			expected: "searched asset not found",
		},
		{
			name:     "ErrReturn error",
			err:      ErrReturn,
			expected: "unexpected value returned by API",
		},
		{
			name:     "ErrRateLimit error",
			err:      ErrRateLimit,
			expected: "you are rate-limited by GitHub. Consider using a token by setting the TENV_GITHUB_TOKEN env variable to increase the rate limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.err.Error())
			assert.Error(t, tt.err)
		})
	}
}
