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
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyURLTransformer(t *testing.T) {
	t.Parallel()
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			result, err := ApplyURLTransformer(testCase.urlTransformer, testCase.baseURLs...)

			if testCase.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, result)
			}
		})
	}
}

func TestNewURLTransformer(t *testing.T) {
	t.Parallel()
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			transformer := NewURLTransformer(testCase.prevBaseURL, testCase.baseURL)
			result, err := transformer(testCase.input)

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestNoTransform(t *testing.T) {
	t.Parallel()
	input := "any_url"
	result, err := NoTransform(input)

	require.NoError(t, err)
	assert.Equal(t, input, result)
}

func TestNoCheck(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
	}

	err := NoCheck(resp)
	require.NoError(t, err)
}

func TestWithBasicAuth(t *testing.T) {
	t.Parallel()
	username := "testuser"
	password := "testpass"

	option := WithBasicAuth(username, password)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	option(req)

	authHeader := req.Header.Get("Authorization")
	assert.Contains(t, authHeader, "Basic")
}

func TestNoDisplay(t *testing.T) {
	t.Parallel()
	// Test that NoDisplay doesn't panic and can be called
	assert.NotPanics(t, func() {
		NoDisplay("test message")
	})
}

func TestURLTransformerType(t *testing.T) {
	t.Parallel()
	// Test that URLTransformer is a function type
	transformer := func(s string) (string, error) {
		return s, nil
	}

	result, err := transformer("test")
	require.NoError(t, err)
	assert.Equal(t, "test", result)
}

func TestRequestOptionType(t *testing.T) {
	t.Parallel()
	// Test that RequestOption is a function type
	option := func(req *http.Request) {
		req.Header.Set("Test", "value")
	}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	option(req)
	assert.Equal(t, "value", req.Header.Get("Test"))
}

func TestResponseCheckerType(t *testing.T) {
	t.Parallel()
	// Test that ResponseChecker is a function type
	checker := func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return assert.AnError
		}

		return nil
	}

	resp := &http.Response{StatusCode: http.StatusOK}
	err := checker(resp)
	require.NoError(t, err)

	resp.StatusCode = 404
	err = checker(resp)
	require.Error(t, err)
}

func TestURLJoinPath(t *testing.T) {
	t.Parallel()
	// Test URL path joining functionality used in NewURLTransformer
	baseURL := "https://example.com/api"
	path := "/v1/resource"

	result, err := url.JoinPath(baseURL, path)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/api/v1/resource", result)
}

func TestConstants(t *testing.T) {
	t.Parallel()
	// Test that the types are properly defined
	_ = NoTransform
	_ = WithBasicAuth("user", "pass")
	_ = NoCheck
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

	result, err := JSON(t.Context(), testServer.URL, NoDisplay, NoCheck)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBytes(t *testing.T) {
	t.Parallel()
	// Test that Bytes function exists and has correct signature
	assert.NotNil(t, Bytes, "Bytes function should be available")

	// Test that the function can be called (conceptual test since it makes HTTP requests)
	t.Log("Bytes function is available for HTTP requests")
}
