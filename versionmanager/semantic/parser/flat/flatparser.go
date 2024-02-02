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

package flatparser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tofuutils/tenv/config"
)

const msgFlatErr = "Failed to read file :"

func RetrieveVersion(versionFileNames []string, conf *config.Config) string {
	for _, fileName := range versionFileNames {
		data, err := os.ReadFile(fileName)
		if err == nil {
			resolvedVersion := string(bytes.TrimSpace(data))
			if conf.Verbose && resolvedVersion != "" {
				fmt.Println("Resolved version from", fileName, ":", resolvedVersion) //nolint
			}

			return resolvedVersion
		}
		if conf.Verbose {
			fmt.Println(msgFlatErr, err) //nolint
		}
	}

	checkedPath := map[string]struct{}{}
	if previousPath, err := os.Getwd(); err == nil {
		currentPath := filepath.Dir(previousPath)
		for currentPath != previousPath {
			version := retrieveVersionFromDir(versionFileNames, currentPath, conf.Verbose)
			if version != "" {
				return version
			}

			checkedPath[currentPath] = struct{}{}
			previousPath = currentPath
			currentPath = filepath.Dir(currentPath)
		}
	} else if conf.Verbose {
		fmt.Println("Failed to resolve working directory :", err) //nolint
	}

	if _, ok := checkedPath[conf.UserPath]; !ok {
		version := retrieveVersionFromDir(versionFileNames, conf.UserPath, conf.Verbose)
		if version != "" {
			return version
		}
	}

	if _, ok := checkedPath[conf.RootPath]; ok {
		return ""
	}

	return retrieveVersionFromDir(versionFileNames, conf.RootPath, conf.Verbose)
}

func retrieveVersionFromDir(versionFileNames []string, dirPath string, verbose bool) string {
	for _, fileName := range versionFileNames {
		filePath := filepath.Join(dirPath, fileName)
		data, err := os.ReadFile(filePath)
		if err == nil {
			resolvedVersion := string(bytes.TrimSpace(data))
			if verbose && resolvedVersion != "" {
				fmt.Println("Resolved version from", filePath, ":", resolvedVersion) //nolint
			}

			return resolvedVersion
		}
		if verbose {
			fmt.Println(msgFlatErr, err) //nolint
		}
	}

	return ""
}
