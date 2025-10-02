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

package cmdconst

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "AgnosticName constant",
			constant: AgnosticName,
			expected: "tf",
		},
		{
			name:     "AtmosName constant",
			constant: AtmosName,
			expected: "atmos",
		},
		{
			name:     "TenvName constant",
			constant: TenvName,
			expected: "tenv",
		},
		{
			name:     "TerraformName constant",
			constant: TerraformName,
			expected: "terraform",
		},
		{
			name:     "TerragruntName constant",
			constant: TerragruntName,
			expected: "terragrunt",
		},
		{
			name:     "TerramateName constant",
			constant: TerramateName,
			expected: "terramate",
		},
		{
			name:     "TofuName constant",
			constant: TofuName,
			expected: "tofu",
		},
		{
			name:     "OpentofuName constant",
			constant: OpentofuName,
			expected: "opentofu",
		},
		{
			name:     "CallSubCmd constant",
			constant: CallSubCmd,
			expected: "call",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

func TestExitCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant int
		expected int
	}{
		{
			name:     "BasicErrorExitCode constant",
			constant: BasicErrorExitCode,
			expected: 1,
		},
		{
			name:     "EarlyErrorExitCode constant",
			constant: EarlyErrorExitCode,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}
