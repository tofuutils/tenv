/*
 *
 * Copyright 2024 gotofuenv authors.
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

package github

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const pageQuery = "?page="

var (
	errNoAsset = errors.New("asset not found for current platform")
	errReturn  = errors.New("unexpected value returned by API")
)

func DownloadAssetUrl(tag string, searchedAssetName string, githubReleaseUrl string, authorizationHeader string) (string, error) {
	releaseUrl, err := url.JoinPath(githubReleaseUrl, "tags", tag)
	if err != nil {
		return "", err
	}

	value, err := apiGetRequest(releaseUrl, authorizationHeader)
	if err != nil {
		return "", err
	}

	object, _ := value.(map[string]any)
	baseAssetsUrl, ok := object["assets_url"].(string)
	if !ok {
		return "", errReturn
	}

	page := 1
	baseAssetsUrl += pageQuery
	for {
		assetsUrl := baseAssetsUrl + strconv.Itoa(page)
		value, err = apiGetRequest(assetsUrl, authorizationHeader)
		if err != nil {
			return "", err
		}

		values, ok := value.([]any)
		if !ok {
			return "", errReturn
		}

		if len(values) == 0 {
			return "", errNoAsset
		}

		for _, value := range values {
			object, _ = value.(map[string]any)
			assetName, ok := object["name"].(string)
			if !ok {
				return "", errReturn
			}

			if assetName != searchedAssetName {
				continue
			}

			downloadUrl, ok := object["browser_download_url"].(string)
			if !ok {
				return "", errReturn
			}
			return downloadUrl, nil
		}
		page++
	}
}

func LatestRelease(githubReleaseUrl string, authorizationHeader string) (string, error) {
	latestUrl, err := url.JoinPath(githubReleaseUrl, "latest")
	if err != nil {
		return "", err
	}

	value, err := apiGetRequest(latestUrl, authorizationHeader)
	if err != nil {
		return "", err
	}

	version, ok := extractVersion(value)
	if !ok {
		return "", errReturn
	}
	return version, nil
}

func ListReleases(githubReleaseUrl string, authorizationHeader string) ([]string, error) {
	basePageUrl := githubReleaseUrl + pageQuery

	page := 1
	var releases []string
	for {
		pageUrl := basePageUrl + strconv.Itoa(page)
		value, err := apiGetRequest(pageUrl, authorizationHeader)
		if err != nil {
			return nil, err
		}

		values, ok := value.([]any)
		if !ok {
			return nil, errReturn
		}

		if len(values) == 0 {
			return releases, nil
		}

		for _, value := range values {
			version, ok := extractVersion(value)
			if !ok {
				return nil, errReturn
			}
			releases = append(releases, version)
		}
		page++
	}
}

func apiGetRequest(callUrl string, authorizationHeader string) (any, error) {
	request, err := http.NewRequest(http.MethodGet, callUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	if authorizationHeader != "" {
		request.Header.Set("Authorization", authorizationHeader)
	}
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var value any
	err = json.Unmarshal(data, &value)
	return value, err
}

func BuildAuthorizationHeader(token string) string {
	if token == "" {
		return ""
	}

	var authorizationBuilder strings.Builder
	authorizationBuilder.WriteString("Bearer ")
	authorizationBuilder.WriteString(token)
	return authorizationBuilder.String()
}

func extractVersion(value any) (string, bool) {
	object, _ := value.(map[string]any)
	version, _ := object["tag_name"].(string)
	if version == "" {
		return "", false
	}

	// version returned without starting 'v'
	if version[0] == 'v' {
		version = version[1:]
	}
	return version, true
}
