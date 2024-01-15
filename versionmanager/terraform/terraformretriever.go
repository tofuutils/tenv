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

package terraformretriever

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"slices"

	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/pkg/apierrors"
	"github.com/dvaumoron/gotofuenv/versionmanager/semantic"
)

const indexJson = "index.json"

type TerraformRetriever struct {
	conf *config.Config
}

func MakeTerraformRetriever(conf *config.Config) TerraformRetriever {
	return TerraformRetriever{conf: conf}
}

func (v TerraformRetriever) DownloadAssetsUrl(version string) (string, string, error) {
	// assume that terraform version do not start with a 'v'
	if version[0] == 'v' {
		version = version[1:]
	}

	versionUrl, err := url.JoinPath(v.conf.TfRemoteUrl, version, indexJson)
	if err != nil {
		return "", "", err
	}

	value, err := apiGetRequest(versionUrl)
	if err != nil {
		return "", "", err
	}

	object, _ := value.(map[string]any)
	builds, ok := object["builds"].([]any)
	if !ok {
		return "", "", apierrors.ErrReturn
	}

	for _, build := range builds {
		object, ok = build.(map[string]any)
		if !ok {
			return "", "", apierrors.ErrReturn
		}

		osStr, ok := object["os"].(string)
		if !ok {
			return "", "", apierrors.ErrReturn
		}

		if osStr != runtime.GOOS {
			continue
		}

		archStr, ok := object["arch"].(string)
		if !ok {
			return "", "", apierrors.ErrReturn
		}

		if archStr != runtime.GOARCH {
			continue
		}

		downloadUrl, ok := object["url"].(string)
		if !ok {
			return "", "", apierrors.ErrReturn
		}
		return downloadUrl, "todo", nil
	}
	return "", "", apierrors.ErrNoAsset
}

func (v TerraformRetriever) LatestRelease() (string, error) {
	// hashicorp release api does not seem to have a shortcut
	versions, err := v.ListReleases()
	if err != nil {
		return "", err
	}

	versionLen := len(versions)
	if versionLen == 0 {
		return "", apierrors.ErrReturn
	}

	slices.SortFunc(versions, semantic.CmpVersion)
	return versions[versionLen-1], nil
}

func (v TerraformRetriever) ListReleases() ([]string, error) {
	releaseUrl, err := url.JoinPath(v.conf.TfRemoteUrl, indexJson)
	if err != nil {
		return nil, err
	}

	value, err := apiGetRequest(releaseUrl)
	if err != nil {
		return nil, err
	}

	object, ok := value.(map[string]any)
	if !ok {
		return nil, apierrors.ErrReturn
	}

	object, ok = object["versions"].(map[string]any)
	if !ok {
		return nil, apierrors.ErrReturn
	}

	var releases []string
	for version := range object {
		releases = append(releases, version)
	}
	return releases, nil
}

func apiGetRequest(callUrl string) (any, error) {
	response, err := http.Get(callUrl)
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
