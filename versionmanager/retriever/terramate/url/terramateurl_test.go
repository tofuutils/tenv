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

package terramateurl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubURLConstant(t *testing.T) {
	// Test that the Github constant is properly constructed
	expected := "https://api.github.com/repos/terramate-io/terramate/releases"
	actual := Github

	assert.Equal(t, expected, actual)
	assert.Contains(t, actual, "github.com")
	assert.Contains(t, actual, "terramate-io/terramate")
	assert.Contains(t, actual, "/releases")
	assert.True(t, strings.HasPrefix(actual, "https://api.github.com/repos/"))
	assert.True(t, strings.HasSuffix(actual, "/releases"))
}

func TestGithubURLStructure(t *testing.T) {
	// Test that the URL follows the expected GitHub API pattern
	url := Github

	// Should be a valid GitHub API releases URL
	assert.Contains(t, url, "api.github.com")
	assert.Contains(t, url, "repos")
	assert.Contains(t, url, "releases")

	// Should point to the correct repository
	assert.Contains(t, url, "terramate-io/terramate")

	// Should be a complete URL
	assert.True(t, len(url) > 20, "URL should be reasonably long")
	assert.True(t, url[:8] == "https://", "URL should start with https://")
}
