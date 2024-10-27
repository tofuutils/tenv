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

package releaseapi_test

import (
	_ "embed"
	"encoding/json"
	"slices"
	"testing"

	releaseapi "github.com/tofuutils/tenv/v3/versionmanager/retriever/terraform/api"
	"github.com/tofuutils/tenv/v3/versionmanager/semantic"
)

//go:embed testdata/release.json
var releaseData []byte

//go:embed testdata/releases.json
var releasesData []byte

func TestExtractAssetUrls(t *testing.T) {
	t.Parallel()

	var releaseValue any
	errRelease := json.Unmarshal(releaseData, &releaseValue)
	if errRelease != nil {
		t.Fatal("Unexpected parsing error : ", errRelease)
	}

	fileName, downloadURL, shaFileName, shaSigFileName, err := releaseapi.ExtractAssetURLs("linux", "386", releaseValue)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if fileName != "terraform_1.7.0_linux_386.zip" {
		t.Error("Unexpected fileName, get :", fileName)
	}
	if downloadURL != "https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_linux_386.zip" {
		t.Error("Unexpected downloadURL, get :", downloadURL)
	}
	if shaFileName != "terraform_1.7.0_SHA256SUMS" {
		t.Error("Unexpected downloadSumsURL, get :", shaFileName)
	}
	if shaSigFileName != "terraform_1.7.0_SHA256SUMS.sig" {
		t.Error("Unexpected downloadSumsSigURL, get :", shaSigFileName)
	}
}

func TestExtractReleases(t *testing.T) {
	t.Parallel()

	var releasesValue any
	err := json.Unmarshal(releasesData, &releasesValue)
	if err != nil {
		t.Fatal("Unexpected parsing error : ", err)
	}

	releases, err := releaseapi.ExtractReleases(releasesValue)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	slices.SortFunc(releases, semantic.CmpVersion)
	if !slices.Equal(releases, []string{"1.6.6", "1.7.0-rc1", "1.7.0-rc2", "1.7.0"}) {
		t.Error("Unmatching results, get :", releases)
	}
}
