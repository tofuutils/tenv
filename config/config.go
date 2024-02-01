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
	autoInstallEnvName = "AUTO_INSTALL"
	forceRemoteEnvName = "FORCE_REMOTE"
	installModeEnvName = "INSTALL_MODE"
	listModeEnvName    = "LIST_MODE"
	listURLEnvName     = "LIST_URL"
	remoteURLEnvName   = "REMOTE"
	rootPathEnvName    = "ROOT"
	verboseEnvName     = "VERBOSE"

	tenvPrefix            = "TENV_"
	tenvRemoteConfEnvName = tenvPrefix + "REMOTE_CONF"
	tenvVerboseEnvName    = tenvPrefix + verboseEnvName

	tfenvPrefix              = "TFENV_"
	tfAutoInstallEnvName     = tfenvPrefix + autoInstallEnvName
	tfForceRemoteEnvName     = tfenvPrefix + forceRemoteEnvName
	tfHashicorpPGPKeyEnvName = tfenvPrefix + "HASHICORP_PGP_KEY"
	tfInstallModeEnvName     = tfenvPrefix + installModeEnvName
	tfListModeEnvName        = tfenvPrefix + listModeEnvName
	tfListURLEnvName         = tfenvPrefix + listURLEnvName
	TfRemoteURLEnvName       = tfenvPrefix + remoteURLEnvName
	tfRootPathEnvName        = tfenvPrefix + rootPathEnvName
	tfVerboseEnvName         = tfenvPrefix + verboseEnvName
	TfVersionEnvName         = tfenvPrefix + "TERRAFORM_VERSION"

	tgPrefix             = "TG_"
	tgInstallModeEnvName = tgPrefix + installModeEnvName
	tgListModeEnvName    = tgPrefix + listModeEnvName
	tgListURLEnvName     = tgPrefix + listURLEnvName
	TgRemoteURLEnvName   = tgPrefix + remoteURLEnvName
	TgVersionEnvName     = tgPrefix + "VERSION"

	tofuenvPrefix             = "TOFUENV_"
	tofuAutoInstallEnvName    = tofuenvPrefix + autoInstallEnvName
	tofuForceRemoteEnvName    = tofuenvPrefix + forceRemoteEnvName
	tofuInstallModeEnvName    = tofuenvPrefix + installModeEnvName
	tofuListModeEnvName       = tofuenvPrefix + listModeEnvName
	tofuListURLEnvName        = tofuenvPrefix + listURLEnvName
	tofuOpenTofuPGPKeyEnvName = tofuenvPrefix + "OPENTOFU_PGP_KEY"
	TofuRemoteURLEnvName      = tofuenvPrefix + remoteURLEnvName
	tofuRootPathEnvName       = tofuenvPrefix + rootPathEnvName
	tofuTokenEnvName          = tofuenvPrefix + "GITHUB_TOKEN"
	tofuVerboseEnvName        = tofuenvPrefix + verboseEnvName
	TofuVersionEnvName        = tofuenvPrefix + "TOFU_VERSION"
)

type Config struct {
	ForceRemote      bool
	GithubToken      string
	NoInstall        bool
	remoteConfLoaded bool
	RemoteConfPath   string
	RootPath         string
	Tf               RemoteConfig
	TfKeyPath        string
	Tg               RemoteConfig
	Tofu             RemoteConfig
	TofuKeyPath      string
	UserPath         string
	Verbose          bool
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

	verbose, err := getenvBoolFallback(false, tenvVerboseEnvName, tofuVerboseEnvName, tfVerboseEnvName)
	if err != nil {
		return Config{}, err
	}

	return Config{
		ForceRemote:    forceRemote,
		GithubToken:    os.Getenv(tofuTokenEnvName),
		NoInstall:      !autoInstall,
		RemoteConfPath: os.Getenv(tenvRemoteConfEnvName),
		RootPath:       rootPath,
		Tf:             makeRemoteConfig(TfRemoteURLEnvName, tfListURLEnvName, tfInstallModeEnvName, tfListModeEnvName, defaultHashicorpURL),
		TfKeyPath:      os.Getenv(tfHashicorpPGPKeyEnvName),
		Tg:             makeRemoteConfig(TgRemoteURLEnvName, tgListURLEnvName, tgInstallModeEnvName, tgListModeEnvName, defaultTerragruntGithubURL),
		Tofu:           makeRemoteConfig(TofuRemoteURLEnvName, tofuListURLEnvName, tofuInstallModeEnvName, tofuListModeEnvName, defaultTofuGithubURL),
		TofuKeyPath:    os.Getenv(tofuOpenTofuPGPKeyEnvName),
		UserPath:       userPath,
		Verbose:        verbose,
	}, nil
}

func (conf *Config) InitRemoteConf() {
	if conf.remoteConfLoaded {
		return
	}
	conf.remoteConfLoaded = true

	remoteConfPath := conf.RemoteConfPath
	if remoteConfPath == "" {
		remoteConfPath = path.Join(conf.RootPath, "remote.yaml")
	}

	data, err := os.ReadFile(remoteConfPath)
	if err != nil {
		if conf.Verbose {
			fmt.Println("Can not read remote conf :", err) //nolint
		}

		return
	}

	var remoteConf map[string]map[string]string
	if err = yaml.Unmarshal(data, &remoteConf); err != nil {
		if conf.Verbose {
			fmt.Println("Can not parse remote conf :", err) //nolint
		}

		return
	}

	conf.Tf.Data = remoteConf["terraform"]
	conf.Tg.Data = remoteConf["terragrunt"]
	conf.Tofu.Data = remoteConf["tofu"]

	return
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
