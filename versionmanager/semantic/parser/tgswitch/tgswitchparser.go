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

package tgswitchparser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/tofuutils/tenv/config"
)

const (
	tomlName = ".tgswitch.toml"

	versionName = "version"
)

func RetrieveTerraguntVersion(conf *config.Config) (string, error) {
	version, err := retrieveVersionFromFile(tomlName, conf.Verbose)
	if err != nil || version != "" {
		return version, err
	}

	version, err = retrieveVersionFromFile(filepath.Join(conf.UserPath, tomlName), conf.Verbose)
	if err != nil || version != "" {
		return version, err
	}

	return retrieveVersionFromFile(filepath.Join(conf.RootPath, tomlName), conf.Verbose)
}

func retrieveVersionFromFile(filePath string, verbose bool) (string, error) {
	data, err := os.ReadFile(tomlName)
	if err != nil {
		if verbose {
			fmt.Println("Failed to read tgswitch file :", err) //nolint
		}

		return "", nil
	}

	if verbose {
		fmt.Println("Readed", tomlName) //nolint
	}

	var parsed map[string]string
	if _, err = toml.Decode(string(data), &parsed); err != nil {
		return "", err
	}

	return parsed[versionName], nil
}
