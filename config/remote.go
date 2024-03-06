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
	"os"
)

const (
	InstallModeDirect = "direct"
	ListModeHTML      = "html"
	ModeAPI           = "api"

	baseGithubURL              = "https://github.com"
	defaultGithubURL           = "https://api.github.com/repos/"
	defaultHashicorpURL        = "https://releases.hashicorp.com"
	defaultTerragruntGithubURL = defaultGithubURL + "gruntwork-io/terragrunt" + slashReleases
	defaultTofuGithubURL       = defaultGithubURL + "opentofu/opentofu" + slashReleases
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

func makeRemoteConfig(remoteURLEnvName string, listURLEnvName string, installModeEnvName string, listModeEnvName string, defaultURL string, defaultBaseURL string) RemoteConfig {
	return RemoteConfig{
		defaultBaseURL: defaultBaseURL, defaultURL: defaultURL, installMode: os.Getenv(installModeEnvName), listMode: os.Getenv(listModeEnvName),
		listURL: os.Getenv(listURLEnvName), RemoteURLEnv: os.Getenv(remoteURLEnvName),
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
	return r.getValueForcedDefault("list_mode", r.listMode, ModeAPI)
}

func (r RemoteConfig) GetListURL() string {
	defaultListURL := r.defaultURL
	if r.GetListMode() == ListModeHTML {
		defaultListURL = r.GetRemoteURL()
	}

	return r.getValueForcedDefault("list_url", r.listURL, defaultListURL)
}

func (r RemoteConfig) GetRemoteURL() string {
	if r.RemoteURL != "" {
		return r.RemoteURL
	}

	return r.getValueForcedDefault("url", r.RemoteURLEnv, r.defaultURL)
}

func (r RemoteConfig) GetRewriteRule() []string {
	oldBase := r.Data["old_base_url"]
	newBase := r.Data["new_base_url"]
	if oldBase != "" && newBase != "" {
		return []string{oldBase, newBase}
	}

	defaultListMode := r.GetListMode() == ModeAPI
	listURL := r.GetListURL()
	remoteURL := r.GetRemoteURL()
	sameURL := remoteURL == listURL
	if defaultListMode && sameURL {
		return nil // no special behaviour, no rewriting
	}

	oneDisabled := defaultListMode || sameURL
	if r.GetInstallMode() == ModeAPI {
		if oldBase == "" {
			oldBase = r.defaultBaseURL
		}

		if newBase == "" {
			if oneDisabled {
				newBase = remoteURL
			} else {
				newBase = listURL
			}
		}

		return []string{oldBase, newBase}
	}

	if oneDisabled {
		return nil // build correct url (direct install mode)
	}

	if oldBase == "" {
		oldBase = remoteURL
	}

	if newBase == "" {
		newBase = listURL
	}

	return []string{oldBase, newBase}
}

func (r RemoteConfig) getValueForcedDefault(name string, forcedValue string, defaultValue string) string {
	if forcedValue != "" {
		return forcedValue
	}

	return MapGetDefault(r.Data, name, defaultValue)
}

func MapGetDefault(m map[string]string, key string, defaultValue string) string {
	if value := m[key]; value != "" {
		return value
	}

	return defaultValue
}
