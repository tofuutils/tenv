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

package tofuversion

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	"github.com/dvaumoron/gotofuenv/config"
)

var (
	errCast          = errors.New("value returned by API has not the expected type")
	errAssetNotFound = errors.New("asset not found for current platform")
)

// version must not start with a 'v'
func githubDownloadUrl(version string, conf *config.Config) (string, error) {
	releaseUrl, err := url.JoinPath(conf.RemoteUrl, "tags", "v"+version)
	if err != nil {
		return "", err
	}

	authorizationHeader := buildAuthorizationHeader(conf.Token)
	value, err := githubApiGetRequest(releaseUrl, authorizationHeader)
	if err != nil {
		return "", err
	}

	object, _ := value.(map[string]any)
	assetsUrl, ok := object["assets_url"].(string)
	if !ok {
		return "", errCast
	}

	value, err = githubApiGetRequest(assetsUrl, authorizationHeader)
	if err != nil {
		return "", err
	}

	values, ok := value.([]any)
	if !ok {
		return "", errCast
	}

	searchedAssetName := buildAssetName(version)
	for _, value := range values {
		object, _ = value.(map[string]any)
		assetName, ok := object["name"].(string)
		if !ok {
			return "", errCast
		}

		if assetName != searchedAssetName {
			continue
		}

		downloadUrl, ok := object["browser_download_url"].(string)
		if !ok {
			return "", errCast
		}
		return downloadUrl, nil
	}
	return "", errAssetNotFound
}

func githubLatestRelease(conf *config.Config) (string, error) {
	latestUrl, err := url.JoinPath(conf.RemoteUrl, "latest")
	if err != nil {
		return "", err
	}

	value, err := githubApiGetRequest(latestUrl, buildAuthorizationHeader(conf.Token))
	if err != nil {
		return "", err
	}

	version, ok := extractCleanVersion(value)
	if !ok {
		return "", errCast
	}
	return version, nil
}

func githubListReleases(conf *config.Config) ([]string, error) {
	basePageUrl := conf.RemoteUrl + "?page="
	authorizationHeader := buildAuthorizationHeader(conf.Token)

	page := 1
	var releases []string
	for {
		pageUrl := basePageUrl + strconv.Itoa(page)
		value, err := githubApiGetRequest(pageUrl, authorizationHeader)
		if err != nil {
			return nil, err
		}

		values, ok := value.([]any)
		if !ok {
			return nil, errCast
		}

		if len(values) == 0 {
			break
		}

		for _, value := range values {
			version, ok := extractCleanVersion(value)
			if !ok {
				return nil, errCast
			}
			releases = append(releases, version)
		}
		page++
	}
	return releases, nil
}

func githubApiGetRequest(callUrl string, authorizationHeader string) (any, error) {
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

func buildAssetName(version string) string {
	var nameBuilder strings.Builder
	nameBuilder.WriteString("tofu_")
	nameBuilder.WriteString(version)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOARCH)
	nameBuilder.WriteString(".zip")
	return nameBuilder.String()
}

func buildAuthorizationHeader(token string) string {
	if token == "" {
		return ""
	}

	var authorizationBuilder strings.Builder
	authorizationBuilder.WriteString("Bearer ")
	authorizationBuilder.WriteString(token)
	return authorizationBuilder.String()
}

func extractCleanVersion(value any) (string, bool) {
	object, _ := value.(map[string]any)
	tagName, _ := object["tag_name"].(string)
	if tagName == "" {
		return "", false
	}

	if tagName[0] == 'v' {
		tagName = tagName[1:]
	}
	return tagName, true
}
