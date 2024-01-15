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
	defaultTfHashicorpUrl = "https://releases.hashicorp.com/terraform"
	defaultTofuGithubUrl  = "https://api.github.com/repos/opentofu/opentofu/releases"

	autoInstallEnvName = "AUTO_INSTALL"
	forceRemoteEnvName = "FORCE_REMOTE"
	remoteUrlEnvName   = "REMOTE"
	rootPathEnvName    = "ROOT"
	verboseEnvName     = "VERBOSE"

	tfenvPrefix          = "TFENV_"
	tfAutoInstallEnvName = tfenvPrefix + autoInstallEnvName
	tfForceRemoteEnvName = tfenvPrefix + forceRemoteEnvName
	TfRemoteUrlEnvName   = tfenvPrefix + remoteUrlEnvName
	tfRootPathEnvName    = tfenvPrefix + rootPathEnvName
	tfVerboseEnvName     = tfenvPrefix + verboseEnvName
	TfVersionEnvName     = tfenvPrefix + "TERRAFORM_VERSION"

	tofuenvPrefix          = "TOFUENV_"
	tofuAutoInstallEnvName = tofuenvPrefix + autoInstallEnvName
	tofuForceRemoteEnvName = tofuenvPrefix + forceRemoteEnvName
	TofuRemoteUrlEnvName   = tofuenvPrefix + remoteUrlEnvName
	tofuRootPathEnvName    = tofuenvPrefix + rootPathEnvName
	tofuTokenEnvName       = tofuenvPrefix + "GITHUB_TOKEN"
	tofuVerboseEnvName     = tofuenvPrefix + verboseEnvName
	TofuVersionEnvName     = tofuenvPrefix + "TOFU_VERSION"
)

type Config struct {
	ForceRemote   bool
	GithubToken   string
	NoInstall     bool
	RootPath      string
	TfRemoteUrl   string
	TofuRemoteUrl string
	UserPath      string
	Verbose       bool
}

func InitConfigFromEnv() (Config, error) {
	userPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	autoInstall := true
	autoInstallStr := getenvFallback(tofuAutoInstallEnvName, tfAutoInstallEnvName)
	if autoInstallStr != "" {
		autoInstall, err = strconv.ParseBool(autoInstallStr)
		if err != nil {
			return Config{}, err
		}
	}

	forceRemote := false
	forceRemoteStr := getenvFallback(tofuForceRemoteEnvName, tfForceRemoteEnvName)
	if forceRemoteStr != "" {
		forceRemote, err = strconv.ParseBool(forceRemoteStr)
		if err != nil {
			return Config{}, err
		}
	}

	tfRemoteUrl := os.Getenv(TfRemoteUrlEnvName)
	if tfRemoteUrl == "" {
		tfRemoteUrl = defaultTfHashicorpUrl
	}

	tofuRemoteUrl := os.Getenv(TofuRemoteUrlEnvName)
	if tofuRemoteUrl == "" {
		tofuRemoteUrl = defaultTofuGithubUrl
	}

	rootPath := getenvFallback(tofuRootPathEnvName, tfRootPathEnvName)
	if rootPath == "" {
		rootPath = path.Join(userPath, ".gotofuenv")
	}

	verbose := false
	verboseStr := getenvFallback(tofuVerboseEnvName, tfVerboseEnvName)
	if verboseStr != "" {
		verbose, err = strconv.ParseBool(verboseStr)
		if err != nil {
			return Config{}, err
		}
	}

	return Config{
		ForceRemote:   forceRemote,
		GithubToken:   os.Getenv(tofuTokenEnvName),
		NoInstall:     !autoInstall,
		RootPath:      rootPath,
		TfRemoteUrl:   tfRemoteUrl,
		TofuRemoteUrl: tofuRemoteUrl,
		UserPath:      userPath,
		Verbose:       verbose,
	}, nil
}

func getenvFallback(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return ""
}
