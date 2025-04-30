/*
 *
 * Copyright 2025 tofuutils authors.
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

package download

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestGetPGPKey(t *testing.T) {
	t.Parallel()

	// Create a temporary test file
	tmpDir := t.TempDir()
	testKeyPath := filepath.Join(tmpDir, "test-key.txt")
	testKeyContent := []byte("test pgp key content")
	if err := os.WriteFile(testKeyPath, testKeyContent, 0644); err != nil {
		t.Fatalf("Failed to create test key file: %v", err)
	}

	tests := []struct {
		name        string
		keyPath     string
		wantErr     bool
		checkResult func([]byte) bool
	}{
		{
			name:    "empty path uses default URL",
			keyPath: "",
			wantErr: false,
			checkResult: func(data []byte) bool {
				return len(data) > 0 // Default URL should return some content
			},
		},
		{
			name:    "local file path",
			keyPath: testKeyPath,
			wantErr: false,
			checkResult: func(data []byte) bool {
				return string(data) == string(testKeyContent)
			},
		},
		{
			name:    "http URL",
			keyPath: "http://example.com/key.txt",
			wantErr: true, // Should fail as URL doesn't exist
		},
		{
			name:    "https URL",
			keyPath: "https://example.com/key.txt",
			wantErr: true, // Should fail as URL doesn't exist
		},
		{
			name:    "non-existent local file",
			keyPath: filepath.Join(tmpDir, "non-existent.txt"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetPGPKey(context.Background(), tt.keyPath, func(string) {})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPGPKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkResult != nil && !tt.checkResult(got) {
				t.Error("GetPGPKey() returned unexpected content")
			}
		})
	}
}
