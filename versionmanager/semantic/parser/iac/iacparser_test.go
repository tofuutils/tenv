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
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/tofuutils/tenv/v4/config"
)

// mockDisplayer implements loghelper.Displayer for testing
type mockDisplayer struct{}

func (m *mockDisplayer) Display(string)                          {}
func (m *mockDisplayer) Log(hclog.Level, string, ...interface{}) {}
func (m *mockDisplayer) IsDebug() bool                           { return false }
func (m *mockDisplayer) IsTrace() bool                           { return false }
func (m *mockDisplayer) Flush(bool)                              {}

// mockParser is a helper function that parses HCL content for testing
func mockParser(path string) (*hcl.File, hcl.Diagnostics) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read file",
			Detail:   err.Error(),
		}}
	}

	return hclsyntax.ParseConfig(content, path, hcl.Pos{Line: 1, Column: 1})
}

func TestGatherRequiredVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		files          map[string]string
		exts           []ExtDescription
		expectedResult []string
		expectError    bool
	}{
		{
			name: "single file with version",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = "~> 1.0.0"
				}`,
			},
			exts: []ExtDescription{
				{
					Value:  ".tf",
					Parser: mockParser,
				},
			},
			expectedResult: []string{"~> 1.0.0"},
		},
		{
			name: "multiple files with versions",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = "~> 1.0.0"
				}`,
				"other.tf": `terraform {
					required_version = ">= 1.2.0"
				}`,
			},
			exts: []ExtDescription{
				{
					Value:  ".tf",
					Parser: mockParser,
				},
			},
			expectedResult: []string{"~> 1.0.0", ">= 1.2.0"},
		},
		{
			name: "file without version",
			files: map[string]string{
				"main.tf": `terraform {
					backend "local" {}
				}`,
			},
			exts: []ExtDescription{
				{
					Value:  ".tf",
					Parser: mockParser,
				},
			},
			expectedResult: nil,
		},
		{
			name:           "no extensions",
			files:          nil,
			exts:           nil,
			expectedResult: nil,
		},
		{
			name: "invalid HCL",
			files: map[string]string{
				"main.tf": `terraform {
					required_version =
				}`,
			},
			exts: []ExtDescription{
				{
					Value:  ".tf",
					Parser: mockParser,
				},
			},
			expectError: true,
		},
		{
			name: "multiple file extensions",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = "~> 1.0.0"
				}`,
				"other.hcl": `terraform {
					required_version = ">= 1.2.0"
				}`,
			},
			exts: []ExtDescription{
				{
					Value:  ".tf",
					Parser: mockParser,
				},
				{
					Value:  ".hcl",
					Parser: mockParser,
				},
			},
			expectedResult: []string{"~> 1.0.0", ">= 1.2.0"},
		},
		{
			name: "non-existent file",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = "~> 1.0.0"
				}`,
			},
			exts: []ExtDescription{
				{
					Value: ".tf",
					Parser: func(path string) (*hcl.File, hcl.Diagnostics) {
						return nil, hcl.Diagnostics{{
							Severity: hcl.DiagError,
							Summary:  "File not found",
						}}
					},
				},
			},
			expectError: true,
		},
		{
			name: "mixed valid and invalid files",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = "~> 1.0.0"
				}`,
				"invalid.tf": `terraform {
					required_version =
				}`,
			},
			exts: []ExtDescription{
				{
					Value:  ".tf",
					Parser: mockParser,
				},
			},
			expectError: true,
		},
		{
			name: "empty terraform block",
			files: map[string]string{
				"main.tf": `terraform {}`,
			},
			exts: []ExtDescription{
				{
					Value:  ".tf",
					Parser: mockParser,
				},
			},
			expectedResult: nil,
		},
		{
			name: "non-string version value",
			files: map[string]string{
				"main.tf": `terraform {
					required_version = 1.0
				}`,
			},
			exts: []ExtDescription{
				{
					Value:  ".tf",
					Parser: mockParser,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()

			// Create test files
			for name, content := range tt.files {
				err := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0600)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Create config
			conf := &config.Config{
				WorkPath:  tempDir,
				Displayer: &mockDisplayer{},
			}

			// Run test
			result, err := GatherRequiredVersion(conf, tt.exts)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if len(result) != len(tt.expectedResult) {
				t.Errorf("expected %d results but got %d", len(tt.expectedResult), len(result))
			}

			for i, v := range result {
				if v != tt.expectedResult[i] {
					t.Errorf("expected %s but got %s at index %d", tt.expectedResult[i], v, i)
				}
			}
		})
	}
}

func TestFilterExts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		fileExts       int
		exts           []ExtDescription
		expectedResult ExtDescription
	}{
		{
			name:     "single extension",
			fileExts: 1,
			exts: []ExtDescription{
				{Value: ".tf"},
				{Value: ".hcl"},
			},
			expectedResult: ExtDescription{Value: ".tf"},
		},
		{
			name:     "second extension",
			fileExts: 2,
			exts: []ExtDescription{
				{Value: ".tf"},
				{Value: ".hcl"},
			},
			expectedResult: ExtDescription{Value: ".hcl"},
		},
		{
			name:     "multiple bits set - first wins",
			fileExts: 3, // 0b11
			exts: []ExtDescription{
				{Value: ".tf"},
				{Value: ".hcl"},
			},
			expectedResult: ExtDescription{Value: ".tf"},
		},
		{
			name:     "high bit set",
			fileExts: 8, // 0b1000
			exts: []ExtDescription{
				{Value: ".tf"},
				{Value: ".hcl"},
				{Value: ".tfvars"},
				{Value: ".auto.tfvars"},
			},
			expectedResult: ExtDescription{Value: ".auto.tfvars"},
		},
		{
			name:     "all bits set",
			fileExts: 15, // 0b1111
			exts: []ExtDescription{
				{Value: ".tf"},
				{Value: ".hcl"},
				{Value: ".tfvars"},
				{Value: ".auto.tfvars"},
			},
			expectedResult: ExtDescription{Value: ".tf"},
		},
		{
			name:     "single extension in list",
			fileExts: 1,
			exts: []ExtDescription{
				{Value: ".tf"},
			},
			expectedResult: ExtDescription{Value: ".tf"},
		},
		{
			name:     "no matching bits",
			fileExts: 16, // 0b10000
			exts: []ExtDescription{
				{Value: ".tf"},
				{Value: ".hcl"},
			},
			expectedResult: ExtDescription{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := filterExts(tt.fileExts, tt.exts)
			if result.Value != tt.expectedResult.Value {
				t.Errorf("expected %s but got %s", tt.expectedResult.Value, result.Value)
			}
		})
	}
}

func TestExtractRequiredVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		hclContent     string
		expectedResult []string
		expectError    bool
	}{
		{
			name: "valid required version",
			hclContent: `terraform {
				required_version = "~> 1.0.0"
			}`,
			expectedResult: []string{"~> 1.0.0"},
		},
		{
			name: "multiple terraform blocks",
			hclContent: `terraform {
				required_version = "~> 1.0.0"
			}
			terraform {
				required_version = ">= 1.2.0"
			}`,
			expectedResult: []string{"~> 1.0.0", ">= 1.2.0"},
		},
		{
			name: "no required version",
			hclContent: `terraform {
				backend "local" {}
			}`,
			expectedResult: nil,
		},
		{
			name:           "empty content",
			hclContent:     "",
			expectedResult: nil,
		},
		{
			name: "complex terraform block",
			hclContent: `terraform {
				required_version = "~> 1.0.0"
				required_providers {
					aws = {
						source  = "hashicorp/aws"
						version = "~> 4.0"
					}
				}
				backend "s3" {
					bucket = "mybucket"
					key    = "path/to/my/key"
					region = "us-east-1"
				}
			}`,
			expectedResult: []string{"~> 1.0.0"},
		},
		{
			name: "invalid version format",
			hclContent: `terraform {
				required_version = 123
			}`,
			expectError: true,
		},
		{
			name: "multiple blocks with mixed content",
			hclContent: `terraform {
				required_version = "~> 1.0.0"
			}
			resource "aws_instance" "example" {
				ami = "ami-123456"
			}
			terraform {
				required_version = ">= 1.2.0"
				backend "local" {}
			}`,
			expectedResult: []string{"~> 1.0.0", ">= 1.2.0"},
		},
		{
			name: "terraform block with comments",
			hclContent: `# Main terraform configuration
			terraform { # Start block
				# Version requirement
				required_version = "~> 1.0.0" # Specify version constraint
			} # End block`,
			expectedResult: []string{"~> 1.0.0"},
		},
		{
			name: "terraform block with heredoc",
			hclContent: `terraform {
				required_version = <<EOF
~> 1.0.0
EOF
			}`,
			expectedResult: []string{"~> 1.0.0"},
		},
		{
			name: "invalid block structure",
			hclContent: `terraform {
				required_version = "~> 1.0.0"
				{
					invalid = true
				}
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Parse HCL content
			file, diags := hclsyntax.ParseConfig([]byte(tt.hclContent), "test.tf", hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				t.Fatal(diags)
			}

			// Create config
			conf := &config.Config{
				Displayer: &mockDisplayer{},
			}

			result := extractRequiredVersion(file.Body, conf)
			if len(result) != len(tt.expectedResult) {
				t.Errorf("expected %d results but got %d", len(tt.expectedResult), len(result))
			}

			for i, v := range result {
				if v != tt.expectedResult[i] {
					t.Errorf("expected %s but got %s at index %d", tt.expectedResult[i], v, i)
				}
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	// Setup temp directory
	tempDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"main.tf": `terraform {
			required_version = "~> 1.0.0"
		}`,
		"other.tf": `terraform {
			required_version = ">= 1.2.0"
		}`,
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0600)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create config
	conf := &config.Config{
		WorkPath:  tempDir,
		Displayer: &mockDisplayer{},
	}

	// Create extensions
	exts := []ExtDescription{
		{
			Value:  ".tf",
			Parser: mockParser,
		},
	}

	// Number of concurrent goroutines
	numGoroutines := 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Run concurrent tests
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			result, err := GatherRequiredVersion(conf, exts)
			if err != nil {
				t.Error(err)
				return
			}
			if len(result) != 2 {
				t.Errorf("expected 2 results but got %d", len(result))
			}
		}()
	}

	wg.Wait()
}

func TestFileErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(string) error
		expectError bool
	}{
		{
			name: "unreadable directory",
			setup: func(dir string) error {
				return os.Chmod(dir, 0000)
			},
			expectError: true,
		},
		{
			name: "unreadable file",
			setup: func(dir string) error {
				filePath := filepath.Join(dir, "main.tf")
				if err := os.WriteFile(filePath, []byte(`terraform {}`), 0600); err != nil {
					return err
				}
				return os.Chmod(filePath, 0000)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()

			// Apply setup
			if err := tt.setup(tempDir); err != nil {
				t.Fatal(err)
			}

			// Create config
			conf := &config.Config{
				WorkPath:  tempDir,
				Displayer: &mockDisplayer{},
			}

			// Run test
			_, err := GatherRequiredVersion(conf, []ExtDescription{{
				Value:  ".tf",
				Parser: mockParser,
			}})

			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

func TestParserErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		parser      func(string) (*hcl.File, hcl.Diagnostics)
		expectError bool
	}{
		{
			name: "nil file no error",
			parser: func(string) (*hcl.File, hcl.Diagnostics) {
				return nil, nil
			},
			expectError: false,
		},
		{
			name: "nil file with error",
			parser: func(string) (*hcl.File, hcl.Diagnostics) {
				return nil, hcl.Diagnostics{{
					Severity: hcl.DiagError,
					Summary:  "Test error",
				}}
			},
			expectError: true,
		},
		{
			name: "parser panic",
			parser: func(string) (*hcl.File, hcl.Diagnostics) {
				panic("parser error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup temp directory
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "main.tf")
			if err := os.WriteFile(filePath, []byte(`terraform {}`), 0600); err != nil {
				t.Fatal(err)
			}

			// Create config
			conf := &config.Config{
				WorkPath:  tempDir,
				Displayer: &mockDisplayer{},
			}

			// Run test with recovery for panics
			var err error
			func() {
				defer func() {
					if r := recover(); r != nil && !tt.expectError {
						t.Errorf("unexpected panic: %v", r)
					}
				}()

				_, err = GatherRequiredVersion(conf, []ExtDescription{{
					Value:  ".tf",
					Parser: tt.parser,
				}})
			}()

			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}
