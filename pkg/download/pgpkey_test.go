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
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestGetPGPKey(t *testing.T) {
	t.Parallel()

	testKeyContent := []byte("test pgp key content")

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/key.txt" {
			_, _ = w.Write(testKeyContent)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer testServer.Close()

	correctURL, err := url.JoinPath(testServer.URL, "key.txt")
	if err != nil {
		t.Fatal("Failed to initialize test URL : ", err)
	}

	notCorrectURL, err := url.JoinPath(testServer.URL, "key2.txt")
	if err != nil {
		t.Fatal("Failed to initialize test URL : ", err)
	}

	tmpDir := t.TempDir()
	correctPath := filepath.Join(tmpDir, "test-key.txt")
	notCorrectPath := filepath.Join(tmpDir, "non-existent-key.txt")

	if err := os.WriteFile(correctPath, testKeyContent, 0o600); err != nil {
		t.Fatalf("Failed to create test key file: %v", err)
	}

	testGet := func(name string, pathOrUrl string, wantErr bool) {
		t.Helper()

		t.Run(name, func(t *testing.T) {
			t.Helper()

			got, err := GetPGPKey(context.Background(), pathOrUrl, NoDisplay)
			switch {
			case wantErr:
				if err == nil {
					t.Error("GetPGPKey() should fail on ", pathOrUrl)
				}
			case err != nil:
				t.Error("GetPGPKey() returned an unexpected error : ", err)
			case !slices.Equal(testKeyContent, got):
				t.Error("GetPGPKey() returned unexpected content : ", string(got))
			}
		})
	}

	testGet("empty path or URL", "", true) // not handled by GetPGPKey, config has default URL
	testGet("http URL", correctURL, false)
	testGet("non-existent http URL", notCorrectURL, true)
	testGet("local file path", correctPath, false)
	testGet("non-existent local file path", notCorrectPath, true)
}
