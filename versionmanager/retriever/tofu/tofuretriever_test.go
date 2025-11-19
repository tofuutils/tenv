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

package tofuretriever

import (
	"testing"

	"github.com/hashicorp/go-version"
)

func TestBuildIdentity(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name     string
		version  *version.Version
		expected string
	}{
		"stable version": {
			name:     "1.11.0",
			version:  version.Must(version.NewVersion("1.11.0")),
			expected: "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v1.11",
		},
		"unstable alpha version": {
			name:     "1.7.0-alpha1",
			version:  version.Must(version.NewVersion("1.7.0-alpha1")),
			expected: "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/main",
		},
		"unstable beta version": {
			name:     "1.7.0-beta1",
			version:  version.Must(version.NewVersion("1.7.0-beta1")),
			expected: "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/main",
		},
		"unstable rc version": {
			name:     "1.7.0-rc1",
			version:  version.Must(version.NewVersion("1.7.0-rc1")),
			expected: "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v1.7",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := buildIdentity(test.version)
			if actual != test.expected {
				t.Errorf("expected %s, got %s", test.expected, actual)
			}
		})
	}
}
