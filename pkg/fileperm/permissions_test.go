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

package fileperm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant int
		expected int
	}{
		{
			name:     "RW constant (read-write)",
			constant: RW,
			expected: 0o644,
		},
		{
			name:     "RWE constant (read-write-execute)",
			constant: RWE,
			expected: 0o755,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}
