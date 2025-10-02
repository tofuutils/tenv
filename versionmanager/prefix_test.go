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

package versionmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvPrefixVersion(t *testing.T) {
	tests := []struct {
		name     string
		prefix   EnvPrefix
		expected string
	}{
		{
			name:     "TF prefix version",
			prefix:   "TF",
			expected: "TFVERSION",
		},
		{
			name:     "TOFU prefix version",
			prefix:   "TOFU",
			expected: "TOFUVERSION",
		},
		{
			name:     "ATMOS prefix version",
			prefix:   "ATMOS",
			expected: "ATMOSVERSION",
		},
		{
			name:     "TERRAGRUNT prefix version",
			prefix:   "TG",
			expected: "TGVERSION",
		},
		{
			name:     "empty prefix version",
			prefix:   "",
			expected: "VERSION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.prefix.Version()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvPrefixConstraint(t *testing.T) {
	tests := []struct {
		name     string
		prefix   EnvPrefix
		expected string
	}{
		{
			name:     "TF prefix constraint",
			prefix:   "TF",
			expected: "TFDEFAULT_CONSTRAINT",
		},
		{
			name:     "TOFU prefix constraint",
			prefix:   "TOFU",
			expected: "TOFUDEFAULT_CONSTRAINT",
		},
		{
			name:     "ATMOS prefix constraint",
			prefix:   "ATMOS",
			expected: "ATMOSDEFAULT_CONSTRAINT",
		},
		{
			name:     "TERRAGRUNT prefix constraint",
			prefix:   "TG",
			expected: "TGDEFAULT_CONSTRAINT",
		},
		{
			name:     "empty prefix constraint",
			prefix:   "",
			expected: "DEFAULT_CONSTRAINT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.prefix.constraint()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvPrefixDefaultVersion(t *testing.T) {
	tests := []struct {
		name     string
		prefix   EnvPrefix
		expected string
	}{
		{
			name:     "TF prefix default version",
			prefix:   "TF",
			expected: "TFDEFAULT_VERSION",
		},
		{
			name:     "TOFU prefix default version",
			prefix:   "TOFU",
			expected: "TOFUDEFAULT_VERSION",
		},
		{
			name:     "ATMOS prefix default version",
			prefix:   "ATMOS",
			expected: "ATMOSDEFAULT_VERSION",
		},
		{
			name:     "TERRAGRUNT prefix default version",
			prefix:   "TG",
			expected: "TGDEFAULT_VERSION",
		},
		{
			name:     "empty prefix default version",
			prefix:   "",
			expected: "DEFAULT_VERSION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.prefix.defaultVersion()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvPrefixStringConversion(t *testing.T) {
	// Test that EnvPrefix can be converted to string
	prefix := EnvPrefix("TEST")
	assert.Equal(t, "TEST", string(prefix))

	// Test that string can be converted to EnvPrefix
	prefix2 := EnvPrefix("TEST2")
	assert.Equal(t, EnvPrefix("TEST2"), prefix2)
}

func TestEnvPrefixMethodChaining(t *testing.T) {
	// Test that methods can be chained and work correctly
	prefix := EnvPrefix("CHAIN")

	version := prefix.Version()
	constraint := prefix.constraint()
	defaultVersion := prefix.defaultVersion()

	assert.Equal(t, "CHAINVERSION", version)
	assert.Equal(t, "CHAINDEFAULT_CONSTRAINT", constraint)
	assert.Equal(t, "CHAINDEFAULT_VERSION", defaultVersion)
}
