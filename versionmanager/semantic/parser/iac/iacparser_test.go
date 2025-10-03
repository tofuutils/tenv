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
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestFilterExts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		fileExts int
		exts     []ExtDescription
		expected ExtDescription
	}{
		{
			name:     "single extension",
			fileExts: 1,
			exts: []ExtDescription{
				{Value: ".tf", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
					return &hcl.File{}, hcl.Diagnostics{}
				}},
			},
			expected: ExtDescription{Value: ".tf", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
				return &hcl.File{}, hcl.Diagnostics{}
			}},
		},
		{
			name:     "multiple extensions - first match",
			fileExts: 1,
			exts: []ExtDescription{
				{Value: ".tf", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
					return &hcl.File{}, hcl.Diagnostics{}
				}},
				{Value: ".hcl", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
					return &hcl.File{}, hcl.Diagnostics{}
				}},
			},
			expected: ExtDescription{Value: ".tf", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
				return &hcl.File{}, hcl.Diagnostics{}
			}},
		},
		{
			name:     "multiple extensions - second match",
			fileExts: 2,
			exts: []ExtDescription{
				{Value: ".tf", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
					return &hcl.File{}, hcl.Diagnostics{}
				}},
				{Value: ".hcl", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
					return &hcl.File{}, hcl.Diagnostics{}
				}},
				{Value: ".tfvars", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
					return &hcl.File{}, hcl.Diagnostics{}
				}},
			},
			expected: ExtDescription{Value: ".hcl", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
				return &hcl.File{}, hcl.Diagnostics{}
			}},
		},
		{
			name:     "no matching extensions",
			fileExts: 0,
			exts:     []ExtDescription{},
			expected: ExtDescription{Parser: nil},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			result := filterExts(testCase.fileExts, testCase.exts)
			// Compare only the Value field since Parser functions have different memory addresses
			assert.Equal(t, testCase.expected.Value, result.Value, "Value should match")
			if testCase.expected.Parser != nil {
				assert.NotNil(t, result.Parser, "Parser should not be nil")
			}
		})
	}
}

func TestExtractRequiredVersion(t *testing.T) {
	t.Parallel()
	// Create a mock config with inert displayer
	conf := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Test with nil body (should return nil due to error handling)
	versions := extractRequiredVersion(nil, conf)
	assert.Nil(t, versions)
}

func TestGatherRequiredVersion(t *testing.T) {
	t.Parallel()
	var err error
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name             string
		files            map[string]string // filename -> content
		exts             []ExtDescription
		expectedVersions []string
		expectError      bool
	}{
		{
			name:             "no extensions provided",
			files:            map[string]string{},
			exts:             []ExtDescription{},
			expectedVersions: nil,
			expectError:      false,
		},
		{
			name: "no matching files",
			files: map[string]string{
				"main.txt":  "not a terraform file",
				"readme.md": "# Documentation",
			},
			exts: []ExtDescription{
				{Value: ".tf", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
					// Return empty file for non-matching files
					return &hcl.File{}, hcl.Diagnostics{}
				}},
			},
			expectedVersions: []string{},
			expectError:      false,
		},
		{
			name:  "empty directory",
			files: map[string]string{},
			exts: []ExtDescription{
				{Value: ".tf", Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
					// Return empty file for non-matching files
					return &hcl.File{}, hcl.Diagnostics{}
				}},
			},
			expectedVersions: []string{},
			expectError:      false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a subdirectory to avoid scanning the entire tenv directory
			testSubDir := filepath.Join(tempDir, "test")
			err = os.MkdirAll(testSubDir, 0o755)
			require.NoError(t, err)

			// Create test files in the subdirectory
			for filename, content := range testCase.files {
				filePath := filepath.Join(testSubDir, filename)
				err := os.WriteFile(filePath, []byte(content), 0o600)
				require.NoError(t, err)
			}

			// Create config pointing to the subdirectory
			conf := &config.Config{
				WorkPath:  testSubDir,
				Displayer: loghelper.InertDisplayer,
			}

			// Test the function
			versions, err := GatherRequiredVersion(conf, testCase.exts)

			if testCase.expectError {
				require.Error(t, err)
				assert.Nil(t, versions)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedVersions, versions)
			}
		})
	}
}

func TestGatherRequiredVersionWithNoExtensions(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	conf := &config.Config{
		WorkPath:  tempDir,
		Displayer: loghelper.InertDisplayer,
	}

	versions, err := GatherRequiredVersion(conf, []ExtDescription{})
	require.NoError(t, err)
	assert.Nil(t, versions)
}

func TestGatherRequiredVersionWithValidTerraformFiles(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		files            map[string]string // filename -> content
		expectedVersions []string
		expectError      bool
	}{
		{
			name: "single terraform file with required_version",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = ">= 1.0.0"
				}`,
			},
			expectedVersions: []string{">= 1.0.0"},
			expectError:      false,
		},
		{
			name: "multiple terraform files with different versions",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = ">= 1.0.0"
				}`,
				"versions.tf": `terraform {
					required_version = "~> 1.5.0"
				}`,
			},
			expectedVersions: []string{">= 1.0.0", "~> 1.5.0"},
			expectError:      false,
		},
		{
			name: "terraform file without required_version",
			files: map[string]string{
				"main.tf": `terraform {
					required_providers {
						aws = "~> 4.0"
					}
				}`,
			},
			expectedVersions: []string{},
			expectError:      false,
		},
		{
			name: "terraform file with empty required_version",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = ""
				}`,
			},
			expectedVersions: []string{""},
			expectError:      false,
		},
		{
			name: "mixed terraform and non-terraform files",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = ">= 1.0.0"
				}`,
				"readme.md": "# Documentation",
				"script.py": "print('hello')",
			},
			expectedVersions: []string{">= 1.0.0"},
			expectError:      false,
		},
		{
			name: "terraform file with complex version constraint",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = ">= 1.0.0, < 2.0.0"
				}`,
			},
			expectedVersions: []string{">= 1.0.0, < 2.0.0"},
			expectError:      false,
		},
		{
			name: "terraform file with exact version",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = "1.5.0"
				}`,
			},
			expectedVersions: []string{"1.5.0"},
			expectError:      false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for this test
			tempDir := t.TempDir()

			// Create test files
			for filename, content := range testCase.files {
				filePath := filepath.Join(tempDir, filename)
				err := os.WriteFile(filePath, []byte(content), 0o600)
				require.NoError(t, err)
			}

			// Create config pointing to the tempDir
			conf := &config.Config{
				WorkPath:  tempDir,
				Displayer: loghelper.InertDisplayer,
			}

			// Create HCL parser extension
			exts := []ExtDescription{
				{
					Value: ".tf",
					Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
						// Parse the HCL content
						parser := hclparse.NewParser()
						content, _ := os.ReadFile(filename)

						return parser.ParseHCL(content, filename)
					},
				},
			}

			// Test the function
			versions, err := GatherRequiredVersion(conf, exts)

			if testCase.expectError {
				require.Error(t, err)
				assert.Nil(t, versions)
			} else {
				require.NoError(t, err)
				assert.ElementsMatch(t, testCase.expectedVersions, versions)
			}
		})
	}
}

func TestGatherRequiredVersionWithInvalidHCL(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	var err error

	// Create a file with invalid HCL syntax
	filePath := filepath.Join(tempDir, "invalid.tf")
	err = os.WriteFile(filePath, []byte(`terraform {
		required_version = ">= 1.0.0"
		invalid_syntax = `), 0o600)
	require.NoError(t, err)

	conf := &config.Config{
		WorkPath:  tempDir,
		Displayer: loghelper.InertDisplayer,
	}

	exts := []ExtDescription{
		{
			Value: ".tf",
			Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
				parser := hclparse.NewParser()
				content, _ := os.ReadFile(filename)

				return parser.ParseHCL(content, filename)
			},
		},
	}

	// Should return error due to invalid HCL
	versions, err := GatherRequiredVersion(conf, exts)
	require.Error(t, err)
	assert.Nil(t, versions)
}

func TestGatherRequiredVersionWithMultipleExtensions(t *testing.T) {
	t.Parallel()
	// Create a subdirectory to avoid scanning the entire tenv directory
	tempDir := t.TempDir()
	var err error

	// Create a subdirectory for test files
	testSubDir := filepath.Join(tempDir, "test")
	err = os.MkdirAll(testSubDir, 0o755)
	require.NoError(t, err)

	// Create files with different extensions
	files := map[string]string{
		"main.tf": `terraform {
			required_version = ">= 1.0.0"
		}`,
		"config.hcl": `terraform {
			required_version = "~> 1.5.0"
		}`,
		"variables.tfvars": `variable "example" {
			default = "value"
		}`,
	}

	for filename, content := range files {
		filePath := filepath.Join(testSubDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0o600)
		require.NoError(t, err)
	}

	conf := &config.Config{
		WorkPath:  testSubDir,
		Displayer: loghelper.InertDisplayer,
	}

	// Test with multiple extensions
	exts := []ExtDescription{
		{
			Value: ".tf",
			Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
				parser := hclparse.NewParser()
				content, _ := os.ReadFile(filename)

				return parser.ParseHCL(content, filename)
			},
		},
		{
			Value: ".hcl",
			Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
				parser := hclparse.NewParser()
				content, _ := os.ReadFile(filename)

				return parser.ParseHCL(content, filename)
			},
		},
	}

	versions, err := GatherRequiredVersion(conf, exts)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{">= 1.0.0", "~> 1.5.0"}, versions)
}

func TestGatherRequiredVersionWithReadDirError(t *testing.T) {
	t.Parallel()
	// Test with non-existent directory
	conf := &config.Config{
		WorkPath:  "/non/existent/directory",
		Displayer: loghelper.InertDisplayer,
	}

	exts := []ExtDescription{
		{
			Value: ".tf",
			Parser: func(filename string) (*hcl.File, hcl.Diagnostics) {
				return &hcl.File{}, hcl.Diagnostics{}
			},
		},
	}

	versions, err := GatherRequiredVersion(conf, exts)
	require.Error(t, err)
	assert.Nil(t, versions)
}

func TestExtractRequiredVersionWithNullValue(t *testing.T) {
	t.Parallel()
	conf := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL with null required_version
	nullHCL := `
terraform {
	required_version = null
}
`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(nullHCL), "test.tf")
	require.False(t, diags.HasErrors())

	versions := extractRequiredVersion(file.Body, conf)
	// Should skip null values
	assert.Empty(t, versions)
}

func TestExtractRequiredVersionWithComplexExpressions(t *testing.T) {
	t.Parallel()
	conf := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Test with expressions that might not be wholly known
	expressionHCL := `
terraform {
	required_version = var.version_constraint
}
`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(expressionHCL), "test.tf")
	require.False(t, diags.HasErrors())

	versions := extractRequiredVersion(file.Body, conf)
	// Should skip unknown values
	assert.Empty(t, versions)
}

func TestExtractRequiredVersionWithInvalidTypeConversion(t *testing.T) {
	t.Parallel()
	conf := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Test with non-string value that can't be converted
	invalidTypeHCL := `
terraform {
	required_version = 123
}
`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(invalidTypeHCL), "test.tf")
	require.False(t, diags.HasErrors())

	versions := extractRequiredVersion(file.Body, conf)
	// Should handle type conversion gracefully
	assert.Empty(t, versions)
}

func TestExtractRequiredVersionWithMultipleBlocks(t *testing.T) {
	t.Parallel()
	conf := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Test with multiple terraform blocks
	multiBlockHCL := `
terraform {
	required_version = ">= 1.0.0"
}

terraform {
	required_version = "~> 1.5.0"
}
`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(multiBlockHCL), "test.tf")
	require.False(t, diags.HasErrors())

	versions := extractRequiredVersion(file.Body, conf)
	assert.Len(t, versions, 2)
	assert.Contains(t, versions, ">= 1.0.0")
	assert.Contains(t, versions, "~> 1.5.0")
}
