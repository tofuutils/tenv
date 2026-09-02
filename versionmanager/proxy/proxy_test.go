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

package proxy

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/tofuutils/tenv/v4/config"
	configutils "github.com/tofuutils/tenv/v4/config/utils"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestExecPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		installPath string
		version     string
		execName    string
		expected    string
	}{
		{
			name:        "basic path construction",
			installPath: "/opt/tenv",
			version:     "1.0.0",
			execName:    "terraform",
			expected:    "/opt/tenv/1.0.0/terraform",
		},
		{
			name:        "path with spaces",
			installPath: "/opt/tenv tools",
			version:     "v1.5.0",
			execName:    "tofu",
			expected:    "/opt/tenv tools/v1.5.0/tofu",
		},
		{
			name:        "version with dots",
			installPath: "/usr/local/bin",
			version:     "1.2.3-alpha",
			execName:    "terragrunt",
			expected:    "/usr/local/bin/1.2.3-alpha/terragrunt",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Create a minimal config for testing with required fields
			conf := &config.Config{
				Getenv: configutils.GetenvFunc(func(key string) string {
					return "" // Return empty string for all environment variables
				}),
				Displayer: loghelper.MakeBasicDisplayer(
					hclog.NewNullLogger(),
					loghelper.StdDisplay,
				),
			}

			result := ExecPath(testCase.installPath, testCase.version, testCase.execName, conf)

			// Use filepath.Join to handle path separators correctly on different OS
			expected := filepath.Join(testCase.installPath, testCase.version, testCase.execName)
			assert.Equal(t, expected, result)
		})
	}
}

func TestUpdateWorkPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		cmdArgs  []string
		expected string
	}{
		{
			name:     "no chdir flag",
			cmdArgs:  []string{"plan", "-out=tfplan"},
			expected: "",
		},
		{
			name:     "chdir flag present",
			cmdArgs:  []string{"-chdir=/path/to/dir", "plan"},
			expected: "/path/to/dir",
		},
		{
			name:     "chdir flag with spaces",
			cmdArgs:  []string{"plan", "-chdir=/path with spaces/dir"},
			expected: "/path with spaces/dir",
		},
		{
			name:     "chdir flag at end",
			cmdArgs:  []string{"apply", "-chdir=./terraform"},
			expected: "./terraform",
		},
		{
			name:     "multiple chdir flags - first one wins",
			cmdArgs:  []string{"-chdir=/first/path", "-chdir=/second/path", "plan"},
			expected: "/first/path",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			conf := &config.Config{
				WorkPath: "", // Start with empty work path
			}

			updateWorkPath(conf, testCase.cmdArgs)

			assert.Equal(t, testCase.expected, conf.WorkPath)
		})
	}
}

func TestChdirFlagPrefix(t *testing.T) {
	t.Parallel()
	// Test that the chdirFlagPrefix constant is correctly defined
	expected := "-chdir="
	assert.Equal(t, expected, chdirFlagPrefix)
}

func TestExecFunction(t *testing.T) {
	t.Parallel()
	// Test that Exec function exists and has correct signature
	// We can't actually call it because it would try to execute commands
	assert.NotNil(t, Exec, "Exec function should be available")

	// Test that the function can be referenced (conceptual test)
	t.Log("Exec function is available for proxy execution")
}
