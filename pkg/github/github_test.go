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
	_ "embed"
	"encoding/json"
	"slices"
	"testing"

	"github.com/dvaumoron/gotofuenv/pkg/apierrors"
	"github.com/dvaumoron/gotofuenv/versionmanager/semantic"
)

//go:embed assets.json
var assetsData []byte

var (
	assetsValue any
	assetsErr   error
)

//go:embed release.json
var releaseData []byte

var (
	releaseValue any
	releaseErr   error
)

//go:embed releases.json
var releasesData []byte

var (
	releasesValue any
	releasesErr   error
)

func init() {
	assetsErr = json.Unmarshal(assetsData, &assetsValue)
	releaseErr = json.Unmarshal(releaseData, &releaseValue)
	releasesErr = json.Unmarshal(releasesData, &releasesValue)
}

func TestExtractAssetsEmpty(t *testing.T) {
	e := struct{}{}
	assets := map[string]string{}
	searchedAssetNames := map[string]struct{}{"tofu_1.6.0_386.deb": e, "tofu_1.6.0_amd64.apk.gpgsig": e}
	err := extractAssets(assets, searchedAssetNames, 2, []any{})
	if err == nil {
		t.Error("Should fail on empty data")
	} else if err != apierrors.ErrAsset {
		t.Error("Unexpected extract error : ", err)
	}
}

func TestExtractAssetsMissing(t *testing.T) {
	if assetsErr != nil {
		t.Fatal("Unexpected parsing error : ", assetsErr)
		return
	}

	e := struct{}{}
	assets := map[string]string{}
	searchedAssetNames := map[string]struct{}{"tofu_1.6.0_386.deb": e, "any_name.zip": e}
	err := extractAssets(assets, searchedAssetNames, 2, assetsValue)
	if err == nil {
		t.Error("Should fail on non exiting fileName")
	} else if err != errContinue {
		t.Error("Unexpected extract error : ", err)
	}
}

func TestExtractAssetsPresent(t *testing.T) {
	if assetsErr != nil {
		t.Fatal("Unexpected parsing error : ", assetsErr)
		return
	}

	e := struct{}{}
	assets := map[string]string{}
	searchedAssetNames := map[string]struct{}{"tofu_1.6.0_386.deb": e, "tofu_1.6.0_amd64.apk.gpgsig": e}
	err := extractAssets(assets, searchedAssetNames, 2, assetsValue)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if res1 := assets["tofu_1.6.0_386.deb"]; res1 != "https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_386.deb" {
		t.Error("Unmatching result 1, get :", res1)
	}
	if res2 := assets["tofu_1.6.0_amd64.apk.gpgsig"]; res2 != "https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_amd64.apk.gpgsig" {
		t.Error("Unmatching result 2, get :", res2)
	}
}

func TestExtractReleasesEmpty(t *testing.T) {
	releases, err := extractReleases([]string{"value"}, []any{})
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	size := len(releases)
	if size == 0 {
		t.Fatal("Unexpected empty results")
	}

	if releases[0] != "value" || size > 1 {
		t.Error("Unexpected result :", releases)
	}
}

func TestExtractReleasesPresent(t *testing.T) {
	if releasesErr != nil {
		t.Fatal("Unexpected parsing error : ", releasesErr)
		return
	}

	var releases []string
	releases, err := extractReleases(releases, releasesValue)
	if err == nil {
		t.Fatal("Should return a errContinue marker")
	} else if err != errContinue {
		t.Fatal("Unexpected extract error : ", err)
	}

	slices.SortFunc(releases, semantic.CmpVersion)
	if !slices.Equal(releases, []string{"1.6.0-alpha5", "1.6.0-beta5", "1.6.0-rc1", "1.6.0"}) {
		t.Error("Unmatching results, get :", releases)
	}
}

func TestExtractVersion(t *testing.T) {
	if releaseErr != nil {
		t.Fatal("Unexpected parsing error : ", releaseErr)
		return
	}

	version, ok := extractVersion(releaseValue)
	if !ok {
		t.Fatal("Unexpected extract failure")
	}
	if version != "1.6.0" {
		t.Error("Unmatching result, get :", version)
	}
}
