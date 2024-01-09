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
	"os"
	"slices"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/dvaumoron/gotofuenv/config"
)

var errNoCompatibleFound = errors.New("no compatible version found")

type Version struct {
	Name string
	Used bool
}

func Install(requestedVersion string, conf *config.Config) error {
	_, err := semver.NewVersion(requestedVersion)
	if err == nil {
		return innerInstall(requestedVersion, conf)
	}

	predicate, reverseOrder, err := parsePredicate(requestedVersion)
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

	versionReceiver, done := iterate(versions, reverseOrder)
	defer done()

	for version := range versionReceiver {
		if predicate(version) {
			return version, nil
		}
	}
	return "", errNoCompatibleFound
}

// requestedVersion should be a specific one
func innerInstall(requestedVersion string, conf *config.Config) error {
	// TODO
	return nil
}

func ListLocal(conf *config.Config) ([]Version, error) {
	entries, err := os.ReadDir(conf.InstallDir())
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(conf.RootFile())
	if err != nil && conf.Verbose {
		fmt.Println("Can not read used version :", err)
	}
	usedVersion := strings.TrimSpace(string(data))

	versions := make([]Version, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			versions = append(versions, Version{
				Name: name,
				Used: usedVersion == name,
			})
		}
	}

	slices.SortFunc(versions, cmpVersion)
	return versions, nil
}

func ListRemote(conf *config.Config) ([]string, error) {
	// TODO
	return nil, nil
}

func Uninstall(requestedVersion string, conf *config.Config) error {
	// TODO
	return nil
}

func Use(requestedVersion string, conf *config.Config) error {
	_, err := semver.NewVersion(requestedVersion)
	if err == nil {
		err = writeVersionFile(requestedVersion, conf)
		if err != nil || conf.NoInstall {
			return err
		}
		return innerInstall(requestedVersion, conf)
	}

	predicate, reverseOrder, err := parsePredicate(requestedVersion)
	if err != nil {
		return err
	}

	versions, err := ListLocal(conf)
	if err != nil {
		return err
	}

	versionReceiver, done := iterate(versions, reverseOrder)
	defer done()

	for version := range versionReceiver {
		if predicate(version.Name) {
			return writeVersionFile(version.Name, conf)
		}
	}

	if conf.NoInstall {
		return errNoCompatibleFound
	} else if conf.Verbose {
		fmt.Println("No compatible version found locally, search a remote one...")
	}

	version, err := searchRemote(predicate, reverseOrder, conf)
	if err != nil {
		return err
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
