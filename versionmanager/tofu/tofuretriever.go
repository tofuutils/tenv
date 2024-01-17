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
	"fmt"
	"runtime"
	"strings"

	"github.com/dvaumoron/gotofuenv/config"
	cosigncheck "github.com/dvaumoron/gotofuenv/pkg/check/cosign"
	pgpcheck "github.com/dvaumoron/gotofuenv/pkg/check/pgp"
	sha256check "github.com/dvaumoron/gotofuenv/pkg/check/sha256"
	"github.com/dvaumoron/gotofuenv/pkg/download"
	"github.com/dvaumoron/gotofuenv/pkg/github"
	"github.com/hashicorp/go-version"
)

const publicKeyUrl = "https://get.opentofu.org/opentofu.asc"

const (
	baseIdentity = "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v"
	issuer       = "https://token.actions.githubusercontent.com"
)

type TofuRetriever struct {
	conf *config.Config
}

func NewTofuRetriever(conf *config.Config) *TofuRetriever {
	return &TofuRetriever{conf: conf}
}

func (r *TofuRetriever) DownloadReleaseZip(versionStr string) ([]byte, error) {
	tag := versionStr
	// assume that opentofu tags start with a 'v'
	// and version in asset name does not
	if tag[0] == 'v' {
		versionStr = versionStr[1:]
	} else {
		tag = "v" + versionStr
	}

	v, err := version.NewVersion(versionStr)
	if err != nil {
		return nil, err
	}
	stable := v.Prerelease() == ""

	assetNames := buildAssetNames(versionStr, stable)
	assets, err := github.DownloadAssetUrl(tag, assetNames, r.conf.TofuRemoteUrl, r.conf.GithubToken)
	if err != nil {
		return nil, err
	}

	data, err := download.DownloadBytes(assets[assetNames[0]])
	if err != nil {
		return nil, err
	}

	if err = checkSumAndSig(v, stable, data, assetNames, assets, r.conf.Verbose); err != nil {
		return nil, err
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

	nameBuilder.Reset()
	nameBuilder.WriteString("tofu_")
	nameBuilder.WriteString(version)
	nameBuilder.WriteString("_SHA256SUMS")
	sumsAssetName := nameBuilder.String()

	if stable {
		return []string{zipAssetName, sumsAssetName, sumsAssetName + ".pem", sumsAssetName + ".sig", sumsAssetName + ".gpgsig"}
	}
	return []string{zipAssetName, sumsAssetName, sumsAssetName + ".pem", sumsAssetName + ".sig"}
}

func buildIdentity(v *version.Version) string {
	cleanedVersion := v.String()
	indexDot := strings.LastIndexByte(cleanedVersion, '.')
	// cleaned, so indexDot can not be -1
	shortVersion := cleanedVersion[:indexDot]
	return baseIdentity + shortVersion
}

func checkSumAndSig(v *version.Version, stable bool, data []byte, assetNames []string, assets map[string]string, verbose bool) error {
	dataSums, err := download.DownloadBytes(assetNames[1])
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, assetNames[0]); err != nil {
		return err
	}

	dataSumsSig, err := download.DownloadBytes(assetNames[3])
	if err != nil {
		return err
	}

	dataSumsCert, err := download.DownloadBytes(assetNames[2])
	if err != nil {
		return err
	}

	identity := buildIdentity(v)
	err = cosigncheck.Check(dataSums, dataSumsSig, dataSumsCert, identity, issuer)
	if err == nil || err != cosigncheck.ErrNotInstalled {
		return err
	}

	if stable {
		if verbose {
			fmt.Println("cosign executable not found, fallback to pgp check")
		}
	} else {
		fmt.Println("cosign executable not found and pgp check not available for unstable version")
		return nil
	}

	dataSumsSig, err = download.DownloadBytes(assetNames[4])
	if err != nil {
		return err
	}

	dataPublicKey, err := download.DownloadBytes(publicKeyUrl)
	if err != nil {
		return err
	}
	return pgpcheck.Check(dataSums, dataSumsSig, dataPublicKey)
}
