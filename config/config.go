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
	autoInstallEnvName = "GOTOFUENV_AUTO_INSTALL"
	debugLevelEnvName  = "GOTOFUENV_DEBUG"
	remoteUrlEnvName   = "GOTOFUENV_REMOTE"
	rootPathEnvName    = "GOTOFUENV_ROOT"
	tokenEnvName       = "GOTOFUENV_GITHUB_TOKEN"
	versionEnvName     = "GOTOFUENV_TOFU_VERSION"
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
