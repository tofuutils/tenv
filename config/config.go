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

	"github.com/tofuutils/tenv/v4/config/cmdconst"
	configutils "github.com/tofuutils/tenv/v4/config/utils"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

const (
	defaultDirName = ".tenv"

	archEnvName        = "ARCH"
	autoInstallEnvName = "AUTO_INSTALL"
	DefaultConstraint  = "DEFAULT_CONSTRAINT"
	DefaultVersion     = "DEFAULT_" + Version
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
	Version            = "VERSION"

	CiEnvName         = "CI"
	PipelineWsEnvName = "PIPELINE_WORKSPACE"

	githubPrefix         = "GITHUB_"
	githubActionsEnvName = githubPrefix + "ACTIONS"
	GithubOutputEnvName  = githubPrefix + "OUTPUT"
	tokenEnvName         = githubPrefix + "TOKEN"

	AtmosPrefix             = "ATMOS_"
	atmosInstallModeEnvName = AtmosPrefix + installModeEnvName
	atmosListModeEnvName    = AtmosPrefix + listModeEnvName
	atmosListURLEnvName     = AtmosPrefix + listURLEnvName
	AtmosRemotePassEnvName  = AtmosPrefix + remotePassEnvName
	AtmosRemoteURLEnvName   = AtmosPrefix + remoteURLEnvName
	AtmosRemoteUserEnvName  = AtmosPrefix + remoteUserEnvName

	tenvPrefix               = "TENV_"
	tenvArchEnvName          = tenvPrefix + archEnvName
	tenvAutoInstallEnvName   = tenvPrefix + autoInstallEnvName
	TenvDetachedProxyEnvName = tenvPrefix + "DETACHED_PROXY"
	tenvForceRemoteEnvName   = tenvPrefix + forceRemoteEnvName
	tenvLogEnvName           = tenvPrefix + logEnvName
	tenvQuietEnvName         = tenvPrefix + quietEnvName
	tenvRemoteConfEnvName    = tenvPrefix + "REMOTE_CONF"
	tenvRootPathEnvName      = tenvPrefix + rootPathEnvName
	tenvTokenEnvName         = tenvPrefix + tokenEnvName

	TfenvPrefix              = "TFENV_"
	tfenvTerraformPrefix     = TfenvPrefix + "TERRAFORM_"
	tfArchEnvName            = TfenvPrefix + archEnvName
	tfAutoInstallEnvName     = TfenvPrefix + autoInstallEnvName
	tfForceRemoteEnvName     = TfenvPrefix + forceRemoteEnvName
	tfHashicorpPGPKeyEnvName = TfenvPrefix + "HASHICORP_PGP_KEY"
	tfInstallModeEnvName     = TfenvPrefix + installModeEnvName
	tfListModeEnvName        = TfenvPrefix + listModeEnvName
	tfListURLEnvName         = TfenvPrefix + listURLEnvName
	TfRemotePassEnvName      = TfenvPrefix + remotePassEnvName
	TfRemoteURLEnvName       = TfenvPrefix + remoteURLEnvName
	TfRemoteUserEnvName      = TfenvPrefix + remoteUserEnvName
	tfRootPathEnvName        = TfenvPrefix + rootPathEnvName

	TgPrefix             = "TG_"
	tgInstallModeEnvName = TgPrefix + installModeEnvName
	tgListModeEnvName    = TgPrefix + listModeEnvName
	tgListURLEnvName     = TgPrefix + listURLEnvName
	TgRemotePassEnvName  = TgPrefix + remotePassEnvName
	TgRemoteURLEnvName   = TgPrefix + remoteURLEnvName
	TgRemoteUserEnvName  = TgPrefix + remoteUserEnvName

	TofuenvPrefix             = "TOFUENV_"
	tofuenvTofuPrefix         = TofuenvPrefix + "TOFU_"
	tofuArchEnvName           = TofuenvPrefix + archEnvName
	tofuAutoInstallEnvName    = TofuenvPrefix + autoInstallEnvName
	tofuForceRemoteEnvName    = TofuenvPrefix + forceRemoteEnvName
	tofuInstallModeEnvName    = TofuenvPrefix + installModeEnvName
	tofuListModeEnvName       = TofuenvPrefix + listModeEnvName
	tofuListURLEnvName        = TofuenvPrefix + listURLEnvName
	tofuOpenTofuPGPKeyEnvName = TofuenvPrefix + "OPENTOFU_PGP_KEY"
	TofuRemotePassEnvName     = TofuenvPrefix + remotePassEnvName
	TofuRemoteURLEnvName      = TofuenvPrefix + remoteURLEnvName
	TofuRemoteUserEnvName     = TofuenvPrefix + remoteUserEnvName
	tofuRootPathEnvName       = TofuenvPrefix + rootPathEnvName
	tofuTokenEnvName          = TofuenvPrefix + tokenEnvName
	TofuURLTemplateEnvName    = TofuenvPrefix + "URL_TEMPLATE"
)

type Config struct {
	Arch             string
	Atmos            RemoteConfig
	Displayer        loghelper.Displayer
	DisplayVerbose   bool
	ForceQuiet       bool
	ForceRemote      bool
	Getenv           configutils.GetenvFunc
	GithubActions    bool
	GithubToken      string
	remoteConfLoaded bool
	RemoteConfPath   string
	RootPath         string
	SkipInstall      bool
	SkipSignature    bool
	Tf               RemoteConfig
	TfKeyPath        string
	Tg               RemoteConfig
	Tofu             RemoteConfig
	TofuKeyPath      string
	UserPath         string
	WorkPath         string
}

func DefaultConfig() (Config, error) {
	userPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Arch:             runtime.GOARCH,
		Atmos:            makeDefaultRemoteConfig(defaultAtmosGithubURL, baseGithubURL),
		Getenv:           EmptyGetenv,
		remoteConfLoaded: true,
		RootPath:         filepath.Join(userPath, defaultDirName),
		SkipInstall:      true,
		Tf:               makeDefaultRemoteConfig(defaultHashicorpURL, defaultHashicorpURL),
		Tg:               makeDefaultRemoteConfig(defaultTerragruntGithubURL, baseGithubURL),
		Tofu:             makeDefaultRemoteConfig(DefaultTofuGithubURL, baseGithubURL),
		UserPath:         userPath,
		WorkPath:         ".",
	}, nil
}

func InitConfigFromEnv() (Config, error) {
	getenv := configutils.GetenvFunc(os.Getenv)

	userPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	arch := getenv.Fallback(tenvArchEnvName, tofuArchEnvName, tfArchEnvName)
	if arch == "" {
		arch = runtime.GOARCH
	}

	autoInstall, err := getenv.BoolFallback(false, tenvAutoInstallEnvName, tofuAutoInstallEnvName, tfAutoInstallEnvName)
	if err != nil {
		return Config{}, err
	}

	forceRemote, err := getenv.BoolFallback(false, tenvForceRemoteEnvName, tofuForceRemoteEnvName, tfForceRemoteEnvName)
	if err != nil {
		return Config{}, err
	}

	rootPath := getenv.Fallback(tenvRootPathEnvName, tofuRootPathEnvName, tfRootPathEnvName)
	if rootPath == "" {
		rootPath = filepath.Join(userPath, defaultDirName)
	}

	quiet, err := getenv.Bool(false, tenvQuietEnvName)
	if err != nil {
		return Config{}, err
	}

	gha, err := getenv.Bool(false, githubActionsEnvName)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Arch:           arch,
		Atmos:          makeRemoteConfig(getenv, AtmosRemoteURLEnvName, atmosListURLEnvName, atmosInstallModeEnvName, atmosListModeEnvName, defaultAtmosGithubURL, baseGithubURL),
		ForceQuiet:     quiet,
		ForceRemote:    forceRemote,
		Getenv:         getenv,
		GithubActions:  gha,
		GithubToken:    getenv.Fallback(tenvTokenEnvName, tofuTokenEnvName),
		RemoteConfPath: getenv(tenvRemoteConfEnvName),
		RootPath:       rootPath,
		SkipInstall:    !autoInstall,
		Tf:             makeRemoteConfig(getenv, TfRemoteURLEnvName, tfListURLEnvName, tfInstallModeEnvName, tfListModeEnvName, defaultHashicorpURL, defaultHashicorpURL),
		TfKeyPath:      getenv(tfHashicorpPGPKeyEnvName),
		Tg:             makeRemoteConfig(getenv, TgRemoteURLEnvName, tgListURLEnvName, tgInstallModeEnvName, tgListModeEnvName, defaultTerragruntGithubURL, baseGithubURL),
		Tofu:           makeRemoteConfig(getenv, TofuRemoteURLEnvName, tofuListURLEnvName, tofuInstallModeEnvName, tofuListModeEnvName, DefaultTofuGithubURL, baseGithubURL),
		TofuKeyPath:    getenv(tofuOpenTofuPGPKeyEnvName),
		UserPath:       userPath,
		WorkPath:       ".",
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
			if logLevelStr := conf.Getenv(tenvLogEnvName); logLevelStr != "" {
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
		conf.SkipInstall = true
	case forceInstall:
		conf.SkipInstall = false
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

func EmptyGetenv(_ string) string {
	return ""
}
