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
	"github.com/dvaumoron/gotofuenv/pkg/sha256check"
)

type TofuRetriever struct {
	assetNames []string
	conf       *config.Config
}

func NewTofuRetriever(conf *config.Config) *TofuRetriever {
	return &TofuRetriever{conf: conf}
}

func (r *TofuRetriever) Check(data []byte, dataSigs []byte) error {
	dataSig, err := sha256check.Extract(dataSigs, r.assetNames[0])
	if err != nil {
		return err
	}
	return sha256check.Check(data, dataSig)
}

func (r *TofuRetriever) DownloadAssetsUrl(version string) (string, string, error) {
	tag := version
	// assume that opentofu tags start with a 'v'
	// and version in asset name does not
	if tag[0] == 'v' {
		version = version[1:]
	} else {
		tag = "v" + version
	}

	r.assetNames = buildAssetNames(version)
	assets, err := github.DownloadAssetUrl(tag, r.assetNames, r.conf.TofuRemoteUrl, r.conf.GithubToken)
	if err != nil {
		return "", "", nil
	}

	// should be safe here (an error would have been returned if one was not found)
	return assets[r.assetNames[0]], assets[r.assetNames[1]], nil
}

func (r *TofuRetriever) LatestRelease() (string, error) {
	return github.LatestRelease(r.conf.TofuRemoteUrl, r.conf.GithubToken)
}

func (r *TofuRetriever) ListReleases() ([]string, error) {
	return github.ListReleases(r.conf.TofuRemoteUrl, r.conf.GithubToken)
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

	nameBuilder.Reset()
	nameBuilder.WriteString("tofu_")
	nameBuilder.WriteString(version)
	nameBuilder.WriteString("_SHA256SUMS")
	return []string{zipAssetName, nameBuilder.String()}
}
