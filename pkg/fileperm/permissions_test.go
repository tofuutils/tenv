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

package fileperm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/fileperm"
)

func TestPermissionConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		perm     os.FileMode
		expected os.FileMode
		desc     string
	}{
		{
			name:     "RW permissions",
			perm:     fileperm.RW,
			expected: 0o644,
			desc:     "Read-write for owner, read-only for group and others",
		},
		{
			name:     "RWE permissions",
			perm:     fileperm.RWE,
			expected: 0o755,
			desc:     "Read-write-execute for owner, read-execute for group and others",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.perm != tt.expected {
				t.Errorf("%s: got %o, want %o (%s)", tt.name, tt.perm, tt.expected, tt.desc)
			}
		})
	}
}

func TestPermissionUsage(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fileperm_test_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	// Test RW permissions
	filePath := filepath.Join(tempDir, "test_rw.txt")
	if err := os.WriteFile(filePath, []byte("test"), fileperm.RW); err != nil {
		t.Fatal("Failed to create test file:", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatal("Failed to stat test file:", err)
	}

	// Check that the actual permissions match the expected ones
	// Note: The actual permissions will be affected by umask
	expectedPerm := os.FileMode(fileperm.RW)
	if info.Mode().Perm()&0o777 != expectedPerm {
		t.Errorf("File permissions mismatch: got %o, want %o", info.Mode().Perm()&0o777, expectedPerm)
	}

	// Test RWE permissions
	dirPath := filepath.Join(tempDir, "test_dir")
	if err := os.Mkdir(dirPath, fileperm.RWE); err != nil {
		t.Fatal("Failed to create test directory:", err)
	}

	info, err = os.Stat(dirPath)
	if err != nil {
		t.Fatal("Failed to stat test directory:", err)
	}

	// Check that the actual permissions match the expected ones
	// Note: The actual permissions will be affected by umask
	expectedPerm = os.FileMode(fileperm.RWE)
	if info.Mode().Perm()&0o777 != expectedPerm {
		t.Errorf("Directory permissions mismatch: got %o, want %o", info.Mode().Perm()&0o777, expectedPerm)
	}
}
