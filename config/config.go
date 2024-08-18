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
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"

	"github.com/tofuutils/tenv/v3/config/cmdconst"
	configutils "github.com/tofuutils/tenv/v3/config/utils"
	"github.com/tofuutils/tenv/v3/pkg/loghelper"
)

const (
	githubActionsEnvName = "GITHUB_ACTIONS"

	archEnvName        = "ARCH"
	autoInstallEnvName = "AUTO_INSTALL"
	defaultConstraint  = "DEFAULT_CONSTRAINT"
	defaultVersion     = "DEFAULT_" + version
	forceRemoteEnvName = "FORCE_REMOTE"
	installModeEnvName = "INSTALL_MODE"
	listModeEnvName    = "LIST_MODE"
	listURLEnvName     = "LIST_URL"
	logEnvName         = "LOG"
	quietEnvName       = "QUIET"
	remotePassEnvName  = "REMOTE_PASSWORD"
	remoteURLEnvName   = "REMOTE"
	remoteUserEnvName  = "REMOTE_USER"
	rootPathEnvName    = "ROOT"
	tokenEnvName       = "GITHUB_TOKEN" //nolint
	version            = "VERSION"

	atmosPrefix                   = "ATMOS_"
	AtmosDefaultConstraintEnvName = atmosPrefix + defaultConstraint
	AtmosDefaultVersionEnvName    = atmosPrefix + defaultVersion
	atmosInstallModeEnvName       = atmosPrefix + installModeEnvName
	atmosListModeEnvName          = atmosPrefix + listModeEnvName
	atmosListURLEnvName           = atmosPrefix + listURLEnvName
	AtmosRemotePassEnvName        = atmosPrefix + remotePassEnvName
	AtmosRemoteURLEnvName         = atmosPrefix + remoteURLEnvName
	AtmosRemoteUserEnvName        = atmosPrefix + remoteUserEnvName
	AtmosVersionEnvName           = atmosPrefix + version

	tenvPrefix             = "TENV_"
	tenvArchEnvName        = tenvPrefix + archEnvName
	tenvAutoInstallEnvName = tenvPrefix + autoInstallEnvName
	tenvForceRemoteEnvName = tenvPrefix + forceRemoteEnvName
	tenvLogEnvName         = tenvPrefix + logEnvName
	tenvQuietEnvName       = tenvPrefix + quietEnvName
	tenvRemoteConfEnvName  = tenvPrefix + "REMOTE_CONF"
	tenvRootPathEnvName    = tenvPrefix + rootPathEnvName
	tenvTokenEnvName       = tenvPrefix + tokenEnvName

	tfenvPrefix                = "TFENV_"
	tfenvTerraformPrefix       = tfenvPrefix + "TERRAFORM_"
	tfArchEnvName              = tfenvPrefix + archEnvName
	tfAutoInstallEnvName       = tfenvPrefix + autoInstallEnvName
	TfDefaultConstraintEnvName = tfenvTerraformPrefix + defaultConstraint
	TfDefaultVersionEnvName    = tfenvTerraformPrefix + defaultVersion
	tfForceRemoteEnvName       = tfenvPrefix + forceRemoteEnvName
	tfHashicorpPGPKeyEnvName   = tfenvPrefix + "HASHICORP_PGP_KEY"
	tfInstallModeEnvName       = tfenvPrefix + installModeEnvName
	tfListModeEnvName          = tfenvPrefix + listModeEnvName
	tfListURLEnvName           = tfenvPrefix + listURLEnvName
	TfRemotePassEnvName        = tfenvPrefix + remotePassEnvName
	TfRemoteURLEnvName         = tfenvPrefix + remoteURLEnvName
	TfRemoteUserEnvName        = tfenvPrefix + remoteUserEnvName
	tfRootPathEnvName          = tfenvPrefix + rootPathEnvName
	TfVersionEnvName           = tfenvTerraformPrefix + version

	tgPrefix                   = "TG_"
	TgDefaultConstraintEnvName = tgPrefix + defaultConstraint
	TgDefaultVersionEnvName    = tgPrefix + defaultVersion
	tgInstallModeEnvName       = tgPrefix + installModeEnvName
	tgListModeEnvName          = tgPrefix + listModeEnvName
	tgListURLEnvName           = tgPrefix + listURLEnvName
	TgRemotePassEnvName        = tgPrefix + remotePassEnvName
	TgRemoteURLEnvName         = tgPrefix + remoteURLEnvName
	TgRemoteUserEnvName        = tgPrefix + remoteUserEnvName
	TgVersionEnvName           = tgPrefix + version

	tofuenvPrefix                = "TOFUENV_"
	tofuenvTofuPrefix            = tofuenvPrefix + "TOFU_"
	tofuArchEnvName              = tofuenvPrefix + archEnvName
	tofuAutoInstallEnvName       = tofuenvPrefix + autoInstallEnvName
	TofuDefaultConstraintEnvName = tofuenvTofuPrefix + defaultConstraint
	TofuDefaultVersionEnvName    = tofuenvTofuPrefix + defaultVersion
	tofuForceRemoteEnvName       = tofuenvPrefix + forceRemoteEnvName
	tofuInstallModeEnvName       = tofuenvPrefix + installModeEnvName
	tofuListModeEnvName          = tofuenvPrefix + listModeEnvName
	tofuListURLEnvName           = tofuenvPrefix + listURLEnvName
	tofuOpenTofuPGPKeyEnvName    = tofuenvPrefix + "OPENTOFU_PGP_KEY"
	TofuRemotePassEnvName        = tofuenvPrefix + remotePassEnvName
	TofuRemoteURLEnvName         = tofuenvPrefix + remoteURLEnvName
	TofuRemoteUserEnvName        = tofuenvPrefix + remoteUserEnvName
	tofuRootPathEnvName          = tofuenvPrefix + rootPathEnvName
	tofuTokenEnvName             = tofuenvPrefix + tokenEnvName
	TofuURLTemplateEnvName       = tofuenvPrefix + "URL_TEMPLATE"
	TofuVersionEnvName           = tofuenvTofuPrefix + version
)

type Config struct {
	Arch             string
	Atmos            RemoteConfig
	Displayer        loghelper.Displayer
	DisplayVerbose   bool
	ForceQuiet       bool
	ForceRemote      bool
	GithubActions    bool
	GithubToken      string
	NoInstall        bool
	remoteConfLoaded bool
	RemoteConfPath   string
	RootPath         string
	SkipSignature    bool
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

	arch := configutils.GetenvFallback(tenvArchEnvName, tofuArchEnvName, tfArchEnvName)
	if arch == "" {
		arch = runtime.GOARCH
	}

	autoInstall, err := configutils.GetenvBoolFallback(false, tenvAutoInstallEnvName, tofuAutoInstallEnvName, tfAutoInstallEnvName)
	if err != nil {
		return Config{}, err
	}

	forceRemote, err := configutils.GetenvBoolFallback(false, tenvForceRemoteEnvName, tofuForceRemoteEnvName, tfForceRemoteEnvName)
	if err != nil {
		return Config{}, err
	}

	rootPath := configutils.GetenvFallback(tenvRootPathEnvName, tofuRootPathEnvName, tfRootPathEnvName)
	if rootPath == "" {
		rootPath = filepath.Join(userPath, ".tenv")
	}

	quiet, err := configutils.GetenvBoolFallback(false, tenvQuietEnvName)
	if err != nil {
		return Config{}, err
	}

	gha, err := configutils.GetenvBool(false, githubActionsEnvName)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Arch:           arch,
		Atmos:          makeRemoteConfig(AtmosRemoteURLEnvName, atmosListURLEnvName, atmosInstallModeEnvName, atmosListModeEnvName, defaultAtmosGithubURL, baseGithubURL),
		ForceQuiet:     quiet,
		ForceRemote:    forceRemote,
		GithubActions:  gha,
		GithubToken:    configutils.GetenvFallback(tenvTokenEnvName, tofuTokenEnvName),
		NoInstall:      !autoInstall,
		RemoteConfPath: os.Getenv(tenvRemoteConfEnvName),
		RootPath:       rootPath,
		Tf:             makeRemoteConfig(TfRemoteURLEnvName, tfListURLEnvName, tfInstallModeEnvName, tfListModeEnvName, defaultHashicorpURL, defaultHashicorpURL),
		TfKeyPath:      os.Getenv(tfHashicorpPGPKeyEnvName),
		Tg:             makeRemoteConfig(TgRemoteURLEnvName, tgListURLEnvName, tgInstallModeEnvName, tgListModeEnvName, defaultTerragruntGithubURL, baseGithubURL),
		Tofu:           makeRemoteConfig(TofuRemoteURLEnvName, tofuListURLEnvName, tofuInstallModeEnvName, tofuListModeEnvName, DefaultTofuGithubURL, baseGithubURL),
		TofuKeyPath:    os.Getenv(tofuOpenTofuPGPKeyEnvName),
		UserPath:       userPath,
	}, nil
}

func (conf *Config) InitDisplayer(proxyCall bool) {
	if conf.ForceQuiet {
		conf.Displayer = loghelper.InertDisplayer
		conf.DisplayVerbose = false
	} else {
		logLevel := hclog.Trace
		if !conf.DisplayVerbose {
			logLevel = hclog.Warn
			if logLevelStr := os.Getenv(tenvLogEnvName); logLevelStr != "" {
				logLevel = hclog.LevelFromString(logLevelStr)
			}
		}
		appLogger := hclog.New(&hclog.LoggerOptions{
			Name: cmdconst.TenvName, Level: logLevel,
		})

		if proxyCall {
			display := loghelper.BuildDisplayFunc(os.Stderr, color.New(color.FgGreen))
			conf.Displayer = loghelper.NewRecordingDisplayer(loghelper.MakeBasicDisplayer(appLogger, display))
		} else {
			conf.Displayer = loghelper.MakeBasicDisplayer(appLogger, loghelper.StdDisplay)
		}
	}
}

func (conf *Config) InitInstall(forceInstall bool, forceNoInstall bool) {
	switch {
	case forceNoInstall: // higher priority to --no-install
		conf.NoInstall = true
	case forceInstall:
		conf.NoInstall = false
	}
}

func (conf *Config) InitRemoteConf() error {
	if conf.remoteConfLoaded {
		return nil
	}
	conf.remoteConfLoaded = true

	remoteConfPath := conf.RemoteConfPath
	if remoteConfPath == "" {
		remoteConfPath = filepath.Join(conf.RootPath, "remote.yaml")
	}

	data, err := os.ReadFile(remoteConfPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		conf.Displayer.Log(hclog.Debug, "Can not read remote configuration file", loghelper.Error, err)

		return nil
	}

	var remoteConf map[string]map[string]string
	if err = yaml.Unmarshal(data, &remoteConf); err != nil {
		return err
	}

	conf.Tf.Data = remoteConf[cmdconst.TerraformName]
	conf.Tg.Data = remoteConf[cmdconst.TerragruntName]
	conf.Tofu.Data = remoteConf[cmdconst.TofuName]
	conf.Atmos.Data = remoteConf[cmdconst.AtmosName]

	return nil
}
