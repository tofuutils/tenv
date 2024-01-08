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

var errNoCompatibleFound = errors.New("no constraint compatible version found")

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

	versions, err := ListRemote(conf)
	if err != nil {
		return err
	}

	if reverseOrder {
		// reverse order, start with latest
		for i := len(versions) - 1; i >= 0; i-- {
			version := versions[i]
			if predicate(version) {
				return innerInstall(version, conf)
			}
		}
	} else {
		// start with oldest
		for _, version := range versions {
			if predicate(version) {
				return innerInstall(version, conf)
			}
		}
	}
	return errNoCompatibleFound
}

// requestedVersion should be a specific one
func innerInstall(requestedVersion string, conf *config.Config) error {
	// TODO
	return nil
}

func List(conf *config.Config) ([]Version, error) {
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
		return innerUse(requestedVersion, conf)
	}

	predicate, reverseOrder, err := parsePredicate(requestedVersion)
	if err != nil {
		return err
	}

	versions, err := List(conf)
	if err != nil {
		return err
	}

	if reverseOrder {
		// reverse order, start with latest
		for i := len(versions) - 1; i >= 0; i-- {
			version := versions[i]
			if predicate(version.Name) {
				return innerUse(version.Name, conf)
			}
		}
	} else {
		// start with oldest
		for _, version := range versions {
			if predicate(version.Name) {
				return innerUse(version.Name, conf)
			}
		}
	}
	return errNoCompatibleFound
}

// requestedVersion should be a specific one
func innerUse(requestedVersion string, conf *config.Config) error {
	err := os.WriteFile(conf.RootFile(), []byte(requestedVersion), 0644)
	if err != nil || !conf.AutoInstall {
		return err
	}
	return innerInstall(requestedVersion, conf)
}
