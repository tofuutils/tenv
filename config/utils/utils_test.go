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

package configutils_test

import (
	"errors"
	"strconv"
	"testing"

	configutils "github.com/tofuutils/tenv/v4/config/utils"
)

func TestGetenvFunc_Bool(t *testing.T) {
	t.Parallel()

	mockGetenv := func(key string) string {
		envMap := map[string]string{
			"TRUE_VAR":  "true",
			"FALSE_VAR": "false",
			"ONE_VAR":   "1",
			"ZERO_VAR":  "0",
			"EMPTY_VAR": "",
		}

		return envMap[key]
	}

	getenvFunc := configutils.GetenvFunc(mockGetenv)

	tests := []struct {
		name        string
		key         string
		defaultVal  bool
		expected    bool
		expectError bool
	}{
		{"true value", "TRUE_VAR", false, true, false},
		{"false value", "FALSE_VAR", true, false, false},
		{"1 value", "ONE_VAR", false, true, false},
		{"0 value", "ZERO_VAR", true, false, false},
		{"empty value uses default", "EMPTY_VAR", true, true, false},
		{"non-existent uses default", "NON_EXISTENT", false, false, false},
		{"invalid bool uses default", "INVALID_VAR", false, false, false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result, err := getenvFunc.Bool(testCase.defaultVal, testCase.key)
			if testCase.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != testCase.expected {
					t.Errorf("Bool() = %v, want %v", result, testCase.expected)
				}
			}
		})
	}
}

func TestGetenvFunc_BoolFallback(t *testing.T) {
	t.Parallel()

	mockGetenv := func(key string) string {
		envMap := map[string]string{
			"FIRST_VAR":  "true",
			"SECOND_VAR": "false",
			"THIRD_VAR":  "1",
		}

		return envMap[key]
	}

	getenvFunc := configutils.GetenvFunc(mockGetenv)

	tests := []struct {
		name        string
		keys        []string
		defaultVal  bool
		expected    bool
		expectError bool
	}{
		{"first key found", []string{"FIRST_VAR", "SECOND_VAR"}, false, true, false},
		{"second key found", []string{"NON_EXISTENT", "SECOND_VAR"}, false, false, false},
		{"third key found", []string{"NON_EXISTENT1", "NON_EXISTENT2", "THIRD_VAR"}, false, true, false},
		{"no keys found uses default", []string{"NON_EXISTENT1", "NON_EXISTENT2"}, true, true, false},
		{"empty keys slice uses default", []string{}, false, false, false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result, err := getenvFunc.BoolFallback(testCase.defaultVal, testCase.keys...)
			if testCase.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != testCase.expected {
					t.Errorf("BoolFallback() = %v, want %v", result, testCase.expected)
				}
			}
		})
	}
}

func TestGetenvFunc_Fallback(t *testing.T) {
	t.Parallel()

	mockGetenv := func(key string) string {
		envMap := map[string]string{
			"FIRST_VAR":  "first_value",
			"SECOND_VAR": "second_value",
			"EMPTY_VAR":  "",
		}

		return envMap[key]
	}

	getenvFunc := configutils.GetenvFunc(mockGetenv)

	tests := []struct {
		name     string
		keys     []string
		expected string
	}{
		{"first key found", []string{"FIRST_VAR", "SECOND_VAR"}, "first_value"},
		{"second key found", []string{"NON_EXISTENT", "SECOND_VAR"}, "second_value"},
		{"third key found", []string{"NON_EXISTENT1", "NON_EXISTENT2", "FIRST_VAR"}, "first_value"},
		{"no keys found returns empty", []string{"NON_EXISTENT1", "NON_EXISTENT2"}, ""},
		{"empty keys slice returns empty", []string{}, ""},
		{"empty value skipped", []string{"EMPTY_VAR", "FIRST_VAR"}, "first_value"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := getenvFunc.Fallback(testCase.keys...)
			if result != testCase.expected {
				t.Errorf("Fallback() = %q, want %q", result, testCase.expected)
			}
		})
	}
}

func TestGetenvFunc_Present(t *testing.T) {
	t.Parallel()

	mockGetenv := func(key string) string {
		envMap := map[string]string{
			"PRESENT_VAR":    "some_value",
			"EMPTY_VAR":      "",
			"WHITESPACE_VAR": "   ",
		}

		return envMap[key]
	}

	getenvFunc := configutils.GetenvFunc(mockGetenv)

	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"present with value", "PRESENT_VAR", true},
		{"empty value not present", "EMPTY_VAR", false},
		{"whitespace only present", "WHITESPACE_VAR", true},
		{"non-existent not present", "NON_EXISTENT", false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := getenvFunc.Present(testCase.key)
			if result != testCase.expected {
				t.Errorf("Present() = %v, want %v", result, testCase.expected)
			}
		})
	}
}

func TestGetenvFunc_WithDefault(t *testing.T) {
	t.Parallel()

	mockGetenv := func(key string) string {
		envMap := map[string]string{
			"PRESENT_VAR":    "actual_value",
			"EMPTY_VAR":      "",
			"WHITESPACE_VAR": "   ",
		}

		return envMap[key]
	}

	getenvFunc := configutils.GetenvFunc(mockGetenv)

	tests := []struct {
		name         string
		key          string
		defaultValue string
		expected     string
	}{
		{"present value returned", "PRESENT_VAR", "default", "actual_value"},
		{"empty value uses default", "EMPTY_VAR", "default", "default"},
		{"whitespace value returned as-is", "WHITESPACE_VAR", "default", "   "},
		{"non-existent uses default", "NON_EXISTENT", "default", "default"},
		{"empty default", "NON_EXISTENT", "", ""},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := getenvFunc.WithDefault(testCase.defaultValue, testCase.key)
			if result != testCase.expected {
				t.Errorf("WithDefault() = %q, want %q", result, testCase.expected)
			}
		})
	}
}

// Test error cases for strconv.ParseBool.
func TestGetenvFunc_Bool_ErrorHandling(t *testing.T) {
	t.Parallel()

	mockGetenv := func(key string) string {
		return "invalid_bool_value"
	}

	getenvFunc := configutils.GetenvFunc(mockGetenv)

	_, err := getenvFunc.Bool(false, "INVALID_VAR")
	if err == nil {
		t.Errorf("Expected error for invalid bool value")
	}

	var parseError *strconv.NumError
	if !errors.As(err, &parseError) {
		t.Errorf("Expected strconv.NumError, got %T", err)
	}
}
