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
	"fmt"
	"os"
	"path"
	"strconv"

	"gopkg.in/yaml.v3"
)

const (
	defaultTerragruntGithubURL = "https://api.github.com/repos/gruntwork-io/terragrunt/releases"
	defaultTofuGithubURL       = "https://api.github.com/repos/opentofu/opentofu/releases"

	autoInstallEnvName = "AUTO_INSTALL"
	forceRemoteEnvName = "FORCE_REMOTE"
	remoteURLEnvName   = "REMOTE"
	rootPathEnvName    = "ROOT"
	verboseEnvName     = "VERBOSE"

	tenvRemoteConf = "TENV_REMOTE_CONF"

	tfenvPrefix              = "TFENV_"
	tfAutoInstallEnvName     = tfenvPrefix + autoInstallEnvName
	tfForceRemoteEnvName     = tfenvPrefix + forceRemoteEnvName
	tfHashicorpPGPKeyEnvName = tfenvPrefix + "HASHICORP_PGP_KEY"
	TfRemoteURLEnvName       = tfenvPrefix + remoteURLEnvName
	tfRootPathEnvName        = tfenvPrefix + rootPathEnvName
	tfVerboseEnvName         = tfenvPrefix + verboseEnvName
	TfVersionEnvName         = tfenvPrefix + "TERRAFORM_VERSION"

	tgPrefix           = "TG_"
	TgRemoteURLEnvName = tgPrefix + remoteURLEnvName
	TgVersionEnvName   = tgPrefix + "VERSION"

	tofuenvPrefix             = "TOFUENV_"
	tofuAutoInstallEnvName    = tofuenvPrefix + autoInstallEnvName
	tofuForceRemoteEnvName    = tofuenvPrefix + forceRemoteEnvName
	tofuOpenTofuPGPKeyEnvName = tofuenvPrefix + "OPENTOFU_PGP_KEY"
	TofuRemoteURLEnvName      = tofuenvPrefix + remoteURLEnvName
	tofuRootPathEnvName       = tofuenvPrefix + rootPathEnvName
	tofuTokenEnvName          = tofuenvPrefix + "GITHUB_TOKEN"
	tofuVerboseEnvName        = tofuenvPrefix + verboseEnvName
	TofuVersionEnvName        = tofuenvPrefix + "TOFU_VERSION"
)

type Config struct {
	ForceRemote    bool
	GithubToken    string
	NoInstall      bool
	RemoteConfPath string
	RootPath       string
	TfKeyPath      string
	TfRemoteURL    string
	TgRemoteURL    string
	TofuKeyPath    string
	TofuRemoteURL  string
	UserPath       string
	Verbose        bool
}

func InitConfigFromEnv() (Config, error) {
	userPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	autoInstall, err := getenvBoolFallback(true, tofuAutoInstallEnvName, tfAutoInstallEnvName)
	if err != nil {
		return Config{}, err
	}

	forceRemote, err := getenvBoolFallback(false, tofuForceRemoteEnvName, tfForceRemoteEnvName)
	if err != nil {
		return Config{}, err
	}

	rootPath := getenvFallback(tofuRootPathEnvName, tfRootPathEnvName)
	if rootPath == "" {
		rootPath = path.Join(userPath, ".tenv")
	}

	verbose, err := getenvBoolFallback(false, tofuVerboseEnvName, tfVerboseEnvName)
	if err != nil {
		return Config{}, err
	}

	return Config{
		ForceRemote:    forceRemote,
		GithubToken:    os.Getenv(tofuTokenEnvName),
		NoInstall:      !autoInstall,
		RemoteConfPath: os.Getenv(tenvRemoteConf),
		RootPath:       rootPath,
		TfKeyPath:      os.Getenv(tfHashicorpPGPKeyEnvName),
		TfRemoteURL:    os.Getenv(TfRemoteURLEnvName),
		TgRemoteURL:    os.Getenv(TgRemoteURLEnvName),
		TofuKeyPath:    os.Getenv(tofuOpenTofuPGPKeyEnvName),
		TofuRemoteURL:  os.Getenv(TofuRemoteURLEnvName),
		UserPath:       userPath,
		Verbose:        verbose,
	}, nil
}

func (conf *Config) ReadRemoteConf(targetName string) map[string]string {
	remoteConfPath := conf.RemoteConfPath
	if remoteConfPath == "" {
		remoteConfPath = path.Join(conf.RootPath, "remote.yaml")
	}

	data, err := os.ReadFile(remoteConfPath)
	if err != nil {
		if conf.Verbose {
			fmt.Println("Can not read remote conf :", err) //nolint
		}

		return nil
	}

	var remoteConf map[string]map[string]string
	if err = yaml.Unmarshal(data, &remoteConf); err != nil {
		if conf.Verbose {
			fmt.Println("Can not parse remote conf :", err) //nolint
		}

		return nil
	}

	return remoteConf[targetName]
}

func MapGetDefault(m map[string]string, key string, defaultValue string) string {
	if value := m[key]; value != "" {
		return value
	}

	return defaultValue
}

func getenvBoolFallback(defaultValue bool, keys ...string) (bool, error) {
	if verboseStr := getenvFallback(keys...); verboseStr != "" {
		return strconv.ParseBool(verboseStr)
	}

	return defaultValue, nil
}

func getenvFallback(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}

	return ""
}
