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

package tofuretriever

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/tofuutils/tenv/config"
	cosigncheck "github.com/tofuutils/tenv/pkg/check/cosign"
	pgpcheck "github.com/tofuutils/tenv/pkg/check/pgp"
	sha256check "github.com/tofuutils/tenv/pkg/check/sha256"
	"github.com/tofuutils/tenv/pkg/download"
	"github.com/tofuutils/tenv/pkg/github"
)

const publicKeyURL = "https://get.opentofu.org/opentofu.asc"

const (
	baseIdentity = "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v"
	issuer       = "https://token.actions.githubusercontent.com"
)

const baseFileName = "tofu_"

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

	v, err := version.NewVersion(versionStr) //nolint
	if err != nil {
		return nil, err
	}
	stable := v.Prerelease() == ""

	assetNames := buildAssetNames(versionStr, stable)
	assets, err := github.DownloadAssetURL(tag, assetNames, r.conf.TofuRemoteURL, r.conf.GithubToken)
	if err != nil {
		return nil, err
	}

	data, err := download.Bytes(assets[assetNames[0]])
	if err != nil {
		return nil, err
	}

	if err = r.checkSumAndSig(v, stable, data, assetNames, assets); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *TofuRetriever) LatestRelease() (string, error) {
	return github.LatestRelease(r.conf.TofuRemoteURL, r.conf.GithubToken)
}

func (r *TofuRetriever) ListReleases() ([]string, error) {
	return github.ListReleases(r.conf.TofuRemoteURL, r.conf.GithubToken)
}

func (r *TofuRetriever) checkSumAndSig(version *version.Version, stable bool, data []byte, assetNames []string, assets map[string]string) error {
	dataSums, err := download.Bytes(assets[assetNames[1]])
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, assetNames[0]); err != nil {
		return err
	}

	dataSumsSig, err := download.Bytes(assets[assetNames[3]])
	if err != nil {
		return err
	}

	dataSumsCert, err := download.Bytes(assets[assetNames[2]])
	if err != nil {
		return err
	}

	identity := buildIdentity(version)
	err = cosigncheck.Check(dataSums, dataSumsSig, dataSumsCert, identity, issuer)
	if err == nil || err != cosigncheck.ErrNotInstalled {
		return err
	}

	if stable {
		if r.conf.Verbose {
			fmt.Println("cosign executable not found, fallback to pgp check") //nolint
		}
	} else {
		fmt.Println("skip signature check : cosign executable not found and pgp check not available for unstable version") //nolint

		return nil
	}

	dataSumsSig, err = download.Bytes(assets[assetNames[4]])
	if err != nil {
		return err
	}

	var dataPublicKey []byte
	if r.conf.TofuKeyPath == "" {
		dataPublicKey, err = download.Bytes(publicKeyURL)
	} else {
		dataPublicKey, err = os.ReadFile(r.conf.TofuKeyPath)
	}

	if err != nil {
		return err
	}

	return pgpcheck.Check(dataSums, dataSumsSig, dataPublicKey)
}

func buildAssetNames(version string, stable bool) []string {
	var nameBuilder strings.Builder
	nameBuilder.WriteString(baseFileName)
	nameBuilder.WriteString(version)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOARCH)
	nameBuilder.WriteString(".zip")
	zipAssetName := nameBuilder.String()

	nameBuilder.Reset()
	nameBuilder.WriteString(baseFileName)
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
