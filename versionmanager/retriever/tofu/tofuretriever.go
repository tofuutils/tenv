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
	"github.com/tofuutils/tenv/pkg/zip"
)

const (
	defaultTofuGithubURL = "https://api.github.com/repos/opentofu/opentofu/releases"
	publicKeyURL         = "https://get.opentofu.org/opentofu.asc"

	baseIdentity = "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v"
	issuer       = "https://token.actions.githubusercontent.com"
)

const baseFileName = "tofu_"

type TofuRetriever struct {
	conf       *config.Config
	notLoaded  bool
	remoteConf map[string]string
}

func NewTofuRetriever(conf *config.Config) *TofuRetriever {
	return &TofuRetriever{conf: conf, notLoaded: true}
}

func (r *TofuRetriever) InstallRelease(versionStr string, targetPath string) error {
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
		return err
	}
	stable := v.Prerelease() == ""

	assetNames := buildAssetNames(versionStr, stable)
	assets, err := github.AssetDownloadURL(tag, assetNames, r.getRemoteURL(), r.conf.GithubToken)
	if err != nil {
		return err
	}

	urlTranformer := download.UrlTranformer(r.readRemoteConf())
	downloadURL, err := urlTranformer(assets[assetNames[0]])
	if err != nil {
		return err
	}

	data, err := download.Bytes(downloadURL)
	if err != nil {
		return err
	}

	if err = r.checkSumAndSig(v, stable, data, assetNames, assets, urlTranformer); err != nil {
		return err
	}

	return zip.UnzipToDir(data, targetPath)
}

func (r *TofuRetriever) LatestRelease() (string, error) {
	return github.LatestRelease(r.getRemoteURL(), r.conf.GithubToken)
}

func (r *TofuRetriever) ListReleases() ([]string, error) {
	return github.ListReleases(r.getRemoteURL(), r.conf.GithubToken)
}

func (r *TofuRetriever) checkSumAndSig(version *version.Version, stable bool, data []byte, assetNames []string, assets map[string]string, urlTranformer func(string) (string, error)) error {
	downloadURL, err := urlTranformer(assets[assetNames[1]])
	if err != nil {
		return err
	}

	dataSums, err := download.Bytes(downloadURL)
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, assetNames[0]); err != nil {
		return err
	}

	downloadURL, err = urlTranformer(assets[assetNames[3]])
	if err != nil {
		return err
	}

	dataSumsSig, err := download.Bytes(downloadURL)
	if err != nil {
		return err
	}

	downloadURL, err = urlTranformer(assets[assetNames[2]])
	if err != nil {
		return err
	}

	dataSumsCert, err := download.Bytes(downloadURL)
	if err != nil {
		return err
	}

	identity := buildIdentity(version)
	err = cosigncheck.Check(dataSums, dataSumsSig, dataSumsCert, identity, issuer)
	if err == nil || err != cosigncheck.ErrNotInstalled {
		return err
	}

	if !stable {
		fmt.Println("skip signature check : cosign executable not found and pgp check not available for unstable version") //nolint

		return nil
	}

	if r.conf.Verbose {
		fmt.Println("cosign executable not found, fallback to pgp check") //nolint
	}

	downloadURL, err = urlTranformer(assets[assetNames[4]])
	if err != nil {
		return err
	}

	dataSumsSig, err = download.Bytes(downloadURL)
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

func (r *TofuRetriever) getRemoteURL() string {
	if r.conf.TofuRemoteURL != "" {
		return r.conf.TofuRemoteURL
	}

	return config.MapGetDefault(r.readRemoteConf(), "url", defaultTofuGithubURL)
}

func (r *TofuRetriever) readRemoteConf() map[string]string {
	if r.notLoaded {
		r.notLoaded = false
		r.remoteConf = r.conf.ReadRemoteConf("tofu")
	}

	return r.remoteConf
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
