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

package zip_test

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	ziputil "github.com/tofuutils/tenv/v4/pkg/zip"
)

func createTestZip(files map[string][]byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write(content); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestUnzipToDir(t *testing.T) {
	t.Parallel()

	// Create test data
	files := map[string][]byte{
		"file1.txt":           []byte("content1"),
		"dir1/file2.txt":      []byte("content2"),
		"dir1/dir2/file3.txt": []byte("content3"),
		"dir1/":               nil, // directory entry
	}

	zipData, err := createTestZip(files)
	if err != nil {
		t.Fatal("Failed to create test zip:", err)
	}

	tests := []struct {
		name    string
		zipData []byte
		filter  func(string) bool
		wantErr bool
	}{
		{
			name:    "basic unzip",
			zipData: zipData,
			filter:  func(string) bool { return true },
			wantErr: false,
		},
		{
			name:    "filtered unzip",
			zipData: zipData,
			filter:  func(path string) bool { return filepath.Base(path) == "file1.txt" },
			wantErr: false,
		},
		{
			name:    "invalid zip data",
			zipData: []byte("not a zip file"),
			filter:  func(string) bool { return true },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create temp directory for each test
			tempDir, err := os.MkdirTemp("", "zip_test_*")
			if err != nil {
				t.Fatal("Failed to create temp dir:", err)
			}
			defer os.RemoveAll(tempDir)

			// Test unzip
			err = ziputil.UnzipToDir(tt.zipData, tempDir, tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnzipToDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify results for successful cases
			for name, content := range files {
				path := filepath.Join(tempDir, name)

				// Skip directory entries in the verification
				if name[len(name)-1] == '/' {
					continue
				}

				if tt.filter(path) {
					// Check file exists and content matches
					gotContent, err := os.ReadFile(path)
					if err != nil {
						t.Errorf("Failed to read file %s: %v", name, err)
						continue
					}
					if !bytes.Equal(gotContent, content) {
						t.Errorf("File %s content mismatch: got %q, want %q", name, gotContent, content)
					}
				} else {
					// Check file doesn't exist if filtered out
					if _, err := os.Stat(path); !os.IsNotExist(err) {
						t.Errorf("File %s should not exist", name)
					}
				}
			}
		})
	}
}

func TestSanitizeArchivePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dirPath  string
		fileName string
		want     string
		wantErr  bool
	}{
		{
			name:     "valid path",
			dirPath:  "/tmp/test",
			fileName: "file.txt",
			want:     "/tmp/test/file.txt",
			wantErr:  false,
		},
		{
			name:     "valid nested path",
			dirPath:  "/tmp/test",
			fileName: "subdir/file.txt",
			want:     "/tmp/test/subdir/file.txt",
			wantErr:  false,
		},
		{
			name:     "path traversal attempt",
			dirPath:  "/tmp/test",
			fileName: "../file.txt",
			wantErr:  true,
		},
		{
			name:     "absolute path attempt",
			dirPath:  "/tmp/test",
			fileName: "/etc/passwd",
			wantErr:  true,
		},
		{
			name:     "multiple path traversal",
			dirPath:  "/tmp/test",
			fileName: "../../../etc/passwd",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			dirPath:  "/tmp/test",
			fileName: "",
			want:     "/tmp/test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ziputil.SanitizeArchivePath(tt.dirPath, tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("sanitizeArchivePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("sanitizeArchivePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnzipWithEmptyFilter(t *testing.T) {
	t.Parallel()

	files := map[string][]byte{
		"file1.txt": []byte("content1"),
	}

	zipData, err := createTestZip(files)
	if err != nil {
		t.Fatal("Failed to create test zip:", err)
	}

	tempDir, err := os.MkdirTemp("", "zip_empty_filter_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with filter that rejects everything
	err = ziputil.UnzipToDir(zipData, tempDir, func(string) bool { return false })
	if err != nil {
		t.Error("Unexpected error with empty filter:", err)
	}

	// Verify no files were extracted
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal("Failed to read directory:", err)
	}

	if len(entries) > 0 {
		t.Error("Expected no files to be extracted")
	}
}

func TestUnzipWithFilePermissions(t *testing.T) {
	t.Parallel()

	// Create a zip with files having different permissions
	files := map[string][]byte{
		"executable.sh": []byte("#!/bin/sh\necho hello"),
		"readonly.txt":  []byte("read only content"),
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Create executable file
	execHeader := &zip.FileHeader{
		Name:   "executable.sh",
		Method: zip.Deflate,
	}
	execHeader.SetMode(0755)
	f1, err := w.CreateHeader(execHeader)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f1.Write(files["executable.sh"]); err != nil {
		t.Fatal(err)
	}

	// Create read-only file
	readHeader := &zip.FileHeader{
		Name:   "readonly.txt",
		Method: zip.Deflate,
	}
	readHeader.SetMode(0444)
	f2, err := w.CreateHeader(readHeader)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f2.Write(files["readonly.txt"]); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp("", "zip_perms_test_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	// Test unzip
	err = ziputil.UnzipToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	if err != nil {
		t.Fatal("Failed to unzip:", err)
	}

	// Verify file permissions
	execPath := filepath.Join(tempDir, "executable.sh")
	info, err := os.Stat(execPath)
	if err != nil {
		t.Fatal("Failed to stat executable file:", err)
	}
	if info.Mode().Perm() != 0755 {
		t.Errorf("Executable file has wrong permissions: got %v, want %v", info.Mode().Perm(), 0755)
	}

	readPath := filepath.Join(tempDir, "readonly.txt")
	info, err = os.Stat(readPath)
	if err != nil {
		t.Fatal("Failed to stat readonly file:", err)
	}
	if info.Mode().Perm() != 0444 {
		t.Errorf("Read-only file has wrong permissions: got %v, want %v", info.Mode().Perm(), 0444)
	}
}

func TestUnzipWithEmptyDirectories(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Create empty directories
	dirs := []string{
		"empty1/",
		"empty2/",
		"nested/empty3/",
		"nested/empty4/",
	}

	for _, dir := range dirs {
		_, err := w.Create(dir)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp("", "zip_empty_dirs_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	// Test unzip
	err = ziputil.UnzipToDir(buf.Bytes(), tempDir, func(string) bool { return true })
	if err != nil {
		t.Fatal("Failed to unzip:", err)
	}

	// Verify directories were created
	for _, dir := range dirs {
		path := filepath.Join(tempDir, dir)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("Failed to stat directory %s: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("Expected %s to be a directory", dir)
		}
	}
}

func TestUnzipWithCorruptedEntries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		zipData []byte
		wantErr bool
	}{
		{
			name:    "truncated zip",
			zipData: []byte("PK\x03\x04"), // Valid zip header but truncated
			wantErr: true,
		},
		{
			name:    "corrupted central directory",
			zipData: append([]byte("PK\x03\x04"), make([]byte, 100)...),
			wantErr: true,
		},
		{
			name:    "empty zip",
			zipData: []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tempDir, err := os.MkdirTemp("", "zip_corrupted_*")
			if err != nil {
				t.Fatal("Failed to create temp dir:", err)
			}
			defer os.RemoveAll(tempDir)

			err = ziputil.UnzipToDir(tt.zipData, tempDir, func(string) bool { return true })
			if (err != nil) != tt.wantErr {
				t.Errorf("UnzipToDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnzipWithDuplicateNames(t *testing.T) {
	t.Parallel()

	// Create a zip with duplicate file names in different cases
	files := map[string][]byte{
		"file.txt":        []byte("content1"),
		"FILE.txt":        []byte("content2"),
		"file.TXT":        []byte("content3"),
		"nested/file.txt": []byte("content4"),
	}

	zipData, err := createTestZip(files)
	if err != nil {
		t.Fatal("Failed to create test zip:", err)
	}

	tempDir, err := os.MkdirTemp("", "zip_duplicates_*")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	// Test unzip
	err = ziputil.UnzipToDir(zipData, tempDir, func(string) bool { return true })
	if err != nil {
		t.Fatal("Failed to unzip:", err)
	}

	// Verify files based on OS case sensitivity
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal("Failed to read directory:", err)
	}

	// Count how many "file.txt" variants we find
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.ToLower(entry.Name()) == "file.txt" {
			count++
		}
	}

	// On case-sensitive systems, we should find 3 files
	// On case-insensitive systems, we should find 1 file
	expectedCount := 3
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		expectedCount = 1
	}

	if count != expectedCount {
		t.Errorf("Expected %d file.txt variants, got %d", expectedCount, count)
	}

	// Nested file should always exist
	nestedPath := filepath.Join(tempDir, "nested", "file.txt")
	if _, err := os.Stat(nestedPath); err != nil {
		t.Errorf("Nested file.txt not found: %v", err)
	}
}
