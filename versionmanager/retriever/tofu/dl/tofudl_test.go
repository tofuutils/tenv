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

package tofudlmirroring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeURLBuilderAndBuild(t *testing.T) {
	t.Parallel()

	builder, err := MakeURLBuilder("https://github.com/opentofu/opentofu/releases/download/v{{ .Version }}/{{ .Artifact }}", "1.6.0")
	require.NoError(t, err)

	url, err := builder.Build("tofu_1.6.0_linux_amd64.zip")
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_linux_amd64.zip", url)
}

func TestMakeURLBuilder_InvalidTemplate(t *testing.T) {
	t.Parallel()

	_, err := MakeURLBuilder("{{ .Version ", "1.6.0")
	require.Error(t, err)
}

func TestBuild_UnknownField(t *testing.T) {
	t.Parallel()

	builder, err := MakeURLBuilder("{{ .NotAField }}", "1.6.0")
	require.NoError(t, err)

	_, err = builder.Build("artifact")
	require.Error(t, err)
}

func TestExtractReleases(t *testing.T) {
	t.Parallel()

	value := map[string]any{
		"versions": []any{
			map[string]any{"id": "1.6.0"},
			map[string]any{"id": "1.7.0"},
		},
	}

	releases, err := ExtractReleases(value)
	require.NoError(t, err)
	assert.Equal(t, []string{"1.6.0", "1.7.0"}, releases)
}

func TestExtractReleases_EmptyVersions(t *testing.T) {
	t.Parallel()

	value := map[string]any{"versions": []any{}}

	releases, err := ExtractReleases(value)
	require.NoError(t, err)
	assert.Empty(t, releases)
}

func TestExtractReleases_MissingVersionsField(t *testing.T) {
	t.Parallel()

	_, err := ExtractReleases(map[string]any{"other": "field"})
	require.Error(t, err)
}

func TestExtractReleases_NotAnObject(t *testing.T) {
	t.Parallel()

	_, err := ExtractReleases("not a map")
	require.Error(t, err)
}

func TestExtractReleases_VersionIDNotAString(t *testing.T) {
	t.Parallel()

	value := map[string]any{
		"versions": []any{
			map[string]any{"id": 42},
		},
	}

	_, err := ExtractReleases(value)
	require.Error(t, err)
}
