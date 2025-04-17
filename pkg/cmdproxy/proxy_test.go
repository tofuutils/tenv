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

package cmdproxy_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/cmdproxy"
)

func TestWriteMultiline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		key       string
		value     string
		wantErr   bool
		wantWrite string
	}{
		{
			name:      "basic write",
			key:       "test",
			value:     "value",
			wantErr:   false,
			wantWrite: "test<<ghadelimeter_",
		},
		{
			name:    "key contains delimiter",
			key:     "test<<delimiter",
			value:   "value",
			wantErr: true,
		},
		{
			name:    "value contains delimiter",
			key:     "test",
			value:   "value<<delimiter",
			wantErr: true,
		},
		{
			name:      "empty values",
			key:       "",
			value:     "",
			wantErr:   false,
			wantWrite: "<<ghadelimeter_",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := cmdproxy.WriteMultiline(&buf, tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteMultiline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got := buf.String()
				if !strings.Contains(got, tt.wantWrite) {
					t.Errorf("WriteMultiline() output = %v, want contains %v", got, tt.wantWrite)
				}
				// Verify the output format
				if !strings.Contains(got, tt.key) || !strings.Contains(got, tt.value) {
					t.Errorf("WriteMultiline() output missing key or value")
				}
			}
		})
	}
}

func TestRunBasicCommand(t *testing.T) {
	t.Parallel()

	// Create a simple command that should succeed
	cmd := exec.Command("echo", "test")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo test")
	}

	// Create a temporary file for output
	tempDir, err := os.MkdirTemp("", "cmdproxy_test_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "output.txt")
	os.Setenv("GITHUB_OUTPUT", outputPath)
	defer os.Unsetenv("GITHUB_OUTPUT")

	// Test without GHA mode
	cmdproxy.Run(cmd, false, os.Getenv)

	// Verify the command executed successfully
	if cmd.ProcessState == nil {
		t.Error("Command process state is nil")
	} else if !cmd.ProcessState.Success() {
		t.Error("Command did not execute successfully")
	}
}

func TestRunWithError(t *testing.T) {
	t.Parallel()

	// Create a command that should fail
	cmd := exec.Command("nonexistent_command")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "nonexistent_command")
	}

	// Create a temporary file for output
	tempDir, err := os.MkdirTemp("", "cmdproxy_test_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "output.txt")
	os.Setenv("GITHUB_OUTPUT", outputPath)
	defer os.Unsetenv("GITHUB_OUTPUT")

	// Test without GHA mode
	cmdproxy.Run(cmd, false, os.Getenv)

	// Verify the command failed
	if cmd.ProcessState == nil {
		t.Error("Command process state is nil")
	} else if cmd.ProcessState.Success() {
		t.Error("Command should have failed")
	}
}

func TestRunWithGHA(t *testing.T) {
	t.Parallel()

	// Create a simple command that should succeed
	cmd := exec.Command("echo", "test")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo test")
	}

	// Create a temporary file for output
	tempDir, err := os.MkdirTemp("", "cmdproxy_test_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "output.txt")
	os.Setenv("GITHUB_OUTPUT", outputPath)
	defer os.Unsetenv("GITHUB_OUTPUT")

	// Test with GHA mode
	cmdproxy.Run(cmd, true, os.Getenv)

	// Verify the command executed successfully
	if cmd.ProcessState == nil {
		t.Error("Command process state is nil")
	} else if !cmd.ProcessState.Success() {
		t.Error("Command did not execute successfully")
	}

	// Verify the output file was created and contains the expected content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal("Failed to read output file:", err)
	}

	// Check for expected keys in the output
	expectedKeys := []string{"stderr", "stdout", "exitcode"}
	for _, key := range expectedKeys {
		if !strings.Contains(string(content), key) {
			t.Errorf("Output file missing key: %s", key)
		}
	}
}

func TestRunWithCustomEnv(t *testing.T) {
	t.Parallel()

	// Create a simple command that should succeed
	cmd := exec.Command("echo", "test")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo test")
	}

	// Create a temporary file for output
	tempDir, err := os.MkdirTemp("", "cmdproxy_test_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "output.txt")

	// Test with custom environment function
	customEnv := func(key string) string {
		if key == "GITHUB_OUTPUT" {
			return outputPath
		}
		return ""
	}

	// Test with GHA mode and custom env
	cmdproxy.Run(cmd, true, customEnv)

	// Verify the command executed successfully
	if cmd.ProcessState == nil {
		t.Error("Command process state is nil")
	} else if !cmd.ProcessState.Success() {
		t.Error("Command did not execute successfully")
	}

	// Verify the output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}
