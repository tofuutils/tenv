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

package tomlparser

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const versionName = "version"

func RetrieveVersion(filePath string, verbose bool) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if verbose {
			fmt.Println("Failed to read tgswitch file :", err) //nolint
		}

		return "", nil
	}

	var parsed map[string]string
	if _, err = toml.Decode(string(data), &parsed); err != nil {
		return "", err
	}

	resolvedVersion := parsed[versionName]
	if resolvedVersion == "" {
		return "", nil
	}
	if verbose {
		fmt.Println("Resolved version from", filePath, ":", resolvedVersion) //nolint
	}

	return resolvedVersion, nil
}
