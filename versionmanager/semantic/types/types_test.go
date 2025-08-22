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
	"github.com/tofuutils/tenv/v4/config"
)

// mockDisplayer implements loghelper.Displayer for testing
type mockDisplayer struct {
	displayedMessages []string
}

func (m *mockDisplayer) Display(msg string) {
	m.displayedMessages = append(m.displayedMessages, msg)
}

func (m *mockDisplayer) Log(hclog.Level, string, ...interface{}) {}
func (m *mockDisplayer) IsDebug() bool                           { return false }
func (m *mockDisplayer) IsTrace() bool                           { return false }
func (m *mockDisplayer) Flush(bool)                              {}

func TestDisplayDetectionInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		version        string
		source         string
		expectedMsg    string
		expectedResult string
	}{
		{
			name:           "valid version and source",
			version:        "1.0.0",
			source:         "test.txt",
			expectedMsg:    "Resolved version from test.txt : 1.0.0",
			expectedResult: "1.0.0",
		},
		{
			name:           "empty version",
			version:        "",
			source:         "test.txt",
			expectedMsg:    "Resolved version from test.txt : ",
			expectedResult: "",
		},
		{
			name:           "empty source",
			version:        "1.0.0",
			source:         "",
			expectedMsg:    "Resolved version from  : 1.0.0",
			expectedResult: "1.0.0",
		},
		{
			name:           "both empty",
			version:        "",
			source:         "",
			expectedMsg:    "Resolved version from  : ",
			expectedResult: "",
		},
		{
			name:           "version with spaces",
			version:        " 1.0.0 ",
			source:         "test.txt",
			expectedMsg:    "Resolved version from test.txt :  1.0.0 ",
			expectedResult: " 1.0.0 ",
		},
		{
			name:           "source with spaces",
			version:        "1.0.0",
			source:         " test.txt ",
			expectedMsg:    "Resolved version from  test.txt  : 1.0.0",
			expectedResult: "1.0.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock displayer
			displayer := &mockDisplayer{}

			// Run test
			result := DisplayDetectionInfo(displayer, tt.version, tt.source)

			// Check result
			if result != tt.expectedResult {
				t.Errorf("expected result %s but got %s", tt.expectedResult, result)
			}

			// Check displayed message
			if len(displayer.displayedMessages) != 1 {
				t.Errorf("expected 1 displayed message but got %d", len(displayer.displayedMessages))
			} else if displayer.displayedMessages[0] != tt.expectedMsg {
				t.Errorf("expected message %s but got %s", tt.expectedMsg, displayer.displayedMessages[0])
			}
		})
	}
}

func TestPredicateInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		predicate     func(string) bool
		reverseOrder  bool
		testValue     string
		expectedValue bool
	}{
		{
			name: "always true predicate",
			predicate: func(s string) bool {
				return true
			},
			reverseOrder:  false,
			testValue:     "test",
			expectedValue: true,
		},
		{
			name: "always false predicate",
			predicate: func(s string) bool {
				return false
			},
			reverseOrder:  false,
			testValue:     "test",
			expectedValue: false,
		},
		{
			name: "length check predicate",
			predicate: func(s string) bool {
				return len(s) > 3
			},
			reverseOrder:  false,
			testValue:     "test",
			expectedValue: true,
		},
		{
			name: "length check predicate with short string",
			predicate: func(s string) bool {
				return len(s) > 3
			},
			reverseOrder:  false,
			testValue:     "hi",
			expectedValue: false,
		},
		{
			name: "reverse order with true predicate",
			predicate: func(s string) bool {
				return true
			},
			reverseOrder:  true,
			testValue:     "test",
			expectedValue: true,
		},
		{
			name: "reverse order with false predicate",
			predicate: func(s string) bool {
				return false
			},
			reverseOrder:  true,
			testValue:     "test",
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create PredicateInfo
			info := PredicateInfo{
				Predicate:    tt.predicate,
				ReverseOrder: tt.reverseOrder,
			}

			// Run test
			result := info.Predicate(tt.testValue)

			// Check result
			if result != tt.expectedValue {
				t.Errorf("expected %v but got %v", tt.expectedValue, result)
			}
		})
	}
}

func TestVersionFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		fileName       string
		parser         func(string, *config.Config) (string, error)
		expectedName   string
		expectedParser bool
	}{
		{
			name:     "valid file with parser",
			fileName: "test.txt",
			parser: func(s string, c *config.Config) (string, error) {
				return "1.0.0", nil
			},
			expectedName:   "test.txt",
			expectedParser: true,
		},
		{
			name:           "valid file without parser",
			fileName:       "test.txt",
			parser:         nil,
			expectedName:   "test.txt",
			expectedParser: false,
		},
		{
			name:           "empty file name",
			fileName:       "",
			parser:         nil,
			expectedName:   "",
			expectedParser: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create VersionFile
			file := VersionFile{
				Name:   tt.fileName,
				Parser: tt.parser,
			}

			// Check name
			if file.Name != tt.expectedName {
				t.Errorf("expected name %s but got %s", tt.expectedName, file.Name)
			}

			// Check parser
			if (file.Parser != nil) != tt.expectedParser {
				t.Errorf("expected parser %v but got %v", tt.expectedParser, file.Parser != nil)
			}

			// If parser exists, test it
			if file.Parser != nil {
				result, err := file.Parser("test.txt", nil)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != "1.0.0" {
					t.Errorf("expected result 1.0.0 but got %s", result)
				}
			}
		})
	}
}

// Mock implementation of ConstraintInfo
type mockConstraintInfo struct {
	constraint string
}

var _ ConstraintInfo = (*mockConstraintInfo)(nil) // Ensure mock implements the interface

// Method implementation
func (m *mockConstraintInfo) ReadDefaultConstraint() string {
	return m.constraint
}

func TestConstraintInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		constraint     string
		expectedResult string
	}{
		{
			name:           "valid constraint",
			constraint:     ">= 1.0.0",
			expectedResult: ">= 1.0.0",
		},
		{
			name:           "empty constraint",
			constraint:     "",
			expectedResult: "",
		},
		{
			name:           "constraint with spaces",
			constraint:     " >= 1.0.0 ",
			expectedResult: " >= 1.0.0 ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock
			mock := &mockConstraintInfo{
				constraint: tt.constraint,
			}

			// Run test
			result := mock.ReadDefaultConstraint()

			// Check result
			if result != tt.expectedResult {
				t.Errorf("expected result %s but got %s", tt.expectedResult, result)
			}
		})
	}
}
