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

package archname_test

import (
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/archname"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"amd64 to x86_64", "amd64", "x86_64"},
		{"386 to i386", "386", "i386"},
		{"x86_64 unchanged", "x86_64", "x86_64"},
		{"i386 unchanged", "i386", "i386"},
		{"arm64 unchanged", "arm64", "arm64"},
		{"ppc64le unchanged", "ppc64le", "ppc64le"},
		{"s390x unchanged", "s390x", "s390x"},
		{"empty string unchanged", "", ""},
		{"unknown arch unchanged", "unknown", "unknown"},
		{"riscv64 unchanged", "riscv64", "riscv64"},
		{"loong64 unchanged", "loong64", "loong64"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := archname.Convert(testCase.input)
			if result != testCase.expected {
				t.Errorf("Convert(%q) = %q, want %q", testCase.input, result, testCase.expected)
			}
		})
	}
}

func TestConvertConsistency(t *testing.T) {
	t.Parallel()

	// Test that converting an already converted value doesn't change it
	original := "amd64"
	converted := archname.Convert(original)
	doubleConverted := archname.Convert(converted)

	if converted != doubleConverted {
		t.Errorf("Convert(%q) = %q, but Convert(%q) = %q", original, converted, converted, doubleConverted)
	}
}

func TestConvertCaseSensitivity(t *testing.T) {
	t.Parallel()

	// Test that the conversion is case-sensitive (should not match)
	upperCase := "AMD64"
	result := archname.Convert(upperCase)

	if result != upperCase {
		t.Errorf("Convert(%q) = %q, expected unchanged %q", upperCase, result, upperCase)
	}
}

func TestConvertMapIntegrity(t *testing.T) {
	t.Parallel()

	// Test that known conversions work correctly
	knownConversions := []struct {
		input    string
		expected string
	}{
		{"amd64", "x86_64"},
		{"386", "i386"},
	}

	for _, conv := range knownConversions {
		result := archname.Convert(conv.input)
		if result != conv.expected {
			t.Errorf("Convert(%q) = %q, want %q", conv.input, result, conv.expected)
		}
	}
}
