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
	pgpcheck "github.com/dvaumoron/gotofuenv/pkg/check/pgp"
	sha256check "github.com/dvaumoron/gotofuenv/pkg/check/sha256"
	"github.com/dvaumoron/gotofuenv/pkg/download"
	"github.com/dvaumoron/gotofuenv/pkg/github"
	"github.com/dvaumoron/gotofuenv/versionmanager/semantic"
)

const publicKeyUrl = "https://get.opentofu.org/opentofu.asc"

type TofuRetriever struct {
	conf *config.Config
}

func NewTofuRetriever(conf *config.Config) *TofuRetriever {
	return &TofuRetriever{conf: conf}
}

func (r *TofuRetriever) DownloadReleaseZip(version string) ([]byte, error) {
	tag := version
	// assume that opentofu tags start with a 'v'
	// and version in asset name does not
	if tag[0] == 'v' {
		version = version[1:]
	} else {
		tag = "v" + version
	}

	stable := semantic.StableVersion(version)
	assetNames := buildAssetNames(version, stable)
	assets, err := github.DownloadAssetUrl(tag, assetNames, r.conf.TofuRemoteUrl, r.conf.GithubToken)
	if err != nil {
		return nil, err
	}

	data, err := download.DownloadBytes(assets[assetNames[0]])
	if err != nil {
		return nil, err
	}

	dataSums, err := download.DownloadBytes(assets[assetNames[1]])
	if err != nil {
		return nil, err
	}

	dataSum, err := sha256check.Extract(dataSums, assetNames[0])
	if err != nil {
		return nil, err
	}

	if err = sha256check.Check(data, dataSum); err != nil {
		return nil, err
	}

	if stable {
		dataSig, err := download.DownloadBytes(assets[assetNames[2]])
		if err != nil {
			return nil, err
		}

		dataPublicKey, err := download.DownloadBytes(publicKeyUrl)
		if err != nil {
			return nil, err
		}

		err = pgpcheck.Check(data, dataSig, dataPublicKey)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (r *TofuRetriever) LatestRelease() (string, error) {
	return github.LatestRelease(r.conf.TofuRemoteUrl, r.conf.GithubToken)
}

func (r *TofuRetriever) ListReleases() ([]string, error) {
	return github.ListReleases(r.conf.TofuRemoteUrl, r.conf.GithubToken)
}

func buildAssetNames(version string, stable bool) []string {
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
	pgpAssetName := nameBuilder.String()

	nameBuilder.Reset()
	nameBuilder.WriteString("tofu_")
	nameBuilder.WriteString(version)
	nameBuilder.WriteString("_SHA256SUMS")

	if stable {
		return []string{zipAssetName, nameBuilder.String(), pgpAssetName}
	}
	return []string{zipAssetName, nameBuilder.String()}
}
