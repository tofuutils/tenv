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

	"github.com/hashicorp/go-version"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
	"github.com/tofuutils/tenv/pkg/reversecmp"
	"github.com/tofuutils/tenv/versionmanager/semantic"
	"github.com/tofuutils/tenv/versionmanager/semantic/parser/types"
)

var (
	errEmptyVersion = errors.New("empty version")
	errNoCompatible = errors.New("no compatible version found")
)

type ReleaseInfoRetriever interface {
	InstallRelease(version string, targetPath string) error
	ListReleases() ([]string, []loghelper.RecordedMessage, error)
}

type VersionManager struct {
	conf             *config.Config
	FolderName       string
	predicateReaders []types.PredicateReader
	retriever        ReleaseInfoRetriever
	VersionEnvName   string
	VersionFiles     []types.VersionFile
}

func MakeVersionManager(conf *config.Config, folderName string, predicateReaders []types.PredicateReader, retriever ReleaseInfoRetriever, versionEnvName string, versionFiles []types.VersionFile) VersionManager {
	return VersionManager{conf: conf, FolderName: folderName, predicateReaders: predicateReaders, retriever: retriever, VersionEnvName: versionEnvName, VersionFiles: versionFiles}
}

// detect version (can install depending on auto install env var).
func (m VersionManager) Detect(multiDisplay func([]loghelper.RecordedMessage)) (string, error) {
	configVersion, err := m.Resolve(semantic.LatestAllowedKey)
	if err != nil {
		multiDisplay(configVersion.Messages)

		return "", err
	}

	return m.detect(configVersion, multiDisplay)
}

func (m VersionManager) Install(requestedVersion types.DetectionInfo) error {
	multiDisplay := loghelper.MultiDisplay(m.conf.AppLogger, m.conf.Display)

	parsedVersion, err := version.NewVersion(requestedVersion.Version)
	if err == nil {
		return m.installSpecificVersion(parsedVersion.String(), requestedVersion.Messages, multiDisplay)
	}

	predicateInfo, err := semantic.ParsePredicate(requestedVersion, m.FolderName, m.predicateReaders, m.conf)
	if err != nil {
		multiDisplay(predicateInfo.Messages)

		return err
	}

	// noInstall is set to false to force install regardless of conf
	_, err = m.searchInstallRemote(predicateInfo, false, multiDisplay)

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

func (m VersionManager) ListRemote(reverseOrder bool) ([]string, []loghelper.RecordedMessage, error) {
	versions, recordeds, err := m.retriever.ListReleases()
	if err != nil {
		return nil, recordeds, err
	}

	cmpFunc := reversecmp.Reverser[string](semantic.CmpVersion, reverseOrder)
	slices.SortFunc(versions, cmpFunc)

	return versions, recordeds, nil
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
	if err == nil {
		m.conf.Display("Removed " + versionFilePath)
	}

	return err
}

func (m VersionManager) Resolve(defaultStrategy string) (types.DetectionInfo, error) {
	if forcedVersion := os.Getenv(m.VersionEnvName); forcedVersion != "" {
		return types.MakeDetectionInfo(forcedVersion, m.VersionEnvName), nil
	}

	if detectionInfo, err := semantic.RetrieveVersion(m.VersionFiles, m.RootVersionFilePath(), m.conf); err != nil || detectionInfo.Version != "" {
		return detectionInfo, err
	}
	detectionMessages := []loghelper.RecordedMessage{{Message: loghelper.Concat("No version files found for ", m.FolderName, ", fallback to ", defaultStrategy, " strategy")}}

	return types.DetectionInfo{Version: defaultStrategy, Messages: detectionMessages}, nil
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
		m.conf.Display(loghelper.Concat("Uninstallation of ", m.FolderName, " ", cleanedVersion, " successful (directory ", targetPath, " removed)"))
	}

	return err
}

func (m VersionManager) Use(requestedVersion string, workingDir bool) error {
	detectedVersion, err := m.detect(types.DetectionInfo{Version: requestedVersion}, loghelper.MultiDisplay(m.conf.AppLogger, m.conf.Display))
	if err != nil {
		return err
	}

	targetFilePath := m.VersionFiles[0].Name
	if !workingDir {
		targetFilePath = m.RootVersionFilePath()
	}
	if err = os.WriteFile(targetFilePath, []byte(detectedVersion), 0644); err == nil {
		m.conf.Display(loghelper.Concat("Written ", detectedVersion, " in ", targetFilePath))
	}

	return err
}

func (m VersionManager) detect(requestedVersion types.DetectionInfo, multiDisplay func([]loghelper.RecordedMessage)) (string, error) {
	parsedVersion, err := version.NewVersion(requestedVersion.Version)
	if err == nil {
		cleanedVersion := parsedVersion.String()
		if m.conf.NoInstall {
			multiDisplay(requestedVersion.Messages)

			return cleanedVersion, nil
		}

		return cleanedVersion, m.installSpecificVersion(cleanedVersion, requestedVersion.Messages, multiDisplay)
	}

	predicateInfo, err := semantic.ParsePredicate(requestedVersion, m.FolderName, m.predicateReaders, m.conf)
	if err != nil {
		multiDisplay(predicateInfo.Messages)

		return "", err
	}

	if !m.conf.ForceRemote {
		versions, err := m.ListLocal(predicateInfo.ReverseOrder)
		if err != nil {
			multiDisplay(predicateInfo.Messages)

			return "", err
		}

		for _, version := range versions {
			if predicateInfo.Predicate(version) {
				recordeds := append(predicateInfo.Messages, loghelper.RecordedMessage{Message: "Found compatible version installed locally : " + version})
				multiDisplay(recordeds)

				return version, nil
			}
		}

		if m.conf.NoInstall {
			multiDisplay(predicateInfo.Messages)

			return "", errNoCompatible
		}
		predicateInfo.Messages = append(predicateInfo.Messages, loghelper.RecordedMessage{Message: "No compatible version found locally, search a remote one..."})
	}

	return m.searchInstallRemote(predicateInfo, m.conf.NoInstall, multiDisplay)
}

func (m VersionManager) installSpecificVersion(version string, recordeds []loghelper.RecordedMessage, multiDisplay func([]loghelper.RecordedMessage)) error {
	if version == "" {
		multiDisplay(recordeds)

		return errEmptyVersion
	}

	installPath := m.InstallPath()
	entries, err := os.ReadDir(installPath)
	if err != nil {
		multiDisplay(recordeds)

		return err
	}

	for _, entry := range entries {
		if entry.IsDir() && version == entry.Name() {
			recordeds = append(recordeds, loghelper.RecordedMessage{Message: loghelper.Concat(m.FolderName, " ", version, " already installed")})
			multiDisplay(recordeds)

			return nil
		}
	}

	// Always normal display when installation is need
	loghelper.MultiDisplay(m.conf.AppLogger, m.conf.Display)(recordeds)
	m.conf.Display(loghelper.Concat("Installing ", m.FolderName, " ", version))

	err = m.retriever.InstallRelease(version, filepath.Join(installPath, version))
	if err == nil {
		m.conf.Display(loghelper.Concat("Installation of ", m.FolderName, " ", version, " successful"))
	}

	return err
}

func (m VersionManager) searchInstallRemote(predicateInfo types.PredicateInfo, noInstall bool, multiDisplay func([]loghelper.RecordedMessage)) (string, error) {
	versions, recordeds2, err := m.ListRemote(predicateInfo.ReverseOrder)
	recordeds := append(predicateInfo.Messages, recordeds2...)
	if err != nil {
		multiDisplay(recordeds)

		return "", err
	}

	for _, version := range versions {
		if predicateInfo.Predicate(version) {
			recordeds = append(recordeds, loghelper.RecordedMessage{Message: "Found compatible version remotely : " + version})
			if noInstall {
				multiDisplay(recordeds)

				return version, nil
			}

			return version, m.installSpecificVersion(version, recordeds, multiDisplay)
		}
	}
	multiDisplay(recordeds)

	return "", errNoCompatible
}
