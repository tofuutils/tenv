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

package terragruntretriever

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/config/cmdconst"
	"github.com/tofuutils/tenv/v3/pkg/apimsg"
	sha256check "github.com/tofuutils/tenv/v3/pkg/check/sha256"
	"github.com/tofuutils/tenv/v3/pkg/download"
	"github.com/tofuutils/tenv/v3/pkg/github"
	"github.com/tofuutils/tenv/v3/pkg/winbin"
	htmlretriever "github.com/tofuutils/tenv/v3/versionmanager/retriever/html"
)

const (
	baseFileName  = "terragrunt_"
	gruntworkName = "gruntwork-io"

	rwePerm = 0o755
)

type TerragruntRetriever struct {
	conf *config.Config
}

func Make(conf *config.Config) TerragruntRetriever {
	return TerragruntRetriever{conf: conf}
}

func (r TerragruntRetriever) Install(ctx context.Context, versionStr string, targetPath string) error {
	err := r.conf.InitRemoteConf()
	if err != nil {
		return err
	}

	tag := versionStr
	// assume that terragrunt tags start with a 'v'
	if tag[0] != 'v' {
		tag = "v" + versionStr
	}

	var assetURLs []string
	fileName, shaFileName := buildAssetNames(r.conf.Arch)
	if r.conf.Displayer.IsDebug() {
		r.conf.Displayer.Log(hclog.Debug, apimsg.MsgSearch, apimsg.AssetsName, []string{fileName, shaFileName})
	}

	switch r.conf.Tg.GetInstallMode() {
	case config.InstallModeDirect:
		baseAssetURL, err2 := url.JoinPath(r.conf.Tg.GetRemoteURL(), gruntworkName, cmdconst.TerragruntName, github.Releases, github.Download, tag)
		if err2 != nil {
			return err2
		}

		assetURLs, err = htmlretriever.BuildAssetURLs(baseAssetURL, fileName, shaFileName)
	case config.ModeAPI:
		assetURLs, err = github.AssetDownloadURL(ctx, tag, []string{fileName, shaFileName}, r.conf.Tg.GetRemoteURL(), r.conf.GithubToken, r.conf.Displayer.Display)
	default:
		return config.ErrInstallMode
	}
	if err != nil {
		return err
	}

	assetURLs, err = download.ApplyURLTransformer(r.conf.Tg.GetRewriteRule(), assetURLs...)
	if err != nil {
		return err
	}

	requestOptions := config.GetBasicAuthOption(r.conf.Getenv, config.TgRemoteUserEnvName, config.TgRemotePassEnvName)
	data, err := download.Bytes(ctx, assetURLs[0], r.conf.Displayer.Display, download.NoCheck, requestOptions...)
	if err != nil {
		return err
	}

	dataSums, err := download.Bytes(ctx, assetURLs[1], r.conf.Displayer.Display, download.NoCheck, requestOptions...)
	if err != nil {
		return err
	}

	if err = sha256check.Check(data, dataSums, fileName); err != nil {
		return err
	}

	err = os.MkdirAll(targetPath, rwePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(targetPath, winbin.GetBinaryName(cmdconst.TerragruntName)), data, rwePerm)
}

func (r TerragruntRetriever) ListVersions(ctx context.Context) ([]string, error) {
	err := r.conf.InitRemoteConf()
	if err != nil {
		return nil, err
	}

	requestOptions := config.GetBasicAuthOption(r.conf.Getenv, config.TgRemoteUserEnvName, config.TgRemotePassEnvName)

	listURL := r.conf.Tg.GetListURL()
	switch r.conf.Tg.GetListMode() {
	case config.ListModeHTML:
		baseURL, err := url.JoinPath(listURL, gruntworkName, cmdconst.TerragruntName, github.Releases, github.Download)
		if err != nil {
			return nil, err
		}

		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + baseURL)

		return htmlretriever.ListReleases(ctx, baseURL, r.conf.Tg.Data, requestOptions)
	case config.ModeAPI:
		r.conf.Displayer.Display(apimsg.MsgFetchAllReleases + listURL)

		return github.ListReleases(ctx, listURL, r.conf.GithubToken)
	default:
		return nil, config.ErrListMode
	}
}

func buildAssetNames(arch string) (string, string) {
	var nameBuilder strings.Builder
	nameBuilder.WriteString(baseFileName)
	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(arch)
	_, _ = winbin.WriteSuffixTo(&nameBuilder)

	return nameBuilder.String(), "SHA256SUMS"
}
