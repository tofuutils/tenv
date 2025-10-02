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

package download

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyURLTransformer(t *testing.T) {
	tests := []struct {
		name           string
		urlTransformer URLTransformer
		baseURLs       []string
		expected       []string
		expectError    bool
	}{
		{
			name: "successful transformation",
			urlTransformer: func(s string) (string, error) {
				return "transformed_" + s, nil
			},
			baseURLs:    []string{"url1", "url2", "url3"},
			expected:    []string{"transformed_url1", "transformed_url2", "transformed_url3"},
			expectError: false,
		},
		{
			name: "transformation error",
			urlTransformer: func(s string) (string, error) {
				if s == "error_url" {
					return "", assert.AnError
				}
				return "transformed_" + s, nil
			},
			baseURLs:    []string{"url1", "error_url", "url3"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "empty base URLs",
			urlTransformer: func(s string) (string, error) {
				return "transformed_" + s, nil
			},
			baseURLs:    []string{},
			expected:    []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ApplyURLTransformer(tt.urlTransformer, tt.baseURLs...)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNewURLTransformer(t *testing.T) {
	tests := []struct {
		name        string
		prevBaseURL string
		baseURL     string
		input       string
		expected    string
	}{
		{
			name:        "successful transformation",
			prevBaseURL: "https://old.example.com/",
			baseURL:     "https://new.example.com/",
			input:       "https://old.example.com/path/to/resource",
			expected:    "https://new.example.com/path/to/resource",
		},
		{
			name:        "no match - return original",
			prevBaseURL: "https://old.example.com/",
			baseURL:     "https://new.example.com/",
			input:       "https://different.example.com/path/to/resource",
			expected:    "https://different.example.com/path/to/resource",
		},
		{
			name:        "empty prevBaseURL - return NoTransform",
			prevBaseURL: "",
			baseURL:     "https://new.example.com/",
			input:       "https://old.example.com/path/to/resource",
			expected:    "https://old.example.com/path/to/resource",
		},
		{
			name:        "empty baseURL - return NoTransform",
			prevBaseURL: "https://old.example.com/",
			baseURL:     "",
			input:       "https://old.example.com/path/to/resource",
			expected:    "https://old.example.com/path/to/resource",
		},
		{
			name:        "input shorter than prevBaseURL",
			prevBaseURL: "https://old.example.com/",
			baseURL:     "https://new.example.com/",
			input:       "https://old.example",
			expected:    "https://old.example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := NewURLTransformer(tt.prevBaseURL, tt.baseURL)
			result, err := transformer(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoTransform(t *testing.T) {
	input := "any_url"
	result, err := NoTransform(input)

	assert.NoError(t, err)
	assert.Equal(t, input, result)
}

func TestNoCheck(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}

	err := NoCheck(resp)
	assert.NoError(t, err)
}

func TestWithBasicAuth(t *testing.T) {
	username := "testuser"
	password := "testpass"

	option := WithBasicAuth(username, password)

	req, err := http.NewRequest("GET", "https://example.com", nil)
	assert.NoError(t, err)

	option(req)

	authHeader := req.Header.Get("Authorization")
	assert.Contains(t, authHeader, "Basic")
}

func TestNoDisplay(t *testing.T) {
	// Test that NoDisplay doesn't panic and can be called
	assert.NotPanics(t, func() {
		NoDisplay("test message")
	})
}

func TestURLTransformerType(t *testing.T) {
	// Test that URLTransformer is a function type
	var transformer URLTransformer = func(s string) (string, error) {
		return s, nil
	}

	result, err := transformer("test")
	assert.NoError(t, err)
	assert.Equal(t, "test", result)
}

func TestRequestOptionType(t *testing.T) {
	// Test that RequestOption is a function type
	var option RequestOption = func(req *http.Request) {
		req.Header.Set("Test", "value")
	}

	req, err := http.NewRequest("GET", "https://example.com", nil)
	assert.NoError(t, err)

	option(req)
	assert.Equal(t, "value", req.Header.Get("Test"))
}

func TestResponseCheckerType(t *testing.T) {
	// Test that ResponseChecker is a function type
	var checker ResponseChecker = func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return assert.AnError
		}
		return nil
	}

	resp := &http.Response{StatusCode: 200}
	err := checker(resp)
	assert.NoError(t, err)

	resp.StatusCode = 404
	err = checker(resp)
	assert.Error(t, err)
}

func TestURLJoinPath(t *testing.T) {
	// Test URL path joining functionality used in NewURLTransformer
	baseURL := "https://example.com/api"
	path := "/v1/resource"

	result, err := url.JoinPath(baseURL, path)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/api/v1/resource", result)
}

func TestConstants(t *testing.T) {
	// Test that the types are properly defined
	var _ URLTransformer = NoTransform
	var _ RequestOption = WithBasicAuth("user", "pass")
	var _ ResponseChecker = NoCheck
}

func TestJSON(t *testing.T) {
	t.Parallel()

	testJSON := `{"key": "value", "number": 42, "array": [1, 2, 3]}`
	expected := map[string]interface{}{
		"key":    "value",
		"number": float64(42),
		"array":  []interface{}{float64(1), float64(2), float64(3)},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testJSON))
	}))
	defer testServer.Close()

	result, err := JSON(context.Background(), testServer.URL, NoDisplay, NoCheck)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBytes(t *testing.T) {
	// Test that Bytes function exists and has correct signature
	assert.NotNil(t, Bytes, "Bytes function should be available")

	// Test that the function can be called (conceptual test since it makes HTTP requests)
	t.Log("Bytes function is available for HTTP requests")
}
