/*
 *
 * Copyright 2025 tofuutils authors.
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

package githuburl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "Base constant",
			constant: Base,
			expected: "https://github.com",
		},
		{
			name:     "Default constant",
			constant: Default,
			expected: "https://api.github.com/repos/",
		},
		{
			name:     "SlashReleasesSuffix constant",
			constant: SlashReleasesSuffix,
			expected: "/releases",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}
