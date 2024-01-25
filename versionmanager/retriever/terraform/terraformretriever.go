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

package terraformretriever

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"slices"

	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/apierrors"
	pgpcheck "github.com/tofuutils/tenv/pkg/check/pgp"
	sha256check "github.com/tofuutils/tenv/pkg/check/sha256"
	"github.com/tofuutils/tenv/pkg/download"
	"github.com/tofuutils/tenv/versionmanager/semantic"
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

	baseVersionUrl, err := url.JoinPath(r.conf.TfRemoteURL, version)
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

	fileName, downloadUrl, downloadSumsUrl, downloadSumsSigUrl, err := extractAssetUrls(baseVersionUrl, runtime.GOOS, runtime.GOARCH, value)
	if err != nil {
		return nil, err
	}

	data, err := download.Bytes(downloadUrl)
	if err != nil {
		return nil, err
	}

	if err = r.checkSumAndSig(fileName, data, downloadSumsUrl, downloadSumsSigUrl); err != nil {
		return nil, err
	}
	return data, nil
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
	releaseUrl, err := url.JoinPath(r.conf.TfRemoteURL, indexJson)
	if err != nil {
		return nil, err
	}

	value, err := apiGetRequest(releaseUrl)
	if err != nil {
		return nil, err
	}
	return extractReleases(value)
}

func (r *TerraformRetriever) checkSumAndSig(fileName string, data []byte, downloadSumsUrl string, downloadSumsSigUrl string) error {
	dataSums, err := download.Bytes(downloadSumsUrl)
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, fileName); err != nil {
		return err
	}

	dataSumsSig, err := download.Bytes(downloadSumsSigUrl)
	if err != nil {
		return err
	}

	var dataPublicKey []byte
	if r.conf.TfKeyPath == "" {
		dataPublicKey, err = download.Bytes(publicKeyUrl)
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

func extractAssetUrls(baseVersionUrl string, searchedOs string, searchedArch string, value any) (string, string, string, string, error) {
	object, _ := value.(map[string]any)
	builds, ok := object["builds"].([]any)
	shaFileName, ok2 := object["shasums"].(string)
	shaSigFileName, ok3 := object["shasums_signature"].(string)
	if !ok || !ok2 || !ok3 {
		return "", "", "", "", apierrors.ErrReturn
	}

	downloadSumsUrl, err := url.JoinPath(baseVersionUrl, shaFileName)
	if err != nil {
		return "", "", "", "", err
	}

	downloadSumsSigUrl, err := url.JoinPath(baseVersionUrl, shaSigFileName)
	if err != nil {
		return "", "", "", "", err
	}

	for _, build := range builds {
		object, _ = build.(map[string]any)
		osStr, ok := object["os"].(string)
		archStr, ok2 := object["arch"].(string)
		downloadUrl, ok3 := object["url"].(string)
		fileName, ok4 := object["filename"].(string)
		if !ok || !ok2 || !ok3 || !ok4 {
			return "", "", "", "", apierrors.ErrReturn
		}

		if osStr != searchedOs || archStr != searchedArch {
			continue
		}
		return fileName, downloadUrl, downloadSumsUrl, downloadSumsSigUrl, nil
	}
	return "", "", "", "", apierrors.ErrAsset
}

func extractReleases(value any) ([]string, error) {
	object, _ := value.(map[string]any)
	object, ok := object["versions"].(map[string]any)
	if !ok {
		return nil, apierrors.ErrReturn
	}

	var releases []string
	for version := range object {
		releases = append(releases, version)
	}
	return releases, nil
}
