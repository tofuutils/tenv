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
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/dvaumoron/gotofuenv/config"
	"golang.org/x/mod/semver"
)

type Version struct {
	Name string
	Used bool
}

func Install(requestedVersion string, conf *config.Config) error {
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

	slices.SortFunc(versions, func(a Version, b Version) int {
		return semver.Compare(cleanVersion(a.Name), cleanVersion(b.Name))
	})
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
	// TODO
	return nil
}

func cleanVersion(version string) string {
	if version == "" || version[0] == 'v' {
		return version
	}
	return "v" + version
}
