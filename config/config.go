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
	"github.com/tofuutils/tenv/v4/config/envname"
	configutils "github.com/tofuutils/tenv/v4/config/utils"
	githuburl "github.com/tofuutils/tenv/v4/pkg/github/url"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	atmosurl "github.com/tofuutils/tenv/v4/versionmanager/retriever/atmos/url"
	terraformurl "github.com/tofuutils/tenv/v4/versionmanager/retriever/terraform/url"
	terragrunturl "github.com/tofuutils/tenv/v4/versionmanager/retriever/terragrunt/url"
	terramateurl "github.com/tofuutils/tenv/v4/versionmanager/retriever/terramate/url"
	tofuurl "github.com/tofuutils/tenv/v4/versionmanager/retriever/tofu/url"
)

const (
	defaultDirName = ".tenv"
)

const (
	SignValidation ValidationMode = iota
	ShaValidation
	NoValidation
)

type ValidationMode uint8

func ParseValidationMode(mode string) ValidationMode {
	switch mode {
	case "none":
		return NoValidation
	case "sha":
		return ShaValidation
	default:
		return SignValidation
	}
}

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
	LockPath         string
	remoteConfLoaded bool
	RemoteConfPath   string
	RootPath         string
	SkipInstall      bool
	Tf               RemoteConfig
	TfKeyPathOrURL   string
	Tg               RemoteConfig
	Tm               RemoteConfig
	Tofu             RemoteConfig
	TofuKeyPathOrURL string
	UserPath         string
	Validation       ValidationMode
	WorkPath         string
}

func DefaultConfig() (Config, error) {
	userPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Arch:             runtime.GOARCH,
		Atmos:            makeDefaultRemoteConfig(atmosurl.Github, githuburl.Base),
		Getenv:           EmptyGetenv,
		LockPath:         filepath.Join(userPath, defaultDirName),
		remoteConfLoaded: true,
		RootPath:         filepath.Join(userPath, defaultDirName),
		SkipInstall:      true,
		Tf:               makeDefaultRemoteConfig(terraformurl.Hashicorp, terraformurl.Hashicorp),
		Tg:               makeDefaultRemoteConfig(terragrunturl.Github, githuburl.Base),
		Tm:               makeDefaultRemoteConfig(terramateurl.Github, githuburl.Base),
		Tofu:             makeDefaultRemoteConfig(tofuurl.Github, githuburl.Base),
		UserPath:         userPath,
		WorkPath:         ".",
		Validation:       SignValidation,
		TfKeyPathOrURL:   terraformurl.PublicKey,
		TofuKeyPathOrURL: tofuurl.PublicKey,
	}, nil
}

func InitConfigFromEnv() (Config, error) {
	getenv := configutils.GetenvFunc(os.Getenv)

	userPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	arch := getenv.Fallback(envname.TenvArch, envname.TofuArch, envname.TfArch)
	if arch == "" {
		arch = runtime.GOARCH
	}

	autoInstall, err := getenv.BoolFallback(false, envname.TenvAutoInstall, envname.TofuAutoInstall, envname.TfAutoInstall)
	if err != nil {
		return Config{}, err
	}

	forceRemote, err := getenv.BoolFallback(false, envname.TenvForceRemote, envname.TofuForceRemote, envname.TfForceRemote)
	if err != nil {
		return Config{}, err
	}

	rootPath := getenv.Fallback(envname.TenvRootPath, envname.TofuRootPath, envname.TfRootPath)
	if rootPath == "" {
		rootPath = filepath.Join(userPath, defaultDirName)
	}

	lockPath := getenv.Fallback(envname.TenvLockPath, envname.TofuRootPath, envname.TfRootPath)
	if lockPath == "" {
		lockPath = rootPath // Default to root path for backward compatibility
	}

	quiet, err := getenv.Bool(false, envname.TenvQuiet)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Arch:             arch,
		Atmos:            makeRemoteConfig(getenv, envname.AtmosRemoteURL, envname.AtmosListURL, envname.AtmosInstallMode, envname.AtmosListMode, atmosurl.Github, githuburl.Base),
		ForceQuiet:       quiet,
		ForceRemote:      forceRemote,
		Getenv:           getenv,
		GithubActions:    getenv.Present(envname.GithubActions),
		GithubToken:      getenv.Fallback(envname.TenvToken, envname.TofuToken),
		LockPath:         lockPath,
		RemoteConfPath:   getenv(envname.TenvRemoteConf),
		RootPath:         rootPath,
		SkipInstall:      !autoInstall,
		Tf:               makeRemoteConfig(getenv, envname.TfRemoteURL, envname.TfListURL, envname.TfInstallMode, envname.TfListMode, terraformurl.Hashicorp, terraformurl.Hashicorp),
		TfKeyPathOrURL:   getenv.WithDefault(terraformurl.PublicKey, envname.TfHashicorpPGPKey),
		Tg:               makeRemoteConfig(getenv, envname.TgRemoteURL, envname.TgListURL, envname.TgInstallMode, envname.TgListMode, terragrunturl.Github, githuburl.Base),
		Tm:               makeRemoteConfig(getenv, envname.TmRemoteURL, envname.TmListURL, envname.TmInstallMode, envname.TmListMode, terramateurl.Github, githuburl.Base),
		Tofu:             makeRemoteConfig(getenv, envname.TofuRemoteURL, envname.TofuListURL, envname.TofuInstallMode, envname.TofuListMode, tofuurl.Github, githuburl.Base),
		TofuKeyPathOrURL: getenv.WithDefault(tofuurl.PublicKey, envname.TofuOpenTofuPGPKey),
		UserPath:         userPath,
		Validation:       ParseValidationMode(getenv(envname.TenvValidation)),
		WorkPath:         ".",
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
			if logLevelStr := conf.Getenv(envname.TenvLog); logLevelStr != "" {
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

func (conf *Config) InitValidation(skipSum bool, skipSign bool) {
	switch {
	case skipSum: // higher priority to --skip-sha
		conf.Validation = NoValidation
	case skipSign && conf.Validation != NoValidation:
		conf.Validation = ShaValidation
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
