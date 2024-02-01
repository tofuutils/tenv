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
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tofuutils/tenv/config"
	sha256check "github.com/tofuutils/tenv/pkg/check/sha256"
	"github.com/tofuutils/tenv/pkg/download"
	"github.com/tofuutils/tenv/pkg/github"
	htmlretriever "github.com/tofuutils/tenv/versionmanager/retriever/html"
	"github.com/tofuutils/tenv/versionmanager/semantic"
)

const (
	defaultTerragruntGithubURL = "https://api.github.com/repos/gruntwork-io/terragrunt/releases"

	baseFileName  = "terragrunt_"
	gruntworkName = "gruntwork-io"
	Name          = "terragrunt"
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

	var err error
	var assetURLs []string
	assetNames := buildAssetNames()
	if r.readRemoteConf()[htmlretriever.InstallMode] == htmlretriever.InstallModeDirect {
		baseAssetURL, err2 := url.JoinPath(r.getRemoteURL(), gruntworkName, Name, github.Releases, github.Download, tag)
		if err2 != nil {
			return err2
		}

		assetURLs, err = htmlretriever.BuildAssetURLs(baseAssetURL, assetNames)
	} else {
		assetURLs, err = github.AssetDownloadURL(tag, assetNames, r.getRemoteURL(), r.conf.GithubToken, r.conf.Verbose)
	}
	if err != nil {
		return err
	}

	urlTranformer := download.UrlTranformer(r.readRemoteConf())
	downloadURL, err := urlTranformer(assetURLs[0])
	if err != nil {
		return err
	}

	data, err := download.Bytes(downloadURL, r.conf.Verbose)
	if err != nil {
		return err
	}

	downloadURL, err = urlTranformer(assetURLs[1])
	if err != nil {
		return err
	}

	dataSums, err := download.Bytes(downloadURL, r.conf.Verbose)
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

	return os.WriteFile(filepath.Join(targetPath, Name), data, 0755)
}

func (r *TerragruntRetriever) LatestRelease() (string, error) {
	if r.readRemoteConf()[htmlretriever.ListMode] == htmlretriever.ListModeHTML {
		versions, err := r.ListReleases()
		if err != nil {
			return "", err
		}

		return semantic.LatestVersionFromList(versions)
	}

	return github.LatestRelease(r.getRemoteURL(), r.conf.GithubToken, r.conf.Verbose)
}

func (r *TerragruntRetriever) ListReleases() ([]string, error) {
	remoteConf := r.readRemoteConf()
	listRemoteURL := config.MapGetDefault(remoteConf, htmlretriever.ListURL, r.getRemoteURL())

	if remoteConf[htmlretriever.ListMode] == htmlretriever.ListModeHTML {
		baseURL, err := url.JoinPath(listRemoteURL, gruntworkName, Name, github.Releases, github.Download)
		if err != nil {
			return nil, err
		}

		return htmlretriever.ListReleases(baseURL, remoteConf, r.conf.Verbose)
	}

	return github.ListReleases(r.getRemoteURL(), r.conf.GithubToken, r.conf.Verbose)
}

func (r *TerragruntRetriever) getRemoteURL() string {
	if r.conf.TgRemoteURL != "" {
		return r.conf.TgRemoteURL
	}

	return config.MapGetDefault(r.readRemoteConf(), htmlretriever.URL, defaultTerragruntGithubURL)
}

func (r *TerragruntRetriever) readRemoteConf() map[string]string {
	if r.notLoaded {
		r.notLoaded = false
		r.remoteConf = r.conf.ReadRemoteConf(Name)
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
