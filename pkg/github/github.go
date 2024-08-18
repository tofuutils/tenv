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
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tofuutils/tenv/v3/pkg/apimsg"
	"github.com/tofuutils/tenv/v3/pkg/download"
	versionfinder "github.com/tofuutils/tenv/v3/versionmanager/semantic/finder"
)

const (
	Download = "download"
	Releases = "releases"

	pageQuery = "?page="
)

var errContinue = errors.New("continue")

func AssetDownloadURL(tag string, searchedAssetNames []string, githubReleaseURL string, githubToken string, display func(string)) ([]string, error) {
	releaseUrl, err := url.JoinPath(githubReleaseURL, "tags", tag) //nolint
	if err != nil {
		return nil, err
	}

	display(apimsg.MsgFetchRelease + releaseUrl)

	authorizationHeader := buildAuthorizationHeader(githubToken)
	value, err := apiGetRequest(releaseUrl, authorizationHeader)
	if err != nil {
		return nil, err
	}

	object, _ := value.(map[string]any)
	baseAssetsURL, ok := object["assets_url"].(string)
	if !ok {
		return nil, apimsg.ErrReturn
	}

	waited := len(searchedAssetNames)
	searchedAssetNameSet := make(map[string]struct{}, waited)
	for _, searchAssetName := range searchedAssetNames {
		searchedAssetNameSet[searchAssetName] = struct{}{}
	}

	page := 1
	assets := make(map[string]string, waited)
	baseAssetsURL += pageQuery
	for {
		assetsURL := baseAssetsURL + strconv.Itoa(page)
		value, err = apiGetRequest(assetsURL, authorizationHeader)
		if err != nil {
			return nil, err
		}

		if err = extractAssets(assets, searchedAssetNameSet, waited, value); err == nil {
			assetURLs := make([]string, 0, waited)
			for _, searchAssetName := range searchedAssetNames {
				assetURLs = append(assetURLs, assets[searchAssetName])
			}

			return assetURLs, nil
		} else if err != errContinue {
			return nil, err
		}
		page++
	}
}

func ListReleases(githubReleaseURL string, githubToken string) ([]string, error) {
	basePageURL := githubReleaseURL + pageQuery
	authorizationHeader := buildAuthorizationHeader(githubToken)

	page := 1
	var releases []string
	for {
		pageURL := basePageURL + strconv.Itoa(page)
		value, err := apiGetRequest(pageURL, authorizationHeader)
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

func apiGetRequest(callURL string, authorizationHeader string) (any, error) {
	return download.JSON(callURL, download.NoDisplay, func(request *http.Request) *http.Request {
		request.Header.Set("Accept", "application/vnd.github+json")
		if authorizationHeader != "" {
			request.Header.Set("Authorization", authorizationHeader)
		}
		request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		return request
	})
}

func buildAuthorizationHeader(token string) string {
	if token == "" {
		return ""
	}

	return "Bearer " + token
}

func extractAssets(assets map[string]string, searchedAssetNameSet map[string]struct{}, waited int, value any) error {
	values, ok := value.([]any)
	if !ok {
		return apimsg.ErrReturn
	}

	if len(values) == 0 {
		return apimsg.ErrAsset
	}

	for _, value := range values {
		object, _ := value.(map[string]any)
		assetName, ok := object["name"].(string) //nolint
		if !ok {
			return apimsg.ErrReturn
		}

		if _, ok := searchedAssetNameSet[assetName]; !ok {
			continue
		}

		downloadURL, ok := object["browser_download_url"].(string)
		if !ok {
			return apimsg.ErrReturn
		}
		assets[assetName] = downloadURL

		if len(assets) == waited {
			return nil
		}
	}

	return errContinue
}

func extractReleases(releases []string, value any) ([]string, error) {
	values, ok := value.([]any)
	if !ok {
		return nil, apimsg.ErrReturn
	}

	if len(values) == 0 {
		return releases, nil
	}

	for _, value := range values {
		version := extractVersion(value)
		if version == "" {
			return nil, apimsg.ErrReturn
		}
		releases = append(releases, version)
	}

	return releases, errContinue
}

func extractVersion(value any) string {
	object, _ := value.(map[string]any)
	version, _ := object["tag_name"].(string)

	return versionfinder.Find(version)
}
