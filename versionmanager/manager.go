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
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/lockfile"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/pkg/reversecmp"
	"github.com/tofuutils/tenv/v4/versionmanager/lastuse"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic"
	flatparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/flat"
	iacparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/iac"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

const (
	rwePerm = 0o755
	rwPerm  = 0o600
)

var (
	errEmptyVersion        = errors.New("empty version")
	errNoCompatible        = errors.New("no compatible version found")
	ErrNoCompatibleLocally = errors.New("no compatible version found locally")
)

type ReleaseRetriever interface {
	Install(ctx context.Context, version string, targetPath string) error
	ListVersions(ctx context.Context) ([]string, error)
}

type DatedVersion struct {
	UseDate time.Time
	Version string
}

type VersionManager struct {
	Conf         *config.Config
	EnvNames     EnvPrefix
	FolderName   string
	iacExts      []iacparser.ExtDescription
	retriever    ReleaseRetriever
	VersionFiles []types.VersionFile
}

func Make(conf *config.Config, envPrefix string, folderName string, iacExts []iacparser.ExtDescription, retriever ReleaseRetriever, versionFiles []types.VersionFile) VersionManager {
	return VersionManager{Conf: conf, EnvNames: EnvPrefix(envPrefix), FolderName: folderName, iacExts: iacExts, retriever: retriever, VersionFiles: versionFiles}
}

// Detect version (resolve and evaluate, can install depending on auto install env var).
func (m VersionManager) Detect(ctx context.Context, proxyCall bool) (string, error) {
	configVersion, err := m.Resolve(semantic.LatestAllowedKey)
	if err != nil {
		m.Conf.Displayer.Flush(proxyCall)

		return "", err
	}

	return m.Evaluate(ctx, configVersion, proxyCall)
}

// Evaluate version resolution strategy or version constraint (can install depending on auto install env var).
func (m VersionManager) Evaluate(ctx context.Context, requestedVersion string, proxyCall bool) (string, error) {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		cleanedVersion := parsedVersion.String() // use a parsable version
		if m.Conf.SkipInstall {
			_, installed, err := m.checkVersionInstallation("", cleanedVersion)
			if err != nil {
				return "", err
			}

			if !installed {
				return cleanedVersion, m.autoInstallDisabledMsg(cleanedVersion)
			}
			m.Conf.Displayer.Flush(proxyCall)

			return cleanedVersion, nil
		}

		return cleanedVersion, m.installSpecificVersion(ctx, cleanedVersion, proxyCall)
	}

	predicateInfo, err := semantic.ParsePredicate(requestedVersion, m.FolderName, m, m.iacExts, m.Conf)
	if err != nil {
		m.Conf.Displayer.Flush(proxyCall)

		return "", err
	}

	installPath, err := m.InstallPath()
	if err != nil {
		m.Conf.Displayer.Flush(proxyCall)

		return "", err
	}

	if !m.Conf.ForceRemote {
		versions, err := m.innerListLocal(installPath, predicateInfo.ReverseOrder)
		if err != nil {
			m.Conf.Displayer.Flush(proxyCall)

			return "", err
		}

		for _, version := range versions {
			if predicateInfo.Predicate(version) {
				m.Conf.Displayer.Display("Found compatible version installed locally : " + version)
				m.Conf.Displayer.Flush(proxyCall)

				return version, nil
			}
		}

		m.Conf.Displayer.Display("No compatible version found locally, search a remote one...")
	}

	return m.searchInstallRemote(ctx, predicateInfo, m.Conf.SkipInstall, proxyCall)
}

func (m VersionManager) Install(ctx context.Context, requestedVersion string) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		return m.installSpecificVersion(ctx, parsedVersion.String(), false) // use a parsable version
	}

	predicateInfo, err := semantic.ParsePredicate(requestedVersion, m.FolderName, m, m.iacExts, m.Conf)
	if err != nil {
		return err
	}

	// noInstall is set to false to force install regardless of conf
	_, err = m.searchInstallRemote(ctx, predicateInfo, false, false)

	return err
}

func (m VersionManager) InstallMultiple(ctx context.Context, versions []string) error {
	installPath, err := m.InstallPath()
	if err != nil {
		return err
	}

	deleteLock := lockfile.Write(installPath, m.Conf.Displayer)
	disableExit := lockfile.CleanAndExitOnInterrupt(deleteLock)
	defer disableExit()
	defer deleteLock()

	for _, version := range versions {
		if err = m.installSpecificVersionWithoutLock(ctx, installPath, version, false); err != nil {
			return err
		}
	}

	return nil
}

// try to ensure the directory exists with a MkdirAll call.
// (made lazy method : not always useful and allows flag override for root path).
func (m VersionManager) InstallPath() (string, error) {
	dirPath := filepath.Join(m.Conf.RootPath, m.FolderName)

	return dirPath, os.MkdirAll(dirPath, rwePerm)
}

func (m VersionManager) ListLocal(reverseOrder bool) ([]DatedVersion, error) {
	installPath, err := m.InstallPath()
	if err != nil {
		return nil, err
	}

	versions, err := m.innerListLocal(installPath, reverseOrder)
	if err != nil {
		return nil, err
	}

	datedVersions := make([]DatedVersion, 0, len(versions))
	for _, version := range versions {
		datedVersions = append(datedVersions, DatedVersion{
			UseDate: lastuse.Read(filepath.Join(installPath, version), m.Conf.Displayer),
			Version: version,
		})
	}

	return datedVersions, nil
}

func (m VersionManager) ListRemote(ctx context.Context, reverseOrder bool) ([]string, error) {
	versions, err := m.retriever.ListVersions(ctx)
	if err != nil {
		return nil, err
	}

	cmpFunc := reversecmp.Reverser[string](semantic.CmpVersion, reverseOrder)
	slices.SortFunc(versions, cmpFunc)

	return versions, nil
}

func (m VersionManager) LocalSet() map[string]struct{} {
	installPath, err := m.InstallPath()
	if err != nil {
		m.Conf.Displayer.Log(hclog.Warn, "Can not create installation directory", loghelper.Error, err)

		return nil
	}

	entries, err := os.ReadDir(installPath)
	if err != nil {
		m.Conf.Displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Can not read installed versions", loghelper.Error, err)

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

func (m VersionManager) ReadDefaultConstraint() string {
	if constraint := m.Conf.Getenv(m.EnvNames.constraint()); constraint != "" {
		return constraint
	}

	constraint, _ := flatparser.Retrieve(m.RootConstraintFilePath(), m.Conf, flatparser.NoMsg)

	return constraint
}

func (m VersionManager) ResetConstraint() error {
	return removeFile(m.RootConstraintFilePath(), m.Conf)
}

func (m VersionManager) ResetVersion() error {
	return removeFile(m.RootVersionFilePath(), m.Conf)
}

// Search the requested version in version files (with fallbacks and env var overloading).
func (m VersionManager) Resolve(defaultStrategy string) (string, error) {
	versionEnvName := m.EnvNames.Version()
	version := m.Conf.Getenv(versionEnvName)
	if version != "" {
		return types.DisplayDetectionInfo(m.Conf.Displayer, version, versionEnvName), nil
	}

	version, err := m.ResolveWithVersionFiles()
	if err != nil || version != "" {
		return version, err
	}

	defaultVersionEnvName := m.EnvNames.defaultVersion()
	if version = m.Conf.Getenv(defaultVersionEnvName); version != "" {
		return types.DisplayDetectionInfo(m.Conf.Displayer, version, defaultVersionEnvName), nil
	}

	if version, err = flatparser.RetrieveVersion(m.RootVersionFilePath(), m.Conf); err != nil || version != "" {
		return version, err
	}
	m.Conf.Displayer.Display(loghelper.Concat("No version files found for ", m.FolderName, ", fallback to ", defaultStrategy, " strategy"))

	return defaultStrategy, nil
}

// Search the requested version in version files.
func (m VersionManager) ResolveWithVersionFiles() (string, error) {
	return semantic.RetrieveVersion(m.VersionFiles, m.Conf)
}

// (made lazy method : not always useful and allows flag override for root path).
func (m VersionManager) RootConstraintFilePath() string {
	return filepath.Join(m.Conf.RootPath, m.FolderName, "constraint")
}

// (made lazy method : not always useful and allows flag override for root path).
func (m VersionManager) RootVersionFilePath() string {
	return filepath.Join(m.Conf.RootPath, m.FolderName, "version")
}

func (m VersionManager) SetConstraint(constraint string) error {
	_, err := version.NewConstraint(constraint) // check the use of a parsable constraint
	if err != nil {
		return err
	}

	return writeFile(m.RootConstraintFilePath(), constraint, m.Conf)
}

func (m VersionManager) Uninstall(requestedVersion string) error {
	installPath, err := m.InstallPath()
	if err != nil {
		return err
	}

	deleteLock := lockfile.Write(installPath, m.Conf.Displayer)
	disableExit := lockfile.CleanAndExitOnInterrupt(deleteLock)
	defer disableExit()
	defer deleteLock()

	parsedVersion, err := version.NewVersion(requestedVersion) // check the use of a parsable version
	if err == nil {
		m.uninstallSpecificVersion(installPath, parsedVersion.String())

		return nil
	}

	versions, err := m.innerListLocal(installPath, true)
	if err != nil {
		return err
	}

	selected, err := semantic.SelectVersionsToUninstall(requestedVersion, installPath, versions, m.Conf.Displayer)
	if err != nil {
		return err
	}

	if len(selected) == 0 {
		m.Conf.Displayer.Display(loghelper.Concat("No matching ", m.FolderName, " versions"))

		return nil
	}

	m.Conf.Displayer.Display(loghelper.Concat("Selected ", m.FolderName, " versions for uninstallation :"))
	m.Conf.Displayer.Display(strings.Join(selected, ", "))
	m.Conf.Displayer.Display("Uninstall ? [y/N]")

	buffer := make([]byte, 1)
	if _, err = os.Stdin.Read(buffer); err != nil {
		return err
	}

	read := buffer[0]

	if doUninstall := read == 'y' || read == 'Y'; !doUninstall {
		return nil
	}

	for _, version := range selected {
		m.uninstallSpecificVersion(installPath, version)
	}

	return nil
}

func (m VersionManager) UninstallMultiple(versions []string) error {
	installPath, err := m.InstallPath()
	if err != nil {
		return err
	}

	deleteLock := lockfile.Write(installPath, m.Conf.Displayer)
	disableExit := lockfile.CleanAndExitOnInterrupt(deleteLock)
	defer disableExit()
	defer deleteLock()

	for _, version := range versions {
		m.uninstallSpecificVersion(installPath, version)
	}

	return nil
}

func (m VersionManager) Use(ctx context.Context, requestedVersion string, workingDir bool) error {
	detectedVersion, err := m.Evaluate(ctx, requestedVersion, false)
	if err != nil {
		if !errors.Is(err, ErrNoCompatibleLocally) {
			return err
		}

		m.Conf.Displayer.Display(err.Error())
	}

	targetFilePath := m.VersionFiles[0].Name
	if !workingDir {
		targetFilePath = m.RootVersionFilePath()
	}

	return writeFile(targetFilePath, detectedVersion, m.Conf)
}

func (m VersionManager) alreadyInstalledMsg(version string, proxyCall bool) {
	m.Conf.Displayer.Display(loghelper.Concat(m.FolderName, " ", version, " already installed"))
	m.Conf.Displayer.Flush(proxyCall)
}

func (m VersionManager) autoInstallDisabledMsg(version string) error {
	cmdName := strings.ToLower(m.FolderName)
	m.Conf.Displayer.Flush(false) // Always normal display when installation is missing
	m.Conf.Displayer.Display(loghelper.Concat("Auto-install is disabled. To install ", m.FolderName, " version ", version, ", you can set environment variable TENV_AUTO_INSTALL=true, or install it via any of the following command: 'tenv ", cmdName, " install', 'tenv ", cmdName, " install ", version, "'"))

	return ErrNoCompatibleLocally
}

func (m VersionManager) checkVersionInstallation(installPath string, version string) (string, bool, error) {
	var err error
	if installPath == "" {
		installPath, err = m.InstallPath()
		if err != nil {
			return "", false, err
		}
	}

	if _, err = os.Stat(filepath.Join(installPath, version)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return installPath, false, nil
		}

		return "", false, err
	}

	return installPath, true, nil
}

func (m VersionManager) innerListLocal(installPath string, reverseOrder bool) ([]string, error) {
	entries, err := os.ReadDir(installPath)
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

func (m VersionManager) installSpecificVersion(ctx context.Context, version string, proxyCall bool) error {
	if version == "" {
		m.Conf.Displayer.Flush(proxyCall)

		return errEmptyVersion
	}

	// first check without lock
	installPath, installed, err := m.checkVersionInstallation("", version)
	if err != nil {
		return err
	}

	if installed {
		m.alreadyInstalledMsg(version, proxyCall)

		return nil
	}

	deleteLock := lockfile.Write(installPath, m.Conf.Displayer)
	disableExit := lockfile.CleanAndExitOnInterrupt(deleteLock)
	defer disableExit()
	defer deleteLock()

	return m.installSpecificVersionWithoutLock(ctx, installPath, version, proxyCall)
}

func (m VersionManager) installSpecificVersionWithoutLock(ctx context.Context, installPath string, version string, proxyCall bool) error {
	// second check with lock to ensure there is no ongoing install
	_, installed, err := m.checkVersionInstallation(installPath, version)
	if err != nil {
		return err
	}

	if installed {
		m.alreadyInstalledMsg(version, proxyCall)

		return nil
	}

	// Always normal display when installation is needed
	m.Conf.Displayer.Flush(false)
	m.Conf.Displayer.Display(loghelper.Concat("Installing ", m.FolderName, " ", version))

	err = m.retriever.Install(ctx, version, filepath.Join(installPath, version))
	if err == nil {
		m.Conf.Displayer.Display(loghelper.Concat("Installation of ", m.FolderName, " ", version, " successful"))
	}

	return err
}

func (m VersionManager) searchInstallRemote(ctx context.Context, predicateInfo types.PredicateInfo, noInstall bool, proxyCall bool) (string, error) {
	versions, err := m.ListRemote(ctx, predicateInfo.ReverseOrder)
	if err != nil {
		m.Conf.Displayer.Flush(proxyCall)

		return "", err
	}

	for _, version := range versions {
		if predicateInfo.Predicate(version) {
			m.Conf.Displayer.Display("Found compatible version remotely : " + version)
			if noInstall {
				return version, m.autoInstallDisabledMsg(version)
			}

			return version, m.installSpecificVersion(ctx, version, proxyCall)
		}
	}
	m.Conf.Displayer.Flush(proxyCall)

	return "", errNoCompatible
}

func (m VersionManager) uninstallSpecificVersion(installPath string, version string) {
	if version == "" {
		m.Conf.Displayer.Display(errEmptyVersion.Error())

		return
	}

	targetPath := filepath.Join(installPath, version)
	err := os.RemoveAll(targetPath)
	if err == nil {
		m.Conf.Displayer.Display(loghelper.Concat("Uninstallation of ", m.FolderName, " ", version, " successful (directory ", targetPath, " removed)"))
	} else {
		m.Conf.Displayer.Display(loghelper.Concat("Uninstallation of ", m.FolderName, " ", version, " failed with error : ", err.Error()))
	}
}

func removeFile(filePath string, conf *config.Config) error {
	err := os.RemoveAll(filePath)
	if err == nil {
		conf.Displayer.Display("Removed " + filePath)
	}

	return err
}

func writeFile(filePath string, content string, conf *config.Config) error {
	err := os.WriteFile(filePath, []byte(content), rwPerm)
	if err == nil {
		conf.Displayer.Display(loghelper.Concat("Written ", content, " in ", filePath))
	}

	return err
}
