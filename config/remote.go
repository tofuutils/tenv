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

import "os"

const (
	baseGithubURL              = "https://github.com"
	defaultGithubURL           = "https://api.github.com/repos/"
	defaultHashicorpURL        = "https://releases.hashicorp.com"
	defaultTerragruntGithubURL = defaultGithubURL + "gruntwork-io/terragrunt" + slashReleases
	defaultTofuGithubURL       = defaultGithubURL + "opentofu/opentofu" + slashReleases
	slashReleases              = "/releases"
)

type RemoteConfig struct {
	Data           map[string]string
	defaultBaseURL string
	defaultURL     string
	installMode    string
	listMode       string
	listURL        string
	RemoteURL      string
}

func makeRemoteConfig(remoteURLEnvName string, listURLEnvName string, installModeEnvName string, listModeEnvName string, defaultURL string, defaultBaseURL string) RemoteConfig {
	return RemoteConfig{
		defaultBaseURL: defaultBaseURL, defaultURL: defaultURL, installMode: os.Getenv(installModeEnvName), listMode: os.Getenv(listModeEnvName),
		listURL: os.Getenv(listURLEnvName), RemoteURL: os.Getenv(remoteURLEnvName),
	}
}

func (r RemoteConfig) GetInstallMode() string {
	return r.getValueForced("install_mode", r.installMode)
}

func (r RemoteConfig) GetListMode() string {
	return r.getValueForced("list_mode", r.listMode)
}

func (r RemoteConfig) GetListURL() string {
	return r.getValueForcedDefault("list_url", r.listURL, r.GetRemoteURL())
}

func (r RemoteConfig) GetRemoteURL() string {
	return r.getValueForcedDefault("url", r.RemoteURL, r.defaultURL)
}

func (r RemoteConfig) GetRewriteRule() []string {
	oldBase := r.Data["old_base_url"]
	newBase := r.Data["new_base_url"]
	if oldBase != "" && newBase != "" {
		return []string{oldBase, newBase}
	}

	emptyListMode := r.GetListMode() == ""
	listURL := r.GetListURL()
	remoteURL := r.GetRemoteURL()
	sameURL := remoteURL == listURL
	if emptyListMode && sameURL {
		return nil // no special behaviour, no rewriting
	}

	oneDisabled := emptyListMode || sameURL
	if r.GetInstallMode() == "" {
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
		return nil // build correct url (install mode activated)
	}

	if oldBase == "" {
		oldBase = remoteURL
	}

	if newBase == "" {
		newBase = listURL
	}

	return []string{oldBase, newBase}
}

func MapGetDefault(m map[string]string, key string, defaultValue string) string {
	if value := m[key]; value != "" {
		return value
	}

	return defaultValue
}

func (r RemoteConfig) getValueForced(name string, forcedValue string) string {
	if forcedValue != "" {
		return forcedValue
	}

	return r.Data[name]
}

func (r RemoteConfig) getValueForcedDefault(name string, forcedValue string, defaultValue string) string {
	if forcedValue != "" {
		return forcedValue
	}

	return MapGetDefault(r.Data, name, defaultValue)
}
