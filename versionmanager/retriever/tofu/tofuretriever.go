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
	"net/url"
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
	htmlretriever "github.com/tofuutils/tenv/versionmanager/retriever/html"
)

const (
	publicKeyURL = "https://get.opentofu.org/opentofu.asc"

	baseIdentity = "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v"
	baseFileName = "tofu_"
	issuer       = "https://token.actions.githubusercontent.com"
	opentofu     = "opentofu"
)

type TofuRetriever struct {
	conf *config.Config
}

func NewTofuRetriever(conf *config.Config) *TofuRetriever {
	return &TofuRetriever{conf: conf}
}

func (r *TofuRetriever) InstallRelease(versionStr string, targetPath string) error {
	r.conf.InitRemoteConf()

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

	var assetURLs []string
	assetNames := buildAssetNames(versionStr, stable)
	if r.conf.Tofu.GetInstallMode() == htmlretriever.InstallModeDirect {
		baseAssetURL, err2 := url.JoinPath(r.conf.Tofu.GetRemoteURL(), opentofu, opentofu, github.Releases, github.Download, tag) //nolint
		if err2 != nil {
			return err2
		}

		assetURLs, err = htmlretriever.BuildAssetURLs(baseAssetURL, assetNames...)
	} else {
		assetURLs, err = github.AssetDownloadURL(tag, assetNames, r.conf.Tofu.GetRemoteURL(), r.conf.GithubToken, r.conf.DisplayVerbose)
	}
	if err != nil {
		return err
	}

	urlTranformer := download.UrlTranformer(r.conf.Tofu.GetRewriteRule())
	assetURLs, err = download.ApplyUrlTranformer(urlTranformer, assetURLs...)
	if err != nil {
		return err
	}

	data, err := download.Bytes(assetURLs[0], r.conf.DisplayVerbose)
	if err != nil {
		return err
	}

	if err = r.checkSumAndSig(v, stable, data, assetNames[0], assetURLs); err != nil {
		return err
	}

	return zip.UnzipToDir(data, targetPath)
}

func (r *TofuRetriever) ListReleases() ([]string, error) {
	r.conf.InitRemoteConf()

	if r.conf.Tofu.GetListMode() == htmlretriever.ListModeHTML {
		baseURL, err := url.JoinPath(r.conf.Tofu.GetListURL(), opentofu, opentofu, github.Releases, github.Download) //nolint
		if err != nil {
			return nil, err
		}

		return htmlretriever.ListReleases(baseURL, r.conf.Tofu.Data, r.conf.DisplayVerbose)
	}

	return github.ListReleases(r.conf.Tofu.GetListURL(), r.conf.GithubToken, r.conf.DisplayVerbose)
}

func (r *TofuRetriever) checkSumAndSig(version *version.Version, stable bool, data []byte, fileName string, assetURLs []string) error {
	dataSums, err := download.Bytes(assetURLs[1], r.conf.DisplayVerbose)
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, fileName); err != nil {
		return err
	}

	dataSumsSig, err := download.Bytes(assetURLs[3], r.conf.DisplayVerbose)
	if err != nil {
		return err
	}

	dataSumsCert, err := download.Bytes(assetURLs[2], r.conf.DisplayVerbose)
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

	if r.conf.DisplayVerbose {
		fmt.Println("cosign executable not found, fallback to pgp check") //nolint
	}

	dataSumsSig, err = download.Bytes(assetURLs[4], r.conf.DisplayVerbose)
	if err != nil {
		return err
	}

	var dataPublicKey []byte
	if r.conf.TofuKeyPath == "" {
		dataPublicKey, err = download.Bytes(publicKeyURL, r.conf.DisplayVerbose)
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
	sumsAssetName := nameBuilder.String() + "SHA256SUMS"

	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOARCH)
	nameBuilder.WriteString(".zip")

	if stable {
		return []string{nameBuilder.String(), sumsAssetName, sumsAssetName + ".pem", sumsAssetName + ".sig", sumsAssetName + ".gpgsig"}
	}

	return []string{nameBuilder.String(), sumsAssetName, sumsAssetName + ".pem", sumsAssetName + ".sig"}
}

func buildIdentity(v *version.Version) string {
	cleanedVersion := v.String()
	indexDot := strings.LastIndexByte(cleanedVersion, '.')
	// cleaned, so indexDot can not be -1
	shortVersion := cleanedVersion[:indexDot]

	return baseIdentity + shortVersion
}
