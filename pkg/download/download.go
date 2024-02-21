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
	"io"
	"net/http"
	"net/url"
)

func ApplyUrlTranformer(urlTransformer func(string) (string, error), baseURLs ...string) ([]string, error) {
	transformedURLs := make([]string, 0, len(baseURLs))
	for _, baseURL := range baseURLs {
		transformedURL, err := urlTransformer(baseURL)
		if err != nil {
			return nil, err
		}

		transformedURLs = append(transformedURLs, transformedURL)
	}

	return transformedURLs, nil
}

func Bytes(url string, display func(string)) ([]byte, error) {
	display("Downloading " + url)

	response, err := http.Get(url) //nolint
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func UrlTranformer(rewriteRule []string) func(string) (string, error) {
	if len(rewriteRule) < 2 {
		return noTransform
	}

	prevBaseURL := rewriteRule[0]
	baseURL := rewriteRule[1]
	prevLen := len(prevBaseURL)
	if prevLen == 0 || baseURL == "" {
		return noTransform
	}

	return func(urlValue string) (string, error) {
		if len(urlValue) < prevLen || urlValue[:prevLen] != prevBaseURL {
			return urlValue, nil
		}

		return url.JoinPath(baseURL, urlValue[prevLen:]) //nolint
	}
}

func noTransform(value string) (string, error) {
	return value, nil
}
