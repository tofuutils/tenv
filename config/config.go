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
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/pkg/loghelper"
	"gopkg.in/yaml.v3"
)

const (
	TenvName       = "tenv"
	TerraformName  = "terraform"
	TerragruntName = "terragrunt"
	TofuName       = "tofu"

	autoInstallEnvName = "AUTO_INSTALL"
	forceRemoteEnvName = "FORCE_REMOTE"
	installModeEnvName = "INSTALL_MODE"
	listModeEnvName    = "LIST_MODE"
	listURLEnvName     = "LIST_URL"
	logEnvName         = "LOG"
	quietEnvName       = "QUIET"
	remoteURLEnvName   = "REMOTE"
	rootPathEnvName    = "ROOT"

	tenvPrefix            = "TENV_"
	tenvLogEnvName        = tenvPrefix + logEnvName
	tenvQuietEnvName      = tenvPrefix + quietEnvName
	tenvRemoteConfEnvName = tenvPrefix + "REMOTE_CONF"
	tenvRootPathEnvName   = tenvPrefix + rootPathEnvName

	tfenvPrefix              = "TFENV_"
	tfAutoInstallEnvName     = tfenvPrefix + autoInstallEnvName
	tfForceRemoteEnvName     = tfenvPrefix + forceRemoteEnvName
	tfHashicorpPGPKeyEnvName = tfenvPrefix + "HASHICORP_PGP_KEY"
	tfInstallModeEnvName     = tfenvPrefix + installModeEnvName
	tfListModeEnvName        = tfenvPrefix + listModeEnvName
	tfListURLEnvName         = tfenvPrefix + listURLEnvName
	TfRemoteURLEnvName       = tfenvPrefix + remoteURLEnvName
	tfRootPathEnvName        = tfenvPrefix + rootPathEnvName
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
	TofuVersionEnvName        = tofuenvPrefix + "TOFU_VERSION"
)

type Config struct {
	AppLogger        hclog.Logger
	DisplayNormal    bool
	DisplayVerbose   bool
	ForceQuiet       bool
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
}

func InitConfigFromEnv() (Config, error) {
	userPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	logLevel := hclog.Warn
	logLevelStr := os.Getenv(tenvLogEnvName)
	if logLevelStr != "" {
		logLevel = hclog.LevelFromString(logLevelStr)
	}
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name: TenvName, Level: logLevel,
	})

	autoInstall, err := getenvBoolFallback(true, tofuAutoInstallEnvName, tfAutoInstallEnvName)
	if err != nil {
		return Config{}, err
	}

	forceRemote, err := getenvBoolFallback(false, tofuForceRemoteEnvName, tfForceRemoteEnvName)
	if err != nil {
		return Config{}, err
	}

	rootPath := getenvFallback(tenvRootPathEnvName, tofuRootPathEnvName, tfRootPathEnvName)
	if rootPath == "" {
		rootPath = path.Join(userPath, ".tenv")
	}

	quiet, err := getenvBoolFallback(false, tenvQuietEnvName)
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppLogger:      appLogger,
		ForceQuiet:     quiet,
		ForceRemote:    forceRemote,
		GithubToken:    os.Getenv(tofuTokenEnvName),
		NoInstall:      !autoInstall,
		RemoteConfPath: os.Getenv(tenvRemoteConfEnvName),
		RootPath:       rootPath,
		Tf:             makeRemoteConfig(TfRemoteURLEnvName, tfListURLEnvName, tfInstallModeEnvName, tfListModeEnvName, defaultHashicorpURL, defaultHashicorpURL),
		TfKeyPath:      os.Getenv(tfHashicorpPGPKeyEnvName),
		Tg:             makeRemoteConfig(TgRemoteURLEnvName, tgListURLEnvName, tgInstallModeEnvName, tgListModeEnvName, defaultTerragruntGithubURL, baseGithubURL),
		Tofu:           makeRemoteConfig(TofuRemoteURLEnvName, tofuListURLEnvName, tofuInstallModeEnvName, tofuListModeEnvName, defaultTofuGithubURL, baseGithubURL),
		TofuKeyPath:    os.Getenv(tofuOpenTofuPGPKeyEnvName),
		UserPath:       userPath,
	}, nil
}

func (conf *Config) InitRemoteConf() error {
	if conf.remoteConfLoaded {
		return nil
	}
	conf.remoteConfLoaded = true

	remoteConfPath := conf.RemoteConfPath
	if remoteConfPath == "" {
		remoteConfPath = path.Join(conf.RootPath, "remote.yaml")
	}

	data, err := os.ReadFile(remoteConfPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		conf.AppLogger.Info("Can not read remote configuration file", loghelper.Error, err)

		return nil
	}

	var remoteConf map[string]map[string]string
	if err = yaml.Unmarshal(data, &remoteConf); err != nil {
		return err
	}

	conf.Tf.Data = remoteConf[TerraformName]
	conf.Tg.Data = remoteConf[TerragruntName]
	conf.Tofu.Data = remoteConf[TofuName]

	return nil
}

func (conf *Config) LogLevelUpdate() {
	if conf.ForceQuiet {
		conf.DisplayVerbose = false
		conf.AppLogger.SetLevel(hclog.Off)
	} else {
		conf.DisplayNormal = true
		if conf.DisplayVerbose {
			conf.AppLogger.SetLevel(hclog.Trace)
		}
	}
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
