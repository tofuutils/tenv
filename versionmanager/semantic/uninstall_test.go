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

package semantic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestSelectVersionsToUninstall(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Test versions (should be in descending order)
	versions := []string{"1.5.0", "1.4.0", "1.3.0", "1.2.0", "1.1.0"}

	tests := []struct {
		name        string
		constraint  string
		installPath string
		expected    []string
		expectError bool
	}{
		{
			name:        "all versions",
			constraint:  "all",
			installPath: "/tmp",
			expected:    []string{"1.5.0", "1.4.0", "1.3.0", "1.2.0", "1.1.0"},
			expectError: false,
		},
		{
			name:        "but last (keep latest)",
			constraint:  "but-last",
			installPath: "/tmp",
			expected:    []string{"1.4.0", "1.3.0", "1.2.0", "1.1.0"},
			expectError: false,
		},
		{
			name:        "not used for 30 days",
			constraint:  "not-used-for:30d",
			installPath: "/tmp",
			expected:    []string{"1.5.0", "1.4.0", "1.3.0", "1.2.0", "1.1.0"}, // All versions match (mock returns zero time)
			expectError: false,
		},
		{
			name:        "not used for 2 months",
			constraint:  "not-used-for:2m",
			installPath: "/tmp",
			expected:    []string{"1.5.0", "1.4.0", "1.3.0", "1.2.0", "1.1.0"}, // All versions match (mock returns zero time)
			expectError: false,
		},
		{
			name:        "not used since specific date",
			constraint:  "not-used-since:2024-01-01",
			installPath: "/tmp",
			expected:    []string{"1.5.0", "1.4.0", "1.3.0", "1.2.0", "1.1.0"}, // All versions match (mock returns zero time)
			expectError: false,
		},
		{
			name:        "invalid duration format",
			constraint:  "not-used-for:30x",
			installPath: "/tmp",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid date format",
			constraint:  "not-used-since:invalid-date",
			installPath: "/tmp",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "version constraint - greater than 1.3.0",
			constraint:  "> 1.3.0",
			installPath: "/tmp",
			expected:    []string{"1.5.0", "1.4.0"},
			expectError: false,
		},
		{
			name:        "version constraint - exact version",
			constraint:  "1.3.0",
			installPath: "/tmp",
			expected:    []string{"1.3.0"},
			expectError: false,
		},
		{
			name:        "version constraint - invalid constraint",
			constraint:  "invalid-constraint",
			installPath: "/tmp",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SelectVersionsToUninstall(tt.constraint, tt.installPath, versions, mockConfig)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestFilterStrings(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		predicate func(string) bool
		expected  []string
	}{
		{
			name:  "filter even numbers",
			input: []string{"1", "2", "3", "4", "5"},
			predicate: func(s string) bool {
				// Check if the string represents an even number
				if len(s) == 1 && s[0] >= '2' && s[0] <= '4' && (s[0]-'0')%2 == 0 {
					return true
				}
				return false
			},
			expected: []string{"2", "4"},
		},
		{
			name:  "filter strings starting with 'a'",
			input: []string{"apple", "banana", "avocado", "grape"},
			predicate: func(s string) bool {
				return len(s) > 0 && s[0] == 'a'
			},
			expected: []string{"apple", "avocado"},
		},
		{
			name:  "filter all",
			input: []string{"a", "b", "c"},
			predicate: func(s string) bool {
				return true
			},
			expected: []string{"a", "b", "c"},
		},
		{
			name:  "filter none",
			input: []string{"a", "b", "c"},
			predicate: func(s string) bool {
				return false
			},
			expected: []string{},
		},
		{
			name:  "empty input",
			input: []string{},
			predicate: func(s string) bool {
				return true
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterStrings(tt.input, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPredicateBeforeDate(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	installPath := "/tmp"
	beforeDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Create predicate
	pred := predicateBeforeDate(installPath, beforeDate, mockConfig)

	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "version with zero time",
			version:  "1.0.0",
			expected: true, // Zero time is before any date
		},
		{
			name:     "another version with zero time",
			version:  "2.0.0",
			expected: true, // Zero time is before any date
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pred(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "all", allKey)
	assert.Equal(t, "but-last", butLast)
	assert.Equal(t, "not-used-for:", notUsedForPrefix)
	assert.Equal(t, "not-used-since:", notUsedSincePrefix)
	assert.Equal(t, len(notUsedForPrefix), notUsedForPrefixLen)
	assert.Equal(t, len(notUsedSincePrefix), notUsedSincePrefixLen)
	assert.Equal(t, "unrecognized duration format", errDurationParsing.Error())
}

// Test helper functions that are not exported but can be tested through public APIs
func TestSelectVersionsToUninstallEdgeCases(t *testing.T) {
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	t.Run("nil versions slice", func(t *testing.T) {
		result, err := SelectVersionsToUninstall("all", "/tmp", nil, mockConfig)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("empty versions slice", func(t *testing.T) {
		result, err := SelectVersionsToUninstall("all", "/tmp", []string{}, mockConfig)
		assert.NoError(t, err)
		assert.Equal(t, []string{}, result)
	})

	t.Run("single version with but-last", func(t *testing.T) {
		result, err := SelectVersionsToUninstall("but-last", "/tmp", []string{"1.0.0"}, mockConfig)
		assert.NoError(t, err)
		assert.Equal(t, []string{}, result)
	})

	t.Run("but-last with empty versions", func(t *testing.T) {
		result, err := SelectVersionsToUninstall("but-last", "/tmp", []string{}, mockConfig)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
