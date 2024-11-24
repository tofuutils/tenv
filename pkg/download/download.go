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
	"io"
	"net/http"
	"net/url"
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

func Bytes(ctx context.Context, url string, display func(string), checker ResponseChecker, requestOptions ...RequestOption) ([]byte, error) {
	display("Downloading " + url)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
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

func JSON(ctx context.Context, url string, display func(string), checker ResponseChecker, requestOptions ...RequestOption) (any, error) {
	data, err := Bytes(ctx, url, display, checker, requestOptions...)
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

func NoCheck(*http.Response) error {
	return nil
}
