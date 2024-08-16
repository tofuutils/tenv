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

package terraformretriever

import (
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/config/cmdconst"
	"github.com/tofuutils/tenv/v3/pkg/apimsg"
	pgpcheck "github.com/tofuutils/tenv/v3/pkg/check/pgp"
	sha256check "github.com/tofuutils/tenv/v3/pkg/check/sha256"
	"github.com/tofuutils/tenv/v3/pkg/download"
	"github.com/tofuutils/tenv/v3/pkg/pathfilter"
	"github.com/tofuutils/tenv/v3/pkg/winbin"
	"github.com/tofuutils/tenv/v3/pkg/zip"
	htmlretriever "github.com/tofuutils/tenv/v3/versionmanager/retriever/html"
	releaseapi "github.com/tofuutils/tenv/v3/versionmanager/retriever/terraform/api"
)

const (
	publicKeyURL = "https://www.hashicorp.com/.well-known/pgp-key.txt"

	baseFileName = "terraform_"
	indexJson    = "index.json"
)

type TerraformRetriever struct {
	conf *config.Config
}

func Make(conf *config.Config) TerraformRetriever {
	return TerraformRetriever{conf: conf}
}

func (r TerraformRetriever) InstallRelease(version string, targetPath string) error {
	err := r.conf.InitRemoteConf()
	if err != nil {
		return err
	}

	// assume that terraform  version do not start with a 'v'
	if version[0] == 'v' {
		version = version[1:]
	}

	baseVersionURL, err := url.JoinPath(r.conf.Tf.GetRemoteURL(), cmdconst.TerraformName, version) //nolint
	if err != nil {
		return err
	}

	ro := config.GetBasicAuthOption(config.TfRemoteUserEnvName, config.TfRemotePassEnvName)

	var fileName, shaFileName, shaSigFileName, downloadURL, downloadSumsURL, downloadSumsSigURL string
	switch r.conf.Tf.GetInstallMode() {
	case config.InstallModeDirect:
		fileName, shaFileName, shaSigFileName = buildAssetNames(version, r.conf.Arch)
		if r.conf.Displayer.IsDebug() {
			r.conf.Displayer.Log(hclog.Debug, apimsg.MsgSearch, apimsg.AssetsName, []string{fileName, shaFileName, shaSigFileName})
		}

		assetURLs, err := htmlretriever.BuildAssetURLs(baseVersionURL, fileName, shaFileName, shaSigFileName)
		if err != nil {
			return err
		}

		downloadURL, downloadSumsURL, downloadSumsSigURL = assetURLs[0], assetURLs[1], assetURLs[2]
	case config.ModeAPI:
		versionUrl, err := url.JoinPath(baseVersionURL, indexJson) //nolint
		if err != nil {
			return err
		}

		r.conf.Displayer.Display(apimsg.MsgFetchRelease + versionUrl)

		value, err := download.JSON(versionUrl, download.NoDisplay, ro...)
		if err != nil {
			return err
		}

		fileName, downloadURL, shaFileName, shaSigFileName, err = releaseapi.ExtractAssetUrls(runtime.GOOS, r.conf.Arch, value)
		if err != nil {
			return err
		}

		if r.conf.Displayer.IsDebug() {
			r.conf.Displayer.Log(hclog.Debug, apimsg.MsgSearch, apimsg.AssetsName, []string{fileName, shaFileName, shaSigFileName})
		}

		assetURLs, err := htmlretriever.BuildAssetURLs(baseVersionURL, shaFileName, shaSigFileName)
		if err != nil {
			return err
		}

		downloadSumsURL, downloadSumsSigURL = assetURLs[0], assetURLs[1]
	default:
		return config.ErrInstallMode
	}

	urlTranformer := download.UrlTranformer(r.conf.Tf.GetRewriteRule())
	assetURLs, err := download.ApplyUrlTranformer(urlTranformer, downloadURL, downloadSumsURL, downloadSumsSigURL)
	if err != nil {
		return err
	}

	data, err := download.Bytes(assetURLs[0], r.conf.Displayer.Display, ro...)
	if err != nil {
		return err
	}

	if err = r.checkSumAndSig(fileName, data, assetURLs[1], assetURLs[2], ro); err != nil {
		return err
	}

	return zip.UnzipToDir(data, targetPath, pathfilter.NameEqual(winbin.GetBinaryName(cmdconst.TerraformName)))
}

func (r TerraformRetriever) ListReleases() ([]string, error) {
	err := r.conf.InitRemoteConf()
	if err != nil {
		return nil, err
	}

	baseURL, err := url.JoinPath(r.conf.Tf.GetListURL(), cmdconst.TerraformName) //nolint
	if err != nil {
		return nil, err
	}

	ro := config.GetBasicAuthOption(config.TfRemoteUserEnvName, config.TfRemotePassEnvName)

	switch r.conf.Tf.GetListMode() {
	case config.ListModeHTML:
		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + baseURL)

		return htmlretriever.ListReleases(baseURL, r.conf.Tf.Data, ro)
	case config.ModeAPI:
		releasesURL, err := url.JoinPath(baseURL, indexJson) //nolint
		if err != nil {
			return nil, err
		}

		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + releasesURL)

		value, err := download.JSON(releasesURL, download.NoDisplay, ro...)
		if err != nil {
			return nil, err
		}

		return releaseapi.ExtractReleases(value)
	default:
		return nil, config.ErrListMode
	}
}

func (r TerraformRetriever) checkSumAndSig(fileName string, data []byte, downloadSumsURL string, downloadSumsSigURL string, ro []download.RequestOption) error {
	dataSums, err := download.Bytes(downloadSumsURL, r.conf.Displayer.Display, ro...)
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, fileName); err != nil {
		return err
	}

	if r.conf.SkipSignature {
		return nil
	}

	dataSumsSig, err := download.Bytes(downloadSumsSigURL, r.conf.Displayer.Display, ro...)
	if err != nil {
		return err
	}

	var dataPublicKey []byte
	if r.conf.TfKeyPath == "" {
		dataPublicKey, err = download.Bytes(publicKeyURL, r.conf.Displayer.Display)
	} else {
		dataPublicKey, err = os.ReadFile(r.conf.TfKeyPath)
	}

	if err != nil {
		return err
	}

	return pgpcheck.Check(dataSums, dataSumsSig, dataPublicKey)
}

func buildAssetNames(version string, arch string) (string, string, string) {
	var nameBuilder strings.Builder
	nameBuilder.WriteString(baseFileName)
	nameBuilder.WriteString(version)
	nameBuilder.WriteByte('_')
	sumsAssetName := nameBuilder.String() + "SHA256SUMS"

	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(arch)
	nameBuilder.WriteString(".zip")

	return nameBuilder.String(), sumsAssetName, sumsAssetName + ".sig"
}
