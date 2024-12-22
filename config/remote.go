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

package config

import (
	"errors"
	"strings"

	configutils "github.com/tofuutils/tenv/v4/config/utils"
	"github.com/tofuutils/tenv/v4/pkg/download"
)

const (
	InstallModeDirect = "direct"
	ListModeHTML      = "html"
	ModeAPI           = "api"

	baseGithubURL              = "https://github.com"
	defaultGithubURL           = "https://api.github.com/repos/"
	defaultHashicorpURL        = "https://releases.hashicorp.com"
	defaultTerragruntGithubURL = defaultGithubURL + "gruntwork-io/terragrunt" + slashReleases
	DefaultTofuGithubURL       = defaultGithubURL + "opentofu/opentofu" + slashReleases
	defaultAtmosGithubURL      = defaultGithubURL + "cloudposse/atmos" + slashReleases
	slashReleases              = "/releases"
)

var (
	ErrInstallMode = errors.New("unknown install mode")
	ErrListMode    = errors.New("unknown list mode")
)

type RemoteConfig struct {
	Data           map[string]string // values from conf file
	defaultBaseURL string
	defaultURL     string
	installMode    string // value from env
	listMode       string // value from env
	listURL        string // value from env
	RemoteURL      string // value from flag
	RemoteURLEnv   string // value from env
}

func makeDefaultRemoteConfig(defaultURL string, defaultBaseURL string) RemoteConfig {
	return RemoteConfig{
		defaultBaseURL: defaultBaseURL, defaultURL: defaultURL, Data: map[string]string{},
	}
}

func makeRemoteConfig(getenv configutils.GetenvFunc, remoteURLEnvName string, listURLEnvName string, installModeEnvName string, listModeEnvName string, defaultURL string, defaultBaseURL string) RemoteConfig {
	return RemoteConfig{
		defaultBaseURL: defaultBaseURL, defaultURL: defaultURL, installMode: getenv(installModeEnvName), listMode: getenv(listModeEnvName),
		listURL: getenv(listURLEnvName), RemoteURLEnv: getenv(remoteURLEnvName),
	}
}

func (r RemoteConfig) GetInstallMode() string {
	defaultInstallMode := ModeAPI
	if r.defaultBaseURL == baseGithubURL && r.GetRemoteURL() != r.defaultURL {
		defaultInstallMode = InstallModeDirect
	}

	return r.getValueForcedDefault("install_mode", r.installMode, defaultInstallMode)
}

func (r RemoteConfig) GetListMode() string {
	defaultListMode := ListModeHTML
	if r.GetListURL() == r.defaultURL {
		defaultListMode = ModeAPI
	}

	return r.getValueForcedDefault("list_mode", r.listMode, defaultListMode)
}

func (r RemoteConfig) GetListURL() string {
	return strings.TrimRight(r.getValueForcedDefault("list_url", r.listURL, r.GetRemoteURL()), "/")
}

func (r RemoteConfig) GetRemoteURL() string {
	remoteURL := r.RemoteURL
	if remoteURL == "" {
		remoteURL = r.getValueForcedDefault("url", r.RemoteURLEnv, r.defaultURL)
	}

	return strings.TrimRight(remoteURL, "/")
}

func (r RemoteConfig) GetRewriteRule() download.URLTransformer {
	oldBase := r.Data["old_base_url"]
	newBase := r.Data["new_base_url"]
	if oldBase != "" && newBase != "" {
		return download.NewURLTransformer(oldBase, newBase)
	}

	if r.GetInstallMode() == InstallModeDirect {
		return download.NoTransform // build correct url
	}

	listURL := r.GetListURL()
	remoteURL := r.GetRemoteURL()
	defaultList := listURL == r.defaultURL
	defaultRemote := remoteURL == r.defaultURL
	if defaultList && defaultRemote {
		return download.NoTransform // no special behaviour, no rewriting
	}

	oldBase = r.defaultBaseURL
	newBase = listURL
	if defaultList {
		newBase = remoteURL
	}

	return download.NewURLTransformer(oldBase, newBase)
}

func (r RemoteConfig) getValueForcedDefault(name string, forcedValue string, defaultValue string) string {
	if forcedValue != "" {
		return forcedValue
	}

	return MapGetDefault(r.Data, name, defaultValue)
}

func MapGetDefault(m map[string]string, key string, defaultValue string) string {
	if value := strings.TrimSpace(m[key]); value != "" {
		return value
	}

	return defaultValue
}

func GetBasicAuthOption(getenv configutils.GetenvFunc, userEnvName string, passEnvName string) []download.RequestOption {
	username, password := getenv(userEnvName), getenv(passEnvName)
	if username == "" || password == "" {
		return nil
	}

	return []download.RequestOption{download.WithBasicAuth(username, password)}
}
