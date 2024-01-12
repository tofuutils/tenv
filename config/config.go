/*
 *
 * Copyright 2024 gotofuenv authors.
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
	"os"
	"path"
	"strconv"
)

const (
	LatestAllowedKey = "latest-allowed"
	LatestStableKey  = "latest-stable"
	LatestKey        = "latest"
	MinRequiredKey   = "min-required"

	TfFolderName      = "Terraform"
	TfVersionFileName = ".terraform-version"

	TofuFolderName      = "OpenTofu"
	TofuVersionFileName = ".opentofu-version"
)

const (
	defaultTfHashicorpUrl = "https://releases.hashicorp.com/terraform/index.json"
	defaultTofuGithubUrl  = "https://api.github.com/repos/opentofu/opentofu/releases"

	autoInstallEnvName = "AUTO_INSTALL"
	remoteUrlEnvName   = "REMOTE"
	rootPathEnvName    = "ROOT"
	verboseEnvName     = "VERBOSE"

	tfenvPrefix          = "TFENV_"
	tfAutoInstallEnvName = tfenvPrefix + autoInstallEnvName
	tfRemoteUrlEnvName   = tfenvPrefix + remoteUrlEnvName
	tfRootPathEnvName    = tfenvPrefix + rootPathEnvName
	tfVerboseEnvName     = tfenvPrefix + verboseEnvName
	tfVersionEnvName     = tfenvPrefix + "TERRAFORM_VERSION"

	tofuenvPrefix          = "TOFUENV_"
	tofuAutoInstallEnvName = tofuenvPrefix + autoInstallEnvName
	tofuRemoteUrlEnvName   = tofuenvPrefix + remoteUrlEnvName
	tofuRootPathEnvName    = tofuenvPrefix + rootPathEnvName
	tofuTokenEnvName       = tofuenvPrefix + "GITHUB_TOKEN"
	tofuVerboseEnvName     = tofuenvPrefix + verboseEnvName
	tofuVersionEnvName     = tofuenvPrefix + "TOFU_VERSION"
)

type Config struct {
	NoInstall     bool
	TfRemoteUrl   string
	TofuRemoteUrl string
	RootPath      string
	TfVersion     string
	TofuVersion   string
	GithubToken   string
	UserPath      string
	Verbose       bool
}

func InitConfigFromEnv() (Config, error) {
	userPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	autoInstall := true
	autoInstallStr := getenv(tofuAutoInstallEnvName, tfAutoInstallEnvName)
	if autoInstallStr != "" {
		var err error
		autoInstall, err = strconv.ParseBool(autoInstallStr)
		if err != nil {
			return Config{}, err
		}
	}

	tfRemoteUrl := os.Getenv(tfRemoteUrlEnvName)
	if tfRemoteUrl == "" {
		tfRemoteUrl = defaultTfHashicorpUrl
	}

	tofuRemoteUrl := os.Getenv(tofuRemoteUrlEnvName)
	if tofuRemoteUrl == "" {
		tofuRemoteUrl = defaultTofuGithubUrl
	}

	rootPath := getenv(tofuRootPathEnvName, tfRootPathEnvName)
	if rootPath == "" {
		rootPath = path.Join(userPath, ".gotofuenv")
	}

	verbose := false
	verboseStr := getenv(tofuVerboseEnvName, tfVerboseEnvName)
	if verboseStr != "" {
		verbose, err = strconv.ParseBool(verboseStr)
		if err != nil {
			return Config{}, err
		}
	}

	return Config{
		NoInstall:     !autoInstall,
		TfRemoteUrl:   tfRemoteUrl,
		TofuRemoteUrl: tofuRemoteUrl,
		RootPath:      rootPath,
		TfVersion:     os.Getenv(tfVersionEnvName),
		TofuVersion:   os.Getenv(tofuVersionEnvName),
		GithubToken:   os.Getenv(tofuTokenEnvName),
		UserPath:      userPath,
		Verbose:       verbose,
	}, nil
}

func getenv(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return ""
}
