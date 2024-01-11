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

package tofuversion

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"slices"

	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/pkg/iterate"
	"github.com/dvaumoron/gotofuenv/pkg/zip"
	"github.com/dvaumoron/gotofuenv/tofuversion/github"
	"github.com/hashicorp/go-version"
)

var errNoCompatible = errors.New("no compatible version found")

func Detect(requestedVersion string, conf *config.Config) (string, error) {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		cleanedVersion := parsedVersion.String()
		if conf.NoInstall {
			return cleanedVersion, nil
		}
		return cleanedVersion, installSpecificVersion(cleanedVersion, conf)
	}

	predicate, reverseOrder, err := parsePredicate(requestedVersion, conf)
	if err != nil {
		return "", err
	}

	versions, err := ListLocal(conf)
	if err != nil {
		return "", err
	}

	versionReceiver, done := iterate.Iterate(versions, reverseOrder)
	defer done()

	for version := range versionReceiver {
		if predicate(version) {
			return version, nil
		}
	}

	if conf.NoInstall {
		return "", errNoCompatible
	} else if conf.Verbose {
		fmt.Println("No compatible version found locally, search a remote one...")
	}

	if requestedVersion == config.LatestKey {
		return installLatest(conf)
	}
	return searchInstallRemote(predicate, reverseOrder, conf)
}

func Install(requestedVersion string, conf *config.Config) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		return installSpecificVersion(parsedVersion.String(), conf)
	}

	if requestedVersion == config.LatestKey {
		_, err = installLatest(conf)
		return err
	}

	predicate, reverseOrder, err := parsePredicate(requestedVersion, conf)
	if err != nil {
		return err
	}
	_, err = searchInstallRemote(predicate, reverseOrder, conf)
	return err
}

func ListLocal(conf *config.Config) ([]string, error) {
	entries, err := os.ReadDir(conf.InstallPath())
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	slices.SortFunc(versions, cmpVersion)
	return versions, nil
}

func ListRemote(conf *config.Config) ([]string, error) {
	versions, err := github.ListReleases(conf)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(versions, cmpVersion)
	return versions, nil
}

func LocalSet(conf *config.Config) map[string]struct{} {
	entries, err := os.ReadDir(conf.InstallPath())
	if err != nil {
		if conf.Verbose {
			fmt.Println("Can not read installed versions :", err)
		}
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

func Reset(conf *config.Config) error {
	return os.Remove(conf.RootVersionFilePath())
}

func Uninstall(requestedVersion string, conf *config.Config) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err != nil {
		return err
	}

	cleanedVersion := parsedVersion.String()
	if conf.Verbose {
		fmt.Println("Uninstallation of OpenTofu", cleanedVersion)
	}
	return os.RemoveAll(path.Join(conf.InstallPath(), cleanedVersion))
}

func Use(requestedVersion string, conf *config.Config) error {
	detectedVersion, err := Detect(requestedVersion, conf)
	if err != nil {
		return err
	}

	targetFilePath := config.VersionFileName
	if !conf.WorkingDir {
		targetFilePath = conf.RootVersionFilePath()
	}
	if conf.Verbose {
		fmt.Println("Write", detectedVersion, "in", targetFilePath)
	}
	return os.WriteFile(targetFilePath, []byte(detectedVersion), 0644)
}

func installLatest(conf *config.Config) (string, error) {
	latestVersion, err := github.LatestRelease(conf)
	if err != nil {
		return "", err
	}
	return latestVersion, installSpecificVersion(latestVersion, conf)
}

func installSpecificVersion(version string, conf *config.Config) error {
	installPath := conf.InstallPath()
	entries, err := os.ReadDir(installPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() && version == entry.Name() {
			if conf.Verbose {
				fmt.Println("OpenTofu", version, "already installed")
			}
			return nil
		}
	}

	if conf.Verbose {
		fmt.Println("Installation of OpenTofu", version)
	}

	downloadUrl, err := github.DownloadAssetUrl(version, conf)
	if err != nil {
		return err
	}

	response, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	targetPath := path.Join(installPath, version)
	return zip.UnzipToDir(response.Body, targetPath)
}

func searchInstallRemote(predicate func(string) bool, reverseOrder bool, conf *config.Config) (string, error) {
	versions, err := ListRemote(conf)
	if err != nil {
		return "", err
	}

	versionReceiver, done := iterate.Iterate(versions, reverseOrder)
	defer done()

	for version := range versionReceiver {
		if predicate(version) {
			return version, installSpecificVersion(version, conf)
		}
	}
	return "", errNoCompatible
}
