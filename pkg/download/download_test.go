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

package download_test

import (
	"testing"

	"context"
	"github.com/tofuutils/tenv/v4/pkg/download"
	"os"
	"path/filepath"
)

func TestURLTransformer(t *testing.T) {
	t.Parallel()

	urlTransformer := download.NewURLTransformer("https://releases.hashicorp.com", "http://localhost:8080")

	value, err := urlTransformer("https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_linux_386.zip")
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	if value != "http://localhost:8080/terraform/1.7.0/terraform_1.7.0_linux_386.zip" {
		t.Error("Unexpected result, get :", value)
	}
}

func TestURLTransformerPrefix(t *testing.T) {
	t.Parallel()

	urlTransformer := download.NewURLTransformer("https://github.com", "https://go.dev")

	initialValue := "https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_darwin_amd64.zip"
	value, err := urlTransformer(initialValue)
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	if value != initialValue {
		t.Error("Unexpected result, get :", value)
	}
}

func TestURLTransformerSlash(t *testing.T) {
	t.Parallel()

	urlTransformer := download.NewURLTransformer("https://releases.hashicorp.com/", "http://localhost")

	value, err := urlTransformer("https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_darwin_amd64.zip")
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	if value != "http://localhost/terraform/1.7.0/terraform_1.7.0_darwin_amd64.zip" {
		t.Error("Unexpected result, get :", value)
	}
}

func TestJSONFromFileAbsolutePath(t *testing.T) {
	t.Parallel()

	//create a temporary file with valid JSON content
	tmpdir := t.TempDir()
	tmpfile := filepath.Join(tmpdir, "test.json")
	testdata := []byte(`{"key": "value"}`)

	if err := os.WriteFile(tmpfile, testdata, 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	//convert to file:// URL
	fileURL := "file://" + tmpfile

	result, err := download.JSON(context.Background(), fileURL, download.NoDisplay, download.NoCheck)
	if err != nil {
		t.Fatalf("JSON() failed: %v", err)
	}

	//verify it parsed correctly
	if result == nil {
		t.Error("JSON() returned nil")
	}
}

func TestJSONFromFileRelativePath(t *testing.T) {
	t.Parallel()

	//create a temporary file with valid JSON content
	tmpdir := t.TempDir()
	testfile := "test.json"
	testpath := filepath.Join(tmpdir, testfile)
	testdata := []byte(`{"key": "value"}`)

	if err := os.WriteFile(testpath, testdata, 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	//change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpdir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	//test with relative path
	fileURL := "file://" + testfile

	result, err := download.JSON(context.Background(), fileURL, download.NoDisplay, download.NoCheck)
	if err != nil {
		t.Fatalf("JSON() failed: %v", err)
	}

	//verify it parsed correctly
	if result == nil {
		t.Error("JSON() returned nil")
	}
}

func TestJSONFromFileMissing(t *testing.T) {
	t.Parallel()

	//test with non-existent file
	fileURL := "file:///nonexistent.json"

	_, err := download.JSON(context.Background(), fileURL, download.NoDisplay, download.NoCheck)
	if err == nil {
		t.Fatal("JSON() should fail for nonexistent file, got nil")
	}
}

func TestJSONFromFileInvalidJSON(t *testing.T) {
	t.Parallel()

	//create a temporary file with invalid JSON content
	tmpdir := t.TempDir()
	tmpfile := filepath.Join(tmpdir, "invalid.json")

	// write invalid JSON
	if err := os.WriteFile(tmpfile, []byte(`{invalid json}`), 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	//convert to file:// URL
	fileURL := "file://" + tmpfile

	_, err := download.JSON(context.Background(), fileURL, download.NoDisplay, download.NoCheck)
	if err == nil {
		t.Fatal("JSON() should fail for invalid JSON, got nil")
	}
}
