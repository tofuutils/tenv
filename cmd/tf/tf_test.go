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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
)

func TestTfMainPackage(t *testing.T) {
	// Test that the required constants are accessible
	assert.Equal(t, "tf", cmdconst.AgnosticName)

	// Test that the package can be imported without issues
	// This is a compile-time test that ensures the imports work correctly
	t.Logf("Successfully imported cmdconst.AgnosticName: %s", cmdconst.AgnosticName)
}

func TestTfConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "AgnosticName constant",
			value:    cmdconst.AgnosticName,
			expected: "tf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.value)
		})
	}
}
