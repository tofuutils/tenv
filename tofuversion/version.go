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

var errNoCompatibleFound = errors.New("no compatible version found")

func Install(requestedVersion string, conf *config.Config) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		return innerInstall(parsedVersion.String(), conf)
	}

	if requestedVersion == "latest" {
		latestVersion, err := github.LatestRelease(conf)
		if err != nil {
			return err
		}
		return innerInstall(latestVersion, conf)
	}

	predicate, reverseOrder, err := parsePredicate(requestedVersion, conf)
	if err != nil {
		return err
	}

	version, err := searchRemote(predicate, reverseOrder, conf)
	if err != nil {
		return err
	}
	return innerInstall(version, conf)
}

func searchRemote(predicate func(string) bool, reverseOrder bool, conf *config.Config) (string, error) {
	versions, err := ListRemote(conf)
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
	return "", errNoCompatibleFound
}

// version should be a specific one (without starting 'v')
func innerInstall(version string, conf *config.Config) error {
	installDir := conf.InstallDir()
	entries, err := os.ReadDir(installDir)
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

	downloadUrl, err := github.DownloadUrl(version, conf)
	if err != nil {
		return err
	}

	response, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	targetDir := path.Join(installDir, version)
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}
	return zip.UnzipToDir(response.Body, targetDir)
}

func ListLocal(conf *config.Config) ([]string, error) {
	entries, err := os.ReadDir(conf.InstallDir())
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

func Uninstall(requestedVersion string, conf *config.Config) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err != nil {
		return err
	}

	cleanedVersion := parsedVersion.String()
	if conf.Verbose {
		fmt.Println("Uninstallation of OpenTofu", cleanedVersion)
	}
	return os.RemoveAll(path.Join(conf.InstallDir(), cleanedVersion))
}

func LocalSet(conf *config.Config) map[string]struct{} {
	entries, err := os.ReadDir(conf.InstallDir())
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

func Use(requestedVersion string, conf *config.Config) error {
	parsedVersion, err := version.NewVersion(requestedVersion)
	if err == nil {
		cleanedVersion := parsedVersion.String()
		err = writeVersionFile(cleanedVersion, conf)
		if err != nil || conf.NoInstall {
			return err
		}
		return innerInstall(cleanedVersion, conf)
	}

	predicate, reverseOrder, err := parsePredicate(requestedVersion, conf)
	if err != nil {
		return err
	}

	versions, err := ListLocal(conf)
	if err != nil {
		return err
	}

	versionReceiver, done := iterate.Iterate(versions, reverseOrder)
	defer done()

	for version := range versionReceiver {
		if predicate(version) {
			return writeVersionFile(version, conf)
		}
	}

	if conf.NoInstall {
		return errNoCompatibleFound
	} else if conf.Verbose {
		fmt.Println("No compatible version found locally, search a remote one...")
	}

	version := ""
	if requestedVersion == "latest" {
		latestVersion, err := github.LatestRelease(conf)
		if err != nil {
			return err
		}
		version = latestVersion
	} else {
		remoteVersion, err := searchRemote(predicate, reverseOrder, conf)
		if err != nil {
			return err
		}
		version = remoteVersion
	}

	if err = writeVersionFile(version, conf); err != nil {
		return err
	}
	return innerInstall(version, conf)

}

func writeVersionFile(version string, conf *config.Config) error {
	targetPath := conf.RootFile()
	if conf.WorkingDir {
		targetPath = config.VersionFileName
	}
	if conf.Verbose {
		fmt.Println("Write", version, "in", targetPath)
	}
	return os.WriteFile(targetPath, []byte(version), 0644)
}
