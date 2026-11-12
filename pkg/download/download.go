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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RequestOption = func(*http.Request)

type ResponseChecker = func(*http.Response) error

func ApplyURLTransformer(urlTransformer URLTransformer, baseURLs ...string) ([]string, error) {
	transformedURLs := make([]string, len(baseURLs))
	for index, baseURL := range baseURLs {
		transformedURL, err := urlTransformer(baseURL)
		if err != nil {
			return nil, err
		}

		transformedURLs[index] = transformedURL
	}

	return transformedURLs, nil
}

func Bytes(ctx context.Context, urlStr string, display func(string), checker ResponseChecker, requestOptions ...RequestOption) ([]byte, error) {
	//Handle file:// URLs
	//renaming urlstr to url to avoid shadowing the package name (net/url)
	scheme, _, ok := strings.Cut(urlStr, "://")
	if ok && scheme == "file" {
		return bytesFromFile(urlStr)
	}

	display("Downloading " + urlStr)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
	if err != nil {
		return nil, err
	}

	for _, option := range requestOptions {
		option(request)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if err = checker(response); err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}

func JSON(ctx context.Context, urlStr string, display func(string), checker ResponseChecker, requestOptions ...RequestOption) (any, error) {
	data, err := Bytes(ctx, urlStr, display, checker, requestOptions...)
	if err != nil {
		return nil, err
	}

	var value any
	err = json.Unmarshal(data, &value)

	return value, err
}

func NoDisplay(string) {}

type URLTransformer = func(string) (string, error)

func NewURLTransformer(prevBaseURL string, baseURL string) URLTransformer {
	prevLen := len(prevBaseURL)
	if prevLen == 0 || baseURL == "" {
		return NoTransform
	}

	return func(urlValue string) (string, error) {
		if len(urlValue) < prevLen || urlValue[:prevLen] != prevBaseURL {
			return urlValue, nil
		}

		return url.JoinPath(baseURL, urlValue[prevLen:])
	}
}

func WithBasicAuth(username string, password string) RequestOption {
	return func(r *http.Request) {
		r.SetBasicAuth(username, password)
	}
}

func NoTransform(value string) (string, error) {
	return value, nil
}

func bytesFromFile(fileURL string) ([]byte, error) {
	//parse file:// URL to get the path
	//file://path/to/file -> /path/to/file
	//file://./relative/path/to/file -> ./relative/path/to/file

	filePath := strings.TrimPrefix(fileURL, "file://")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from %s: %w", fileURL, err)
	}
	return data, nil
}

func NoCheck(*http.Response) error {
	return nil
}
