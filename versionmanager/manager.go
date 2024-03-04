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
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
	"github.com/tofuutils/tenv/pkg/reversecmp"
	"github.com/tofuutils/tenv/versionmanager/semantic"
	flatparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/flat"
	"github.com/tofuutils/tenv/versionmanager/semantic/parser/types"
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
	conf                  *config.Config
	FolderName            string
	predicateReaders      []types.PredicateReader
	retriever             ReleaseInfoRetriever
	VersionEnvName        string
	defaultVersionEnvName string
	VersionFiles          []types.VersionFile
}

func Make(conf *config.Config, folderName string, predicateReaders []types.PredicateReader, retriever ReleaseInfoRetriever, versionEnvName string, defaultVersionEnvName string, versionFiles []types.VersionFile) VersionManager {
	return VersionManager{conf: conf, FolderName: folderName, predicateReaders: predicateReaders, retriever: retriever, VersionEnvName: versionEnvName, defaultVersionEnvName: defaultVersionEnvName, VersionFiles: versionFiles}
}

// Detect version (resolve and evaluate, can install depending on auto install env var).
func (m VersionManager) Detect(proxyCall bool) (string, error) {
	configVersion, err := m.Resolve(semantic.LatestAllowedKey)
	if err != nil {
		m.conf.Displayer.Flush(proxyCall)

		return "", err
	}

	return m.Evaluate(configVersion, proxyCall)
}

// Evaluate version resolution strategy or version constraint (can install depending on auto install env var).
func (m VersionManager) Evaluate(requestedVersion string, proxyCall bool) (string, error) {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		cleanedVersion := parsedVersion.String()
		if m.conf.NoInstall {
			m.conf.Displayer.Flush(proxyCall)

			return cleanedVersion, nil
		}

		return cleanedVersion, m.installSpecificVersion(cleanedVersion, proxyCall)
	}

	predicateInfo, err := semantic.ParsePredicate(requestedVersion, m.FolderName, m.predicateReaders, m.conf)
	if err != nil {
		m.conf.Displayer.Flush(proxyCall)

		return "", err
	}

	if !m.conf.ForceRemote {
		versions, err := m.ListLocal(predicateInfo.ReverseOrder)
		if err != nil {
			m.conf.Displayer.Flush(proxyCall)

			return "", err
		}

		for _, version := range versions {
			if predicateInfo.Predicate(version) {
				m.conf.Displayer.Display("Found compatible version installed locally : " + version)
				m.conf.Displayer.Flush(proxyCall)

				return version, nil
			}
		}

		if m.conf.NoInstall {
			m.conf.Displayer.Flush(proxyCall)

			return "", errNoCompatible
		}
		m.conf.Displayer.Display("No compatible version found locally, search a remote one...")
	}

	return m.searchInstallRemote(predicateInfo, m.conf.NoInstall, proxyCall)
}

func (m VersionManager) Install(requestedVersion string) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		return m.installSpecificVersion(parsedVersion.String(), false)
	}

	predicateInfo, err := semantic.ParsePredicate(requestedVersion, m.FolderName, m.predicateReaders, m.conf)
	if err != nil {
		return err
	}

	// noInstall is set to false to force install regardless of conf
	_, err = m.searchInstallRemote(predicateInfo, false, false)

	return err
}

// try to ensure the directory exists with a MkdirAll call.
// (made lazy method : not always useful and allows flag override for root path).
func (m VersionManager) InstallPath() string {
	dir := filepath.Join(m.conf.RootPath, m.FolderName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		m.conf.Displayer.Log(hclog.Warn, "Can not create installation directory", loghelper.Error, err)
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
		m.conf.Displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Can not read installed versions", loghelper.Error, err)

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
	if err == nil {
		m.conf.Displayer.Display("Removed " + versionFilePath)
	}

	return err
}

// Search the requested version in version files.
func (m VersionManager) Resolve(defaultStrategy string) (string, error) {
	version, err := m.ResolveWithoutFallback()
	if err != nil || version != "" {
		return version, err
	}

	if version = os.Getenv(m.defaultVersionEnvName); version != "" {
		return types.DisplayDetectionInfo(m.conf.Displayer, version, m.defaultVersionEnvName), nil
	}

	if version, err = flatparser.RetrieveVersion(m.RootVersionFilePath(), m.conf); err != nil || version != "" {
		return version, err
	}
	m.conf.Displayer.Display(loghelper.Concat("No version files found for ", m.FolderName, ", fallback to ", defaultStrategy, " strategy"))

	return defaultStrategy, nil
}

// Search the requested version in version files.
func (m VersionManager) ResolveWithoutFallback() (string, error) {
	if forcedVersion := os.Getenv(m.VersionEnvName); forcedVersion != "" {
		return types.DisplayDetectionInfo(m.conf.Displayer, forcedVersion, m.VersionEnvName), nil
	}

	return semantic.RetrieveVersion(m.VersionFiles, m.conf)
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
	if err = os.RemoveAll(targetPath); err == nil {
		m.conf.Displayer.Display(loghelper.Concat("Uninstallation of ", m.FolderName, " ", cleanedVersion, " successful (directory ", targetPath, " removed)"))
	}

	return err
}

func (m VersionManager) Use(requestedVersion string, workingDir bool) error {
	detectedVersion, err := m.Evaluate(requestedVersion, false)
	if err != nil {
		return err
	}

	targetFilePath := m.VersionFiles[0].Name
	if !workingDir {
		targetFilePath = m.RootVersionFilePath()
	}
	if err = os.WriteFile(targetFilePath, []byte(detectedVersion), 0644); err == nil {
		m.conf.Displayer.Display(loghelper.Concat("Written ", detectedVersion, " in ", targetFilePath))
	}

	return err
}

func (m VersionManager) installSpecificVersion(version string, proxyCall bool) error {
	if version == "" {
		m.conf.Displayer.Flush(proxyCall)

		return errEmptyVersion
	}

	installPath := m.InstallPath()
	entries, err := os.ReadDir(installPath)
	if err != nil {
		m.conf.Displayer.Flush(proxyCall)

		return err
	}

	for _, entry := range entries {
		if entry.IsDir() && version == entry.Name() {
			m.conf.Displayer.Display(loghelper.Concat(m.FolderName, " ", version, " already installed"))
			m.conf.Displayer.Flush(proxyCall)

			return nil
		}
	}

	// Always normal display when installation is need
	m.conf.Displayer.Flush(false)
	m.conf.Displayer.Display(loghelper.Concat("Installing ", m.FolderName, " ", version))

	err = m.retriever.InstallRelease(version, filepath.Join(installPath, version))
	if err == nil {
		m.conf.Displayer.Display(loghelper.Concat("Installation of ", m.FolderName, " ", version, " successful"))
	}

	return err
}

func (m VersionManager) searchInstallRemote(predicateInfo types.PredicateInfo, noInstall bool, proxyCall bool) (string, error) {
	versions, err := m.ListRemote(predicateInfo.ReverseOrder)
	if err != nil {
		m.conf.Displayer.Flush(proxyCall)

		return "", err
	}

	for _, version := range versions {
		if predicateInfo.Predicate(version) {
			m.conf.Displayer.Display("Found compatible version remotely : " + version)
			if noInstall {
				m.conf.Displayer.Flush(proxyCall)

				return version, nil
			}

			return version, m.installSpecificVersion(version, proxyCall)
		}
	}
	m.conf.Displayer.Flush(proxyCall)

	return "", errNoCompatible
}
