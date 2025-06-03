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

package htmlquery

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

//go:embed testdata/artifactory.html
var artifactoryData []byte

func TestExtractAttr(t *testing.T) {
	t.Parallel()

	extractor := SelectionExtractor("style")
	extracted, err := extractList(artifactoryData, "address", extractor)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if len(extracted) != 1 || extracted[0] != "font-size:small;" {
		t.Error("Unmatching results, get :", extracted)
	}
}

func TestExtractText(t *testing.T) {
	t.Parallel()

	extractor := SelectionExtractor("#text")
	extracted, err := extractList(artifactoryData, "address", extractor)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if len(extracted) != 1 || extracted[0] != "Artifactory/7.68.20 Server" {
		t.Error("Unmatching results, get :", extracted)
	}
}

func TestExtractTexts(t *testing.T) {
	t.Parallel()

	extracted, err := extractList(artifactoryData, "a", selectionTextExtractor)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	expected := []string{"../", "1.7.0/", "1.7.0-alpha20231025/", "1.7.0-rc2/", "1.7.1/", "index.json"}
	if !slices.Equal(extracted, expected) {
		t.Error("Unmatching results, get :", extracted)
	}
}

func TestRequest(t *testing.T) {
	t.Parallel()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(artifactoryData)
	}))
	defer server.Close()

	// Test successful request
	ctx := context.Background()
	extracted, err := Request(ctx, server.URL, "a", selectionTextExtractor)
	if err != nil {
		t.Fatal("Unexpected request error:", err)
	}

	expected := []string{"../", "1.7.0/", "1.7.0-alpha20231025/", "1.7.0-rc2/", "1.7.1/", "index.json"}
	if !slices.Equal(extracted, expected) {
		t.Error("Unmatching results, get:", extracted)
	}

	// Test with invalid URL
	_, err = Request(ctx, "invalid://url", "a", selectionTextExtractor)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}

	// Test with non-existent selector
	extracted, err = Request(ctx, server.URL, "non-existent", selectionTextExtractor)
	if err != nil {
		t.Fatal("Unexpected error for non-existent selector:", err)
	}
	if len(extracted) != 0 {
		t.Error("Expected empty results for non-existent selector, got:", extracted)
	}
}

func TestSelectionExtractor(t *testing.T) {
	t.Parallel()

	// Test text extractor
	extractor := SelectionExtractor("#text")
	if extractor == nil {
		t.Error("Expected non-nil text extractor")
	}

	// Test attribute extractor
	extractor = SelectionExtractor("style")
	if extractor == nil {
		t.Error("Expected non-nil attribute extractor")
	}

	// Test with empty part
	extractor = SelectionExtractor("")
	if extractor == nil {
		t.Error("Expected non-nil extractor for empty part")
	}
}

func TestExtractListError(t *testing.T) {
	t.Parallel()

	// Test with invalid HTML
	invalidHTML := []byte("<invalid>html")
	_, err := extractList(invalidHTML, "a", selectionTextExtractor)
	if err == nil {
		t.Error("Expected error for invalid HTML")
	}

	// Test with nil extractor
	_, err = extractList(artifactoryData, "a", nil)
	if err == nil {
		t.Error("Expected error for nil extractor")
	}
}
