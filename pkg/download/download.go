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
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var ErrPrefix = errors.New("prefix does not match")

func Bytes(url string, verbose bool) ([]byte, error) {
	if verbose {
		fmt.Println("Downloading", url) //nolint
	}

	response, err := http.Get(url) //nolint
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func UrlTranformer(remoteConf map[string]string) func(string) (string, error) {
	prevBaseURL := remoteConf["old_base_url"]
	prevLen := len(prevBaseURL)
	baseURL := remoteConf["new_base_url"]
	if prevLen == 0 || baseURL == "" {
		return noTransform
	}

	return func(urlValue string) (string, error) {
		if len(urlValue) < prevLen || urlValue[:prevLen] != prevBaseURL {
			return "", ErrPrefix
		}

		return url.JoinPath(baseURL, urlValue[prevLen:]) //nolint
	}
}

func noTransform(value string) (string, error) {
	return value, nil
}
