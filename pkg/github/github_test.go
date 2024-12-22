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
	_ "embed"
	"encoding/json"
	"errors"
	"slices"
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/apimsg"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic"
)

//go:embed testdata/assets.json
var assetsData []byte

//go:embed testdata/release.json
var releaseData []byte

//go:embed testdata/releases.json
var releasesData []byte

func TestExtractAssets(t *testing.T) {
	t.Parallel()

	var assetsValue any
	err := json.Unmarshal(assetsData, &assetsValue)
	if err != nil {
		t.Fatal("Unexpected parsing error : ", err)
	}

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		assets := map[string]string{}
		searchedAssetNames := map[string]struct{}{"tofu_1.6.0_386.deb": {}, "tofu_1.6.0_amd64.apk.gpgsig": {}}
		err = extractAssets(assets, searchedAssetNames, 2, []any{})
		if err == nil {
			t.Error("Should fail on empty data")
		} else if !errors.Is(err, apimsg.ErrAsset) {
			t.Error("Unexpected extract error : ", err)
		}
	})

	t.Run("Missing", func(t *testing.T) {
		t.Parallel()

		assets := map[string]string{}
		searchedAssetNames := map[string]struct{}{"tofu_1.6.0_386.deb": {}, "any_name.zip": {}}
		err = extractAssets(assets, searchedAssetNames, 2, assetsValue)
		if err == nil {
			t.Error("Should fail on non exiting fileName")
		} else if !errors.Is(err, errContinue) {
			t.Error("Unexpected extract error : ", err)
		}
	})

	t.Run("Present", func(t *testing.T) {
		t.Parallel()

		assets := map[string]string{}
		searchedAssetNames := map[string]struct{}{"tofu_1.6.0_386.deb": {}, "tofu_1.6.0_amd64.apk.gpgsig": {}}
		err = extractAssets(assets, searchedAssetNames, 2, assetsValue)
		if err != nil {
			t.Fatal("Unexpected extract error : ", err)
		}

		if res1 := assets["tofu_1.6.0_386.deb"]; res1 != "https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_386.deb" {
			t.Error("Unmatching result 1, get :", res1)
		}
		if res2 := assets["tofu_1.6.0_amd64.apk.gpgsig"]; res2 != "https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_amd64.apk.gpgsig" {
			t.Error("Unmatching result 2, get :", res2)
		}
	})
}

func TestExtractReleases(t *testing.T) {
	t.Parallel()

	var releasesValue any
	err := json.Unmarshal(releasesData, &releasesValue)
	if err != nil {
		t.Fatal("Unexpected parsing error : ", err)
	}

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

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
	})

	t.Run("Present", func(t *testing.T) {
		t.Parallel()

		var releases []string
		releases, err = extractReleases(releases, releasesValue)
		if err == nil {
			t.Fatal("Should return a errContinue marker")
		} else if !errors.Is(err, errContinue) {
			t.Fatal("Unexpected extract error : ", err)
		}

		slices.SortFunc(releases, semantic.CmpVersion)
		if !slices.Equal(releases, []string{"1.6.0-alpha5", "1.6.0-beta5", "1.6.0-rc1", "1.6.0"}) {
			t.Error("Unmatching results, get :", releases)
		}
	})
}

func TestExtractVersion(t *testing.T) {
	t.Parallel()

	var releaseValue any
	err := json.Unmarshal(releaseData, &releaseValue)
	if err != nil {
		t.Fatal("Unexpected parsing error : ", err)
	}

	version := extractVersion(releaseValue)
	if version == "" {
		t.Fatal("Unexpected empty result")
	}
	if version != "1.6.0" {
		t.Error("Unmatching result, get :", version)
	}
}
