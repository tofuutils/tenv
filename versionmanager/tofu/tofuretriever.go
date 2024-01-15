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

package tofuretriever

import (
	"runtime"
	"strings"

	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/pkg/github"
)

type TofuRetriever struct {
	conf *config.Config
}

func MakeTofuRetriever(conf *config.Config) TofuRetriever {
	return TofuRetriever{conf: conf}
}

func (v TofuRetriever) DownloadAssetsUrl(version string) (string, string, error) {
	tag := version
	// assume that opentofu tags start with a 'v'
	// and version in asset name does not
	if tag[0] == 'v' {
		version = version[1:]
	} else {
		tag = "v" + version
	}

	assetNames := buildAssetNames(version)
	assets, err := github.DownloadAssetUrl(tag, assetNames, v.conf.TofuRemoteUrl, v.conf.GithubToken)
	if err != nil {
		return "", "", nil
	}

	// should be safe here (an error would have been returned if one was not found)
	return assets[assetNames[0]], assets[assetNames[1]], nil
}

func (v TofuRetriever) LatestRelease() (string, error) {
	return github.LatestRelease(v.conf.TofuRemoteUrl, v.conf.GithubToken)
}

func (v TofuRetriever) ListReleases() ([]string, error) {
	return github.ListReleases(v.conf.TofuRemoteUrl, v.conf.GithubToken)
}

func buildAssetNames(version string) []string {
	var nameBuilder strings.Builder
	nameBuilder.WriteString("tofu_")
	nameBuilder.WriteString(version)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOARCH)
	nameBuilder.WriteString(".zip")
	zipAssetName := nameBuilder.String()

	nameBuilder.WriteString(".gpgsig")
	return []string{zipAssetName, nameBuilder.String()}
}
