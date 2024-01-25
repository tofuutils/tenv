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

package github

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tofuutils/tenv/pkg/apierrors"
)

const pageQuery = "?page="

var errContinue = errors.New("continue")

func DownloadAssetUrl(tag string, searchedAssetNames []string, githubReleaseUrl string, githubToken string) (map[string]string, error) {
	releaseUrl, err := url.JoinPath(githubReleaseUrl, "tags", tag) //nolint
	if err != nil {
		return nil, err
	}

	authorizationHeader := buildAuthorizationHeader(githubToken)
	value, err := apiGetRequest(releaseUrl, authorizationHeader)
	if err != nil {
		return nil, err
	}

	object, _ := value.(map[string]any)
	baseAssetsUrl, ok := object["assets_url"].(string)
	if !ok {
		return nil, apierrors.ErrReturn
	}

	waited := len(searchedAssetNames)
	searchedAssetNameSet := make(map[string]struct{}, waited)
	for _, searchAssetName := range searchedAssetNames {
		searchedAssetNameSet[searchAssetName] = struct{}{}
	}

	page := 1
	assets := make(map[string]string, waited)
	baseAssetsUrl += pageQuery
	for {
		assetsUrl := baseAssetsUrl + strconv.Itoa(page)
		value, err = apiGetRequest(assetsUrl, authorizationHeader)
		if err != nil {
			return nil, err
		}

		if err = extractAssets(assets, searchedAssetNameSet, waited, value); err == nil {
			return assets, nil
		} else if err != errContinue {
			return nil, err
		}
		page++
	}
}

func LatestRelease(githubReleaseUrl string, githubToken string) (string, error) {
	latestUrl, err := url.JoinPath(githubReleaseUrl, "latest") //nolint
	if err != nil {
		return "", err
	}

	authorizationHeader := buildAuthorizationHeader(githubToken)
	value, err := apiGetRequest(latestUrl, authorizationHeader)
	if err != nil {
		return "", err
	}

	version, ok := extractVersion(value)
	if !ok {
		return "", apierrors.ErrReturn
	}

	return version, nil
}

func ListReleases(githubReleaseUrl string, githubToken string) ([]string, error) {
	basePageUrl := githubReleaseUrl + pageQuery
	authorizationHeader := buildAuthorizationHeader(githubToken)

	page := 1
	var releases []string
	for {
		pageUrl := basePageUrl + strconv.Itoa(page)
		value, err := apiGetRequest(pageUrl, authorizationHeader)
		if err != nil {
			return nil, err
		}

		releases, err = extractReleases(releases, value)
		if err == nil {
			return releases, nil
		} else if err != errContinue {
			return nil, err
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

func buildAuthorizationHeader(token string) string {
	if token == "" {
		return ""
	}

	var authorizationBuilder strings.Builder
	authorizationBuilder.WriteString("Bearer ")
	authorizationBuilder.WriteString(token)

	return authorizationBuilder.String()
}

func extractAssets(assets map[string]string, searchedAssetNameSet map[string]struct{}, waited int, value any) error {
	values, ok := value.([]any)
	if !ok {
		return apierrors.ErrReturn
	}

	if len(values) == 0 {
		return apierrors.ErrAsset
	}

	for _, value := range values {
		object, _ := value.(map[string]any)
		assetName, ok := object["name"].(string)
		if !ok {
			return apierrors.ErrReturn
		}

		if _, ok := searchedAssetNameSet[assetName]; !ok {
			continue
		}

		downloadUrl, ok := object["browser_download_url"].(string)
		if !ok {
			return apierrors.ErrReturn
		}
		assets[assetName] = downloadUrl

		if len(assets) == waited {
			return nil
		}
	}

	return errContinue
}

func extractReleases(releases []string, value any) ([]string, error) {
	values, ok := value.([]any)
	if !ok {
		return nil, apierrors.ErrReturn
	}

	if len(values) == 0 {
		return releases, nil
	}

	for _, value := range values {
		version, ok := extractVersion(value)
		if !ok {
			return nil, apierrors.ErrReturn
		}
		releases = append(releases, version)
	}

	return releases, errContinue
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
