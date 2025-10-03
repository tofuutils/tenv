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

package types

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

// Mock implementation of ConstraintInfo interface for testing.
type mockConstraintInfo struct {
	constraint string
}

func (m *mockConstraintInfo) ReadDefaultConstraint() string {
	return m.constraint
}

func TestConstraintInfoInterface(t *testing.T) {
	t.Parallel()

	// Test that the interface is properly defined
	var _ ConstraintInfo = (*mockConstraintInfo)(nil)

	// Test with different constraint values
	tests := []struct {
		name       string
		constraint string
		expected   string
	}{
		{
			name:       "simple version constraint",
			constraint: ">= 1.0.0",
			expected:   ">= 1.0.0",
		},
		{
			name:       "exact version",
			constraint: "1.2.3",
			expected:   "1.2.3",
		},
		{
			name:       "empty constraint",
			constraint: "",
			expected:   "",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockConstraintInfo{constraint: testCase.constraint}
			result := mock.ReadDefaultConstraint()
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestDisplayDetectionInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		version  string
		source   string
		expected string
	}{
		{
			name:     "terraform version from .terraform-version",
			version:  "1.5.0",
			source:   ".terraform-version",
			expected: "1.5.0",
		},
		{
			name:     "tofu version from .opentofu-version",
			version:  "1.6.0",
			source:   ".opentofu-version",
			expected: "1.6.0",
		},
		{
			name:     "terragrunt version from .terragrunt-version",
			version:  "0.45.0",
			source:   ".terragrunt-version",
			expected: "0.45.0",
		},
		{
			name:     "empty version",
			version:  "",
			source:   ".tool-versions",
			expected: "",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock displayer that captures the output
			var capturedMessage string
			mockDisplayer := &mockDisplayer{capture: &capturedMessage}

			// Test DisplayDetectionInfo
			result := DisplayDetectionInfo(mockDisplayer, testCase.version, testCase.source)

			// Verify the returned version
			assert.Equal(t, testCase.expected, result)

			// Verify the display message was called
			assert.NotEmpty(t, capturedMessage)
			assert.Contains(t, capturedMessage, "Resolved version from")
			assert.Contains(t, capturedMessage, testCase.source)
			assert.Contains(t, capturedMessage, testCase.version)
		})
	}
}

func TestPredicateInfo(t *testing.T) {
	t.Parallel()

	// Test PredicateInfo struct
	tests := []struct {
		name         string
		predicate    func(string) bool
		reverseOrder bool
		testVersion  string
		expected     bool
	}{
		{
			name: "always true predicate",
			predicate: func(s string) bool {
				return true
			},
			reverseOrder: false,
			testVersion:  "1.0.0",
			expected:     true,
		},
		{
			name: "always false predicate",
			predicate: func(s string) bool {
				return false
			},
			reverseOrder: true,
			testVersion:  "1.0.0",
			expected:     false,
		},
		{
			name: "version starts with 1 predicate",
			predicate: func(s string) bool {
				return len(s) > 0 && s[0] == '1'
			},
			reverseOrder: false,
			testVersion:  "1.5.0",
			expected:     true,
		},
		{
			name: "version starts with 2 predicate",
			predicate: func(s string) bool {
				return len(s) > 0 && s[0] == '2'
			},
			reverseOrder: true,
			testVersion:  "1.5.0",
			expected:     false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			predicateInfo := PredicateInfo{
				Predicate:    testCase.predicate,
				ReverseOrder: testCase.reverseOrder,
			}

			// Test that the predicate works
			result := predicateInfo.Predicate(testCase.testVersion)
			assert.Equal(t, testCase.expected, result)

			// Test that the reverse order flag is set correctly
			assert.Equal(t, testCase.reverseOrder, predicateInfo.ReverseOrder)
		})
	}
}

func TestVersionFile(t *testing.T) {
	t.Parallel()

	// Test VersionFile struct
	tests := []struct {
		name     string
		fileName string
		parser   func(filePath string, conf *config.Config) (string, error)
	}{
		{
			name:     ".terraform-version file",
			fileName: ".terraform-version",
			parser: func(filePath string, conf *config.Config) (string, error) {
				return "1.5.0", nil
			},
		},
		{
			name:     ".opentofu-version file",
			fileName: ".opentofu-version",
			parser: func(filePath string, conf *config.Config) (string, error) {
				return "1.6.0", nil
			},
		},
		{
			name:     ".tool-versions file",
			fileName: ".tool-versions",
			parser: func(filePath string, conf *config.Config) (string, error) {
				return "1.4.0", nil
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			versionFile := VersionFile{
				Name:   testCase.fileName,
				Parser: testCase.parser,
			}

			// Test that the name is set correctly
			assert.Equal(t, testCase.fileName, versionFile.Name)

			// Test that the parser is set correctly
			assert.NotNil(t, versionFile.Parser)

			// Test that the parser function works
			mockConfig := &config.Config{
				Displayer: loghelper.InertDisplayer,
			}
			version, err := versionFile.Parser("test-file", mockConfig)
			require.NoError(t, err)
			assert.NotEmpty(t, version)
		})
	}
}

func TestVersionFileWithNilParser(t *testing.T) {
	t.Parallel()

	// Test VersionFile with nil parser
	versionFile := VersionFile{
		Name:   ".test-version",
		Parser: nil,
	}

	// Should not panic when accessing fields
	assert.Equal(t, ".test-version", versionFile.Name)
	assert.Nil(t, versionFile.Parser)
}

// mockDisplayer is a mock implementation of loghelper.Displayer for testing.
type mockDisplayer struct {
	capture *string
}

func (m *mockDisplayer) Display(message string) {
	if m.capture != nil {
		*m.capture = message
	}
}

func (m *mockDisplayer) DisplayVerbose(message string) {
	if m.capture != nil {
		*m.capture = message
	}
}

func (m *mockDisplayer) IsDebug() bool {
	return false
}

func (m *mockDisplayer) Log(level hclog.Level, message string, args ...any) {
	if m.capture != nil {
		*m.capture = message
	}
}

func (m *mockDisplayer) Flush(force bool) {
	// No-op for testing
}
