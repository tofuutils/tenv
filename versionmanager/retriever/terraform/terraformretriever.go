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
	"os"
	"runtime"
	"slices"

	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/pkg/apierrors"
	pgpcheck "github.com/dvaumoron/gotofuenv/pkg/check/pgp"
	sha256check "github.com/dvaumoron/gotofuenv/pkg/check/sha256"
	"github.com/dvaumoron/gotofuenv/pkg/download"
	"github.com/dvaumoron/gotofuenv/versionmanager/semantic"
)

const publicKeyUrl = "https://www.hashicorp.com/.well-known/pgp-key.txt"

const indexJson = "index.json"

type TerraformRetriever struct {
	conf *config.Config
}

func NewTerraformRetriever(conf *config.Config) *TerraformRetriever {
	return &TerraformRetriever{conf: conf}
}

func (r *TerraformRetriever) DownloadReleaseZip(version string) ([]byte, error) {
	// assume that terraform version do not start with a 'v'
	if version[0] == 'v' {
		version = version[1:]
	}

	baseVersionUrl, err := url.JoinPath(r.conf.TfRemoteUrl, version)
	if err != nil {
		return nil, err
	}

	versionUrl, err := url.JoinPath(baseVersionUrl, indexJson)
	if err != nil {
		return nil, err
	}

	value, err := apiGetRequest(versionUrl)
	if err != nil {
		return nil, err
	}

	object, _ := value.(map[string]any)
	builds, ok := object["builds"].([]any)
	if !ok {
		return nil, apierrors.ErrReturn
	}

	shaFileName, ok := object["shasums"].(string)
	if !ok {
		return nil, apierrors.ErrReturn
	}

	shaSigFileName, ok := object["shasums_signature"].(string)
	if !ok {
		return nil, apierrors.ErrReturn
	}

	downloadSumsUrl, err := url.JoinPath(baseVersionUrl, shaFileName)
	if err != nil {
		return nil, err
	}

	downloadSumsSigUrl, err := url.JoinPath(baseVersionUrl, shaSigFileName)
	if err != nil {
		return nil, err
	}

	for _, build := range builds {
		object, ok = build.(map[string]any)
		if !ok {
			return nil, apierrors.ErrReturn
		}

		osStr, ok := object["os"].(string)
		if !ok {
			return nil, apierrors.ErrReturn
		}

		if osStr != runtime.GOOS {
			continue
		}

		archStr, ok := object["arch"].(string)
		if !ok {
			return nil, apierrors.ErrReturn
		}

		if archStr != runtime.GOARCH {
			continue
		}

		downloadUrl, ok := object["url"].(string)
		if !ok {
			return nil, apierrors.ErrReturn
		}

		fileName, ok := object["filename"].(string)
		if !ok {
			return nil, apierrors.ErrReturn
		}

		data, err := download.DownloadBytes(downloadUrl)
		if err != nil {
			return nil, err
		}

		if err = r.checkSumAndSig(fileName, data, downloadSumsUrl, downloadSumsSigUrl); err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, apierrors.ErrAsset
}

func (r *TerraformRetriever) LatestRelease() (string, error) {
	// hashicorp release api does not seem to have a shortcut
	versions, err := r.ListReleases()
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

func (r *TerraformRetriever) ListReleases() ([]string, error) {
	releaseUrl, err := url.JoinPath(r.conf.TfRemoteUrl, indexJson)
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

func (r *TerraformRetriever) checkSumAndSig(fileName string, data []byte, downloadSumsUrl string, downloadSumsSigUrl string) error {
	dataSums, err := download.DownloadBytes(downloadSumsUrl)
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, fileName); err != nil {
		return err
	}

	dataSumsSig, err := download.DownloadBytes(downloadSumsSigUrl)
	if err != nil {
		return err
	}

	var dataPublicKey []byte
	if r.conf.TfKeyPath == "" {
		dataPublicKey, err = download.DownloadBytes(publicKeyUrl)
	} else {
		dataPublicKey, err = os.ReadFile(r.conf.TfKeyPath)
	}

	if err != nil {
		return err
	}
	return pgpcheck.Check(dataSums, dataSumsSig, dataPublicKey)
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
