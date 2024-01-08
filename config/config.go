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
	defaultAutoInstall = true
	defaultRemoteUrl   = "https://github.com/opentofu/opentofu/releases"
)

const (
	envPrefix = "GOTOFUENV_"

	autoInstallEnvName = envPrefix + "AUTO_INSTALL"
	debugLevelEnvName  = envPrefix + "DEBUG"
	remoteUrlEnvName   = envPrefix + "REMOTE"
	rootPathEnvName    = envPrefix + "ROOT"
	tokenEnvName       = envPrefix + "GITHUB_TOKEN"
	versionEnvName     = envPrefix + "TOFU_VERSION"
)

type Config struct {
	AutoInstall bool
	DebugLevel  int
	RemoteUrl   string
	RootPath    string
	Token       string
	Version     string
}

func InitConfig() (Config, error) {
	autoInstall := defaultAutoInstall
	autoInstallStr := os.Getenv(autoInstallEnvName)
	if autoInstallStr != "" {
		var err error
		autoInstall, err = strconv.ParseBool(autoInstallStr)
		if err != nil {
			return Config{}, err
		}
	}

	debugLevel := 0
	debugLevelStr := os.Getenv(debugLevelEnvName)
	if debugLevelStr != "" {
		var err error
		debugLevel, err = strconv.Atoi(debugLevelStr)
		if err != nil {
			return Config{}, err
		}
	}

	remoteUrl := os.Getenv(remoteUrlEnvName)
	if remoteUrl == "" {
		remoteUrl = defaultRemoteUrl
	}

	rootPath := os.Getenv(rootPathEnvName)
	if rootPath == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return Config{}, err
		}
		rootPath = path.Join(userHome, ".gotofuenv")
	}

	return Config{
		AutoInstall: autoInstall,
		DebugLevel:  debugLevel,
		RemoteUrl:   remoteUrl,
		RootPath:    rootPath,
		Token:       os.Getenv(tokenEnvName),
		Version:     os.Getenv(versionEnvName),
	}, nil
}
