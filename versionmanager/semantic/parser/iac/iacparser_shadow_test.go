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

package iacparser

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func gatherFromFiles(t *testing.T, files map[string]string, tofu bool) []string {
	t.Helper()

	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	hclParser := hclparse.NewParser()
	exts := []ExtDescription{
		{Value: ".tf", Parser: hclParser.ParseHCLFile},
		{Value: ".tf.json", Parser: hclParser.ParseJSONFile},
	}
	if tofu {
		exts = append([]ExtDescription{
			{Value: ".tofu", Parser: hclParser.ParseHCLFile},
			{Value: ".tofu.json", Parser: hclParser.ParseJSONFile},
		}, exts...)
	}

	conf := &config.Config{WorkPath: dir, Displayer: loghelper.InertDisplayer}
	versions, err := GatherRequiredVersion(conf, exts)
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(versions)

	return versions
}

func TestGatherRequiredVersionSameBaseName(t *testing.T) {
	t.Parallel()

	const (
		hclFile     = "terraform {\n  required_version = \">= 1.5.0\"\n}\n"
		hclTofuFile = "terraform {\n  required_version = \">= 1.8.0\"\n}\n"
		jsonFile    = `{"terraform": {"required_version": "< 2.0.0"}}`
		jsonTofu    = `{"terraform": {"required_version": "< 1.10.0"}}`
	)

	tests := []struct {
		name     string
		files    map[string]string
		tofu     bool
		expected []string
	}{
		{
			name:     "TfAndTfJSONAreBothRead",
			files:    map[string]string{"main.tf": "", "main.tf.json": jsonFile},
			expected: []string{"< 2.0.0"},
		},
		{
			name:     "TfAndTfJSONBothWithVersion",
			files:    map[string]string{"main.tf": hclFile, "main.tf.json": jsonFile},
			expected: []string{"< 2.0.0", ">= 1.5.0"},
		},
		{
			name:     "TofuHidesTf",
			files:    map[string]string{"main.tf": hclFile, "main.tofu": hclTofuFile},
			tofu:     true,
			expected: []string{">= 1.8.0"},
		},
		{
			name:     "TofuJSONHidesTfJSONOnly",
			files:    map[string]string{"main.tf": hclFile, "main.tf.json": jsonFile, "main.tofu.json": jsonTofu},
			tofu:     true,
			expected: []string{"< 1.10.0", ">= 1.5.0"},
		},
		{
			name:     "TofuDoesNotHideTfJSON",
			files:    map[string]string{"main.tofu": hclTofuFile, "main.tf.json": jsonFile},
			tofu:     true,
			expected: []string{"< 2.0.0", ">= 1.8.0"},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			versions := gatherFromFiles(t, testCase.files, testCase.tofu)
			if !slices.Equal(versions, testCase.expected) {
				t.Fatalf("unexpected versions %q, want %q", versions, testCase.expected)
			}
		})
	}
}
