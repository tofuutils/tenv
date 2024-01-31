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
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tofuutils/tenv/config"
	sha256check "github.com/tofuutils/tenv/pkg/check/sha256"
	"github.com/tofuutils/tenv/pkg/download"
	"github.com/tofuutils/tenv/pkg/github"
)

const (
	defaultTerragruntGithubURL = "https://api.github.com/repos/gruntwork-io/terragrunt/releases"

	baseFileName = "terragrunt_"
)

type TerragruntRetriever struct {
	conf       *config.Config
	notLoaded  bool
	remoteConf map[string]string
}

func NewTerragruntRetriever(conf *config.Config) *TerragruntRetriever {
	return &TerragruntRetriever{conf: conf, notLoaded: true}
}

func (r *TerragruntRetriever) InstallRelease(versionStr string, targetPath string) error {
	tag := versionStr
	// assume that terragrunt tags start with a 'v'
	if tag[0] != 'v' {
		tag = "v" + versionStr
	}

	assetNames := buildAssetNames()
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

	downloadURL, err = urlTranformer(assets[assetNames[1]])
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

	err = os.MkdirAll(targetPath, 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(targetPath, "terragrunt"), data, 0755)
}

func (r *TerragruntRetriever) LatestRelease() (string, error) {
	return github.LatestRelease(r.getRemoteURL(), r.conf.GithubToken)
}

func (r *TerragruntRetriever) ListReleases() ([]string, error) {
	return github.ListReleases(r.getRemoteURL(), r.conf.GithubToken)
}

func (r *TerragruntRetriever) getRemoteURL() string {
	if r.conf.TgRemoteURL != "" {
		return r.conf.TgRemoteURL
	}

	return config.MapGetDefault(r.readRemoteConf(), "url", defaultTerragruntGithubURL)
}

func (r *TerragruntRetriever) readRemoteConf() map[string]string {
	if r.notLoaded {
		r.notLoaded = false
		r.remoteConf = r.conf.ReadRemoteConf("terragrunt")
	}

	return r.remoteConf
}

func buildAssetNames() []string {
	var nameBuilder strings.Builder
	nameBuilder.WriteString(baseFileName)
	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOARCH)
	if runtime.GOOS == "windows" {
		nameBuilder.WriteString(".exe")
	}

	return []string{nameBuilder.String(), "SHA256SUMS"}
}
