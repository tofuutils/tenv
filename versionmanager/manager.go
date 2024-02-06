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

package versionmanager

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
	"github.com/tofuutils/tenv/pkg/reversecmp"
	"github.com/tofuutils/tenv/versionmanager/semantic"
)

var (
	errEmptyVersion = errors.New("empty version")
	errNoCompatible = errors.New("no compatible version found")
)

type ReleaseInfoRetriever interface {
	InstallRelease(version string, targetPath string) error
	ListReleases() ([]string, error)
}

type VersionManager struct {
	conf             *config.Config
	FolderName       string
	predicateReaders []func(*config.Config) (func(string) bool, bool, error)
	retriever        ReleaseInfoRetriever
	VersionEnvName   string
	VersionFiles     []semantic.VersionFile
}

func MakeVersionManager(conf *config.Config, folderName string, predicateReaders []func(*config.Config) (func(string) bool, bool, error), retriever ReleaseInfoRetriever, versionEnvName string, versionFiles []semantic.VersionFile) VersionManager {
	return VersionManager{conf: conf, FolderName: folderName, predicateReaders: predicateReaders, retriever: retriever, VersionEnvName: versionEnvName, VersionFiles: versionFiles}
}

// detect version (can install depending on auto install env var).
func (m VersionManager) Detect() (string, error) {
	configVersion, err := m.Resolve(semantic.LatestAllowedKey)
	if err != nil {
		return "", err
	}

	return m.detect(configVersion)
}

func (m VersionManager) Install(requestedVersion string) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		return m.installSpecificVersion(parsedVersion.String())
	}

	predicate, reverseOrder, err := semantic.ParsePredicate(requestedVersion, m.FolderName, m.predicateReaders, m.conf)
	if err != nil {
		return err
	}

	// noInstall is set to false to force install regardless of conf
	_, err = m.searchInstallRemote(predicate, reverseOrder, false)

	return err
}

// try to ensure the directory exists with a MkdirAll call.
// (made lazy method : not always useful and allows flag override for root path).
func (m VersionManager) InstallPath() string {
	dir := filepath.Join(m.conf.RootPath, m.FolderName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		m.conf.AppLogger.Warn("Can not create installation directory", loghelper.Error, err)
	}

	return dir
}

func (m VersionManager) ListLocal(reverseOrder bool) ([]string, error) {
	entries, err := os.ReadDir(m.InstallPath())
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	cmpFunc := reversecmp.Reverser[string](semantic.CmpVersion, reverseOrder)
	slices.SortFunc(versions, cmpFunc)

	return versions, nil
}

func (m VersionManager) ListRemote(reverseOrder bool) ([]string, error) {
	versions, err := m.retriever.ListReleases()
	if err != nil {
		return nil, err
	}

	cmpFunc := reversecmp.Reverser[string](semantic.CmpVersion, reverseOrder)
	slices.SortFunc(versions, cmpFunc)

	return versions, nil
}

func (m VersionManager) LocalSet() map[string]struct{} {
	entries, err := os.ReadDir(m.InstallPath())
	if err != nil {
		m.conf.AppLogger.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Can not read installed versions", loghelper.Error, err)

		return nil
	}

	versionSet := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			versionSet[entry.Name()] = struct{}{}
		}
	}

	return versionSet
}

func (m VersionManager) Reset() error {
	versionFilePath := m.RootVersionFilePath()
	err := os.RemoveAll(versionFilePath)
	if err == nil && m.conf.DisplayNormal {
		fmt.Println("Removed", versionFilePath) //nolint
	}

	return err
}

// (made lazy method : not always useful and allows flag override for root path).
func (m VersionManager) Resolve(defaultStrategy string) (string, error) {
	if forcedVersion := os.Getenv(m.VersionEnvName); forcedVersion != "" {
		if m.conf.DisplayNormal {
			fmt.Println("Resolved version from", m.VersionEnvName, ":", color.GreenString(forcedVersion)) //nolint
		}

		return forcedVersion, nil
	}

	if version, err := semantic.RetrieveVersion(m.VersionFiles, m.RootVersionFilePath(), m.conf); err != nil || version != "" {
		return version, err
	}

	if m.conf.DisplayNormal {
		fmt.Println("No", m.FolderName, "version found in basic files, fallback to", color.GreenString(defaultStrategy), "strategy") //nolint
	}

	return defaultStrategy, nil
}

// (made lazy method : not always useful and allows flag override for root path).
func (m VersionManager) RootVersionFilePath() string {
	return filepath.Join(m.conf.RootPath, m.FolderName, "version")
}

func (m VersionManager) Uninstall(requestedVersion string) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err != nil {
		return err
	}

	cleanedVersion := parsedVersion.String()
	targetPath := filepath.Join(m.InstallPath(), cleanedVersion)
	if err = os.RemoveAll(targetPath); err == nil && m.conf.DisplayNormal {
		fmt.Println("Uninstallation of", m.FolderName, cleanedVersion, "successful (directory", targetPath, "removed)") //nolint
	}

	return err
}

func (m VersionManager) Use(requestedVersion string, workingDir bool) error {
	detectedVersion, err := m.detect(requestedVersion)
	if err != nil {
		return err
	}

	targetFilePath := m.VersionFiles[0].Name
	if !workingDir {
		targetFilePath = m.RootVersionFilePath()
	}
	if err = os.WriteFile(targetFilePath, []byte(detectedVersion), 0644); err == nil && m.conf.DisplayNormal {
		fmt.Println("Written", detectedVersion, "in", targetFilePath) //nolint
	}

	return err
}

func (m VersionManager) detect(requestedVersion string) (string, error) {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		cleanedVersion := parsedVersion.String()
		if m.conf.NoInstall {
			return cleanedVersion, nil
		}

		return cleanedVersion, m.installSpecificVersion(cleanedVersion)
	}

	predicate, reverseOrder, err := semantic.ParsePredicate(requestedVersion, m.FolderName, m.predicateReaders, m.conf)
	if err != nil {
		return "", err
	}

	if !m.conf.ForceRemote {
		versions, err := m.ListLocal(reverseOrder)
		if err != nil {
			return "", err
		}

		for _, version := range versions {
			if predicate(version) {
				if m.conf.DisplayNormal {
					fmt.Println("Found compatible version installed locally :", color.GreenString(version)) //nolint
				}

				return version, nil
			}
		}

		if m.conf.NoInstall {
			return "", errNoCompatible
		} else if m.conf.DisplayNormal {
			fmt.Println("No compatible version found locally, search a remote one...") //nolint
		}
	}

	return m.searchInstallRemote(predicate, reverseOrder, m.conf.NoInstall)
}

func (m VersionManager) installSpecificVersion(version string) error {
	if version == "" {
		return errEmptyVersion
	}

	installPath := m.InstallPath()
	entries, err := os.ReadDir(installPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() && version == entry.Name() {
			if m.conf.DisplayNormal {
				fmt.Println(m.FolderName, version, "already installed") //nolint
			}

			return nil
		}
	}

	if m.conf.DisplayNormal {
		fmt.Println("Installing", m.FolderName, version) //nolint
	}
	if err = m.retriever.InstallRelease(version, filepath.Join(installPath, version)); err == nil && m.conf.DisplayNormal {
		fmt.Println("Installation of", m.FolderName, version, "successful") //nolint
	}

	return err
}

func (m VersionManager) searchInstallRemote(predicate func(string) bool, reverseOrder bool, noInstall bool) (string, error) {
	versions, err := m.ListRemote(reverseOrder)
	if err != nil {
		return "", err
	}

	for _, version := range versions {
		if predicate(version) {
			if m.conf.DisplayNormal {
				fmt.Println("Found compatible version remotely :", color.GreenString(version)) //nolint
			}
			if noInstall {
				return version, nil
			}

			return version, m.installSpecificVersion(version)
		}
	}

	return "", errNoCompatible
}
