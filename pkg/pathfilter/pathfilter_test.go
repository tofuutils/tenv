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

package pathfilter_test

import (
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/pathfilter"
)

func TestNameEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		targetName string
		testCases  []struct {
			path     string
			expected bool
		}
	}{
		{
			name:       "simple filename match",
			targetName: "test.txt",
			testCases: []struct {
				path     string
				expected bool
			}{
				{"test.txt", true},
				{"other.txt", false},
				{"path/test.txt", true},
				{"path/other.txt", false},
				{"deep/nested/path/test.txt", true},
				{"deep/nested/path/other.txt", false},
			},
		},
		{
			name:       "filename with extension",
			targetName: "main.go",
			testCases: []struct {
				path     string
				expected bool
			}{
				{"main.go", true},
				{"utils.go", false},
				{"src/main.go", true},
				{"src/utils.go", false},
				{"cmd/app/main.go", true},
				{"cmd/app/utils.go", false},
			},
		},
		{
			name:       "filename without extension",
			targetName: "README",
			testCases: []struct {
				path     string
				expected bool
			}{
				{"README", true},
				{"readme", false},
				{"docs/README", true},
				{"docs/readme", false},
				{"README.md", false},
				{"docs/README.md", false},
			},
		},
		{
			name:       "hidden file",
			targetName: ".gitignore",
			testCases: []struct {
				path     string
				expected bool
			}{
				{".gitignore", true},
				{"gitignore", false},
				{".gitignore.bak", false},
				{"path/.gitignore", true},
				{"path/gitignore", false},
			},
		},
		{
			name:       "empty target name",
			targetName: "",
			testCases: []struct {
				path     string
				expected bool
			}{
				{"", true},
				{"any", false},
				{"path/", true},
				{"path/file", false},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			filter := pathfilter.NameEqual(testCase.targetName)

			for _, testCaseInner := range testCase.testCases {
				result := filter(testCaseInner.path)
				if result != testCaseInner.expected {
					t.Errorf("NameEqual(%q)(%q) = %v, want %v", testCase.targetName, testCaseInner.path, result, testCaseInner.expected)
				}
			}
		})
	}
}

func TestNameEqual_UnixPaths(t *testing.T) {
	t.Parallel()

	filter := pathfilter.NameEqual("test.txt")

	tests := []struct {
		path     string
		expected bool
	}{
		{"test.txt", true},
		{"/test.txt", true},
		{"/path/test.txt", true},
		{"/path/to/test.txt", true},
		{"/path/to/other.txt", false},
		{"path/test.txt", true},
		{"path/to/test.txt", true},
	}

	for _, tt := range tests {
		result := filter(tt.path)
		if result != tt.expected {
			t.Errorf("NameEqual(\"test.txt\")(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestNameEqual_WindowsPaths(t *testing.T) {
	t.Parallel()

	filter := pathfilter.NameEqual("test.txt")

	tests := []struct {
		path     string
		expected bool
	}{
		{`test.txt`, true},
		{`C:\test.txt`, true},
		{`C:\path\test.txt`, true},
		{`C:\path\to\test.txt`, true},
		{`C:\path\to\other.txt`, false},
		{`path\test.txt`, true},
		{`path\to\test.txt`, true},
	}

	for _, tt := range tests {
		result := filter(tt.path)
		if result != tt.expected {
			t.Errorf("NameEqual(\"test.txt\")(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestNameEqual_MixedSeparators(t *testing.T) {
	t.Parallel()

	filter := pathfilter.NameEqual("test.txt")

	tests := []struct {
		path     string
		expected bool
	}{
		{`C:\unix\mixed\test.txt`, true}, // Last separator is \, so filename is test.txt
		{`C:\unix\mixed\other.txt`, false},
		{`/windows/mixed/test.txt`, true}, // Last separator is /, so filename is test.txt
		{`/windows/mixed/other.txt`, false},
	}

	for _, tt := range tests {
		result := filter(tt.path)
		if result != tt.expected {
			t.Errorf("NameEqual(\"test.txt\")(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestNameEqual_EdgeCases(t *testing.T) {
	t.Parallel()

	filter := pathfilter.NameEqual("test")

	tests := []struct {
		path     string
		expected bool
	}{
		{"test", true},
		{"test.txt", false},
		{"test.bak", false},
		{"mytest", false},
		{"test", true},
		{"/test", true},
		{`C:\test`, true},
		{"path/test", true},
		{"path/test/extra", false},
	}

	for _, tt := range tests {
		result := filter(tt.path)
		if result != tt.expected {
			t.Errorf("NameEqual(\"test\")(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestNameEqual_Consistency(t *testing.T) {
	t.Parallel()

	// Test that the same filter gives consistent results
	filter := pathfilter.NameEqual("consistent.txt")

	testPath := "some/path/consistent.txt"
	result1 := filter(testPath)
	result2 := filter(testPath)

	if result1 != result2 {
		t.Errorf("Filter gave inconsistent results for same input: %v != %v", result1, result2)
	}
}
