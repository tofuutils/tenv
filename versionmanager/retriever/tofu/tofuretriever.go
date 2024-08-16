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
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/config/cmdconst"
	"github.com/tofuutils/tenv/v3/pkg/apimsg"
	cosigncheck "github.com/tofuutils/tenv/v3/pkg/check/cosign"
	pgpcheck "github.com/tofuutils/tenv/v3/pkg/check/pgp"
	sha256check "github.com/tofuutils/tenv/v3/pkg/check/sha256"
	"github.com/tofuutils/tenv/v3/pkg/download"
	"github.com/tofuutils/tenv/v3/pkg/github"
	"github.com/tofuutils/tenv/v3/pkg/pathfilter"
	"github.com/tofuutils/tenv/v3/pkg/winbin"
	"github.com/tofuutils/tenv/v3/pkg/zip"
	htmlretriever "github.com/tofuutils/tenv/v3/versionmanager/retriever/html"
	tofudlmirroring "github.com/tofuutils/tenv/v3/versionmanager/retriever/tofu/dl"
)

const (
	modeMirroring = "mirror"

	getTofuURL              = "https://get.opentofu.org/"
	defaultTofuMirroringURL = getTofuURL + "tofu/api.json"
	publicKeyURL            = getTofuURL + "opentofu.asc"

	defaultTofuURLTemplate = "https://github.com/opentofu/opentofu/releases/download/v{{ .Version }}/{{ .Artifact }}"

	baseIdentity     = "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v"
	issuer           = "https://token.actions.githubusercontent.com"
	unstableIdentity = "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/main"

	baseFileName = "tofu_"
	opentofu     = "opentofu"
)

type TofuRetriever struct {
	conf *config.Config
}

func Make(conf *config.Config) TofuRetriever {
	return TofuRetriever{conf: conf}
}

func (r TofuRetriever) InstallRelease(versionStr string, targetPath string) error {
	err := r.conf.InitRemoteConf()
	if err != nil {
		return err
	}

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
	assetNames := buildAssetNames(versionStr, r.conf.Arch, stable)
	if r.conf.Displayer.IsDebug() {
		r.conf.Displayer.Log(hclog.Debug, apimsg.MsgSearch, apimsg.AssetsName, assetNames)
	}

	switch r.conf.Tofu.GetInstallMode() {
	case config.InstallModeDirect:
		baseAssetURL, err2 := url.JoinPath(r.conf.Tofu.GetRemoteURL(), opentofu, opentofu, github.Releases, github.Download, tag) //nolint
		if err2 != nil {
			return err2
		}

		assetURLs, err = htmlretriever.BuildAssetURLs(baseAssetURL, assetNames...)
	case config.ModeAPI:
		assetURLs, err = github.AssetDownloadURL(tag, assetNames, r.conf.Tofu.GetRemoteURL(), r.conf.GithubToken, r.conf.Displayer.Display)
	case modeMirroring:
		urlTemplate := os.Getenv(config.TofuURLTemplateEnvName)
		if urlTemplate == "" {
			urlTemplate = defaultTofuURLTemplate
		}

		builder, err2 := tofudlmirroring.MakeURLBuilder(urlTemplate, versionStr)
		if err2 != nil {
			return err2
		}

		assetURLs, err = download.ApplyUrlTranformer(builder.Build, assetNames...)
	default:
		return config.ErrInstallMode
	}
	if err != nil {
		return err
	}

	urlTranformer := download.UrlTranformer(r.conf.Tofu.GetRewriteRule())
	assetURLs, err = download.ApplyUrlTranformer(urlTranformer, assetURLs...)
	if err != nil {
		return err
	}

	ro := config.GetBasicAuthOption(config.TofuRemoteUserEnvName, config.TofuRemotePassEnvName)
	data, err := download.Bytes(assetURLs[0], r.conf.Displayer.Display, ro...)
	if err != nil {
		return err
	}

	if err = r.checkSumAndSig(v, stable, data, assetNames[0], assetURLs, ro); err != nil {
		return err
	}

	return zip.UnzipToDir(data, targetPath, pathfilter.NameEqual(winbin.GetBinaryName(cmdconst.TofuName)))
}

func (r TofuRetriever) ListReleases() ([]string, error) {
	err := r.conf.InitRemoteConf()
	if err != nil {
		return nil, err
	}

	ro := config.GetBasicAuthOption(config.TofuRemoteUserEnvName, config.TofuRemotePassEnvName)

	listURL := r.conf.Tofu.GetListURL()
	switch r.conf.Tofu.GetListMode() {
	case config.ListModeHTML:
		baseURL, err := url.JoinPath(listURL, opentofu, opentofu, github.Releases, github.Download) //nolint
		if err != nil {
			return nil, err
		}

		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + baseURL)

		return htmlretriever.ListReleases(baseURL, r.conf.Tofu.Data, ro)
	case config.ModeAPI:
		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + listURL)

		return github.ListReleases(listURL, r.conf.GithubToken)
	case modeMirroring:
		if listURL == config.DefaultTofuGithubURL {
			listURL = defaultTofuMirroringURL
		}

		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + listURL)

		value, err := download.JSON(listURL, download.NoDisplay, ro...)
		if err != nil {
			return nil, err
		}

		return tofudlmirroring.ExtractReleases(value)
	default:
		return nil, config.ErrListMode
	}
}

func (r TofuRetriever) checkSumAndSig(version *version.Version, stable bool, data []byte, fileName string, assetURLs []string, ro []download.RequestOption) error {
	dataSums, err := download.Bytes(assetURLs[1], r.conf.Displayer.Display, ro...)
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, fileName); err != nil {
		return err
	}

	if r.conf.SkipSignature {
		return nil
	}

	dataSumsSig, err := download.Bytes(assetURLs[3], r.conf.Displayer.Display, ro...)
	if err != nil {
		return err
	}

	dataSumsCert, err := download.Bytes(assetURLs[2], r.conf.Displayer.Display, ro...)
	if err != nil {
		return err
	}

	identity := buildIdentity(version, stable)
	err = cosigncheck.Check(dataSums, dataSumsSig, dataSumsCert, identity, issuer, r.conf.Displayer)
	if err == nil || err != cosigncheck.ErrNotInstalled {
		return err
	}

	if !stable {
		r.conf.Displayer.Display("skip signature check : cosign executable not found and pgp check not available for unstable version")

		return nil
	}

	r.conf.Displayer.Display("cosign executable not found, fallback to pgp check")

	dataSumsSig, err = download.Bytes(assetURLs[4], r.conf.Displayer.Display, ro...)
	if err != nil {
		return err
	}

	var dataPublicKey []byte
	if r.conf.TofuKeyPath == "" {
		dataPublicKey, err = download.Bytes(publicKeyURL, r.conf.Displayer.Display)
	} else {
		dataPublicKey, err = os.ReadFile(r.conf.TofuKeyPath)
	}

	if err != nil {
		return err
	}

	return pgpcheck.Check(dataSums, dataSumsSig, dataPublicKey)
}

func buildAssetNames(version string, arch string, stable bool) []string {
	var nameBuilder strings.Builder
	nameBuilder.WriteString(baseFileName)
	nameBuilder.WriteString(version)
	nameBuilder.WriteByte('_')
	sumsAssetName := nameBuilder.String() + "SHA256SUMS"

	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(arch)
	nameBuilder.WriteString(".zip")

	if stable {
		return []string{nameBuilder.String(), sumsAssetName, sumsAssetName + ".pem", sumsAssetName + ".sig", sumsAssetName + ".gpgsig"}
	}

	return []string{nameBuilder.String(), sumsAssetName, sumsAssetName + ".pem", sumsAssetName + ".sig"}
}

func buildIdentity(v *version.Version, stable bool) string {
	if !stable {
		return unstableIdentity
	}

	cleanedVersion := v.String()
	indexDot := strings.LastIndexByte(cleanedVersion, '.')
	// cleaned, so indexDot can not be -1
	shortVersion := cleanedVersion[:indexDot]

	return baseIdentity + shortVersion
}
