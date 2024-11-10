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
	"context"
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
	indexJSON    = "index.json"
)

type TerraformRetriever struct {
	conf *config.Config
}

func Make(conf *config.Config) TerraformRetriever {
	return TerraformRetriever{conf: conf}
}

func (r TerraformRetriever) Install(ctx context.Context, version string, targetPath string) error {
	err := r.conf.InitRemoteConf()
	if err != nil {
		return err
	}

	// assume that terraform  version do not start with a 'v'
	if version[0] == 'v' {
		version = version[1:]
	}

	baseVersionURL, err := url.JoinPath(r.conf.Tf.GetRemoteURL(), cmdconst.TerraformName, version)
	if err != nil {
		return err
	}

	requestOptions := config.GetBasicAuthOption(r.conf.Getenv, config.TfRemoteUserEnvName, config.TfRemotePassEnvName)

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
		versionURL, err := url.JoinPath(baseVersionURL, indexJSON)
		if err != nil {
			return err
		}

		r.conf.Displayer.Display(apimsg.MsgFetchRelease + versionURL)

		value, err := download.JSON(ctx, versionURL, download.NoDisplay, requestOptions...)
		if err != nil {
			return err
		}

		fileName, downloadURL, shaFileName, shaSigFileName, err = releaseapi.ExtractAssetURLs(runtime.GOOS, r.conf.Arch, value)
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

	assetURLs, err := download.ApplyURLTransformer(r.conf.Tf.GetRewriteRule(), downloadURL, downloadSumsURL, downloadSumsSigURL)
	if err != nil {
		return err
	}

	data, err := download.Bytes(ctx, assetURLs[0], r.conf.Displayer.Display, requestOptions...)
	if err != nil {
		return err
	}

	if err = r.checkSumAndSig(ctx, fileName, data, assetURLs[1], assetURLs[2], requestOptions); err != nil {
		return err
	}

	return zip.UnzipToDir(data, targetPath, pathfilter.NameEqual(winbin.GetBinaryName(cmdconst.TerraformName)))
}

func (r TerraformRetriever) ListVersions(ctx context.Context) ([]string, error) {
	err := r.conf.InitRemoteConf()
	if err != nil {
		return nil, err
	}

	baseURL, err := url.JoinPath(r.conf.Tf.GetListURL(), cmdconst.TerraformName)
	if err != nil {
		return nil, err
	}

	requestOptions := config.GetBasicAuthOption(r.conf.Getenv, config.TfRemoteUserEnvName, config.TfRemotePassEnvName)

	switch r.conf.Tf.GetListMode() {
	case config.ListModeHTML:
		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + baseURL)

		return htmlretriever.ListReleases(ctx, baseURL, r.conf.Tf.Data, requestOptions)
	case config.ModeAPI:
		releasesURL, err := url.JoinPath(baseURL, indexJSON)
		if err != nil {
			return nil, err
		}

		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + releasesURL)

		value, err := download.JSON(ctx, releasesURL, download.NoDisplay, requestOptions...)
		if err != nil {
			return nil, err
		}

		return releaseapi.ExtractReleases(value)
	default:
		return nil, config.ErrListMode
	}
}

func (r TerraformRetriever) checkSumAndSig(ctx context.Context, fileName string, data []byte, downloadSumsURL string, downloadSumsSigURL string, options []download.RequestOption) error {
	dataSums, err := download.Bytes(ctx, downloadSumsURL, r.conf.Displayer.Display, options...)
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, fileName); err != nil {
		return err
	}

	if r.conf.SkipSignature {
		return nil
	}

	dataSumsSig, err := download.Bytes(ctx, downloadSumsSigURL, r.conf.Displayer.Display, options...)
	if err != nil {
		return err
	}

	var dataPublicKey []byte
	if r.conf.TfKeyPath == "" {
		dataPublicKey, err = download.Bytes(ctx, publicKeyURL, r.conf.Displayer.Display)
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
