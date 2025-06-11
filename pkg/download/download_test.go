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
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/download"
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

func TestApplyURLTransformer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		prevBaseURL   string
		baseURL       string
		inputURLs     []string
		expectedURLs  []string
		expectedError bool
	}{
		{
			name:        "transform multiple URLs",
			prevBaseURL: "https://releases.hashicorp.com",
			baseURL:     "http://localhost:8080",
			inputURLs: []string{
				"https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_linux_386.zip",
				"https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_darwin_amd64.zip",
			},
			expectedURLs: []string{
				"http://localhost:8080/terraform/1.7.0/terraform_1.7.0_linux_386.zip",
				"http://localhost:8080/terraform/1.7.0/terraform_1.7.0_darwin_amd64.zip",
			},
			expectedError: false,
		},
		{
			name:          "empty base URL",
			prevBaseURL:   "https://releases.hashicorp.com",
			baseURL:       "",
			inputURLs:     []string{"https://releases.hashicorp.com/test"},
			expectedURLs:  nil,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transformer := download.NewURLTransformer(tt.prevBaseURL, tt.baseURL)
			transformedURLs, err := download.ApplyURLTransformer(transformer, tt.inputURLs...)

			if (err != nil) != tt.expectedError {
				t.Errorf("ApplyURLTransformer() error = %v, wantErr %v", err, tt.expectedError)
				return
			}

			if !tt.expectedError && !slices.Equal(transformedURLs, tt.expectedURLs) {
				t.Errorf("ApplyURLTransformer() = %v, want %v", transformedURLs, tt.expectedURLs)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	t.Parallel()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test data"))
	}))
	defer server.Close()

	// Test successful download
	ctx := context.Background()
	data, err := download.Bytes(ctx, server.URL, download.NoDisplay, download.NoCheck)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if string(data) != "test data" {
		t.Error("Unexpected data:", string(data))
	}

	// Test with invalid URL
	_, err = download.Bytes(ctx, "invalid://url", download.NoDisplay, download.NoCheck)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}

	// Test with context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = download.Bytes(cancelCtx, server.URL, download.NoDisplay, download.NoCheck)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"key": "value"}`))
	}))
	defer server.Close()

	// Test successful JSON download
	ctx := context.Background()
	value, err := download.JSON(ctx, server.URL, download.NoDisplay, download.NoCheck)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	// Verify JSON structure
	jsonMap, ok := value.(map[string]interface{})
	if !ok {
		t.Error("Expected map[string]interface{}")
	}
	if jsonMap["key"] != "value" {
		t.Error("Unexpected JSON value:", jsonMap["key"])
	}

	// Test with invalid JSON
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	_, err = download.JSON(ctx, server.URL, download.NoDisplay, download.NoCheck)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestWithBasicAuth(t *testing.T) {
	t.Parallel()

	// Create a test server that checks auth
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "user" || password != "pass" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte("authenticated"))
	}))
	defer server.Close()

	// Test with basic auth
	ctx := context.Background()
	authOption := download.WithBasicAuth("user", "pass")
	data, err := download.Bytes(ctx, server.URL, download.NoDisplay, download.NoCheck, authOption)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if string(data) != "authenticated" {
		t.Error("Unexpected response:", string(data))
	}
}

func TestNoTransform(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "normal URL",
			input: "https://example.com",
		},
		{
			name:  "empty URL",
			input: "",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := download.NoTransform(tt.input)
			if err != nil {
				t.Error("Unexpected error:", err)
			}
			if result != tt.input {
				t.Errorf("NoTransform() = %v, want %v", result, tt.input)
			}
		})
	}
}
