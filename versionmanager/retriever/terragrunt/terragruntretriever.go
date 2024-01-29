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

const baseFileName = "terragrunt_"

type TerragruntRetriever struct {
	conf *config.Config
}

func NewTerragruntRetriever(conf *config.Config) *TerragruntRetriever {
	return &TerragruntRetriever{conf: conf}
}

func (r *TerragruntRetriever) InstallRelease(versionStr string, targetPath string) error {
	tag := versionStr
	// assume that terragrunt tags start with a 'v'
	if tag[0] != 'v' {
		tag = "v" + versionStr
	}

	assetNames := buildAssetNames()
	assets, err := github.DownloadAssetURL(tag, assetNames, r.conf.TgRemoteURL, r.conf.GithubToken)
	if err != nil {
		return err
	}

	data, err := download.Bytes(assets[assetNames[0]])
	if err != nil {
		return err
	}

	dataSums, err := download.Bytes(assets[assetNames[1]])
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
	return github.LatestRelease(r.conf.TgRemoteURL, r.conf.GithubToken)
}

func (r *TerragruntRetriever) ListReleases() ([]string, error) {
	return github.ListReleases(r.conf.TgRemoteURL, r.conf.GithubToken)
}

func buildAssetNames() []string {
	var nameBuilder strings.Builder
	nameBuilder.WriteString(baseFileName)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOARCH)
	if runtime.GOOS == "windows" {
		nameBuilder.WriteString(".exe")
	}

	return []string{nameBuilder.String(), "SHA256SUMS"}
}
