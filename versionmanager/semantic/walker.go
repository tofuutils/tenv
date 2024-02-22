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

package semantic

import (
	"os"
	"path/filepath"

	"github.com/tofuutils/tenv/config"
	flatparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/flat"
	"github.com/tofuutils/tenv/versionmanager/semantic/parser/types"
)

func RetrieveVersion(versionFiles []types.VersionFile, rootVersionFilePath string, conf *config.Config) (string, error) {
	for _, versionFile := range versionFiles {
		if version, err := versionFile.Parser(versionFile.Name, conf); err != nil || version != "" {
			return version, err
		}
	}

	previousPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	userPathNotDone := true
	for currentPath := filepath.Dir(previousPath); currentPath != previousPath; previousPath, currentPath = currentPath, filepath.Dir(currentPath) {
		if version, err := retrieveVersionFromDir(versionFiles, currentPath, conf); err != nil || version != "" {
			return version, err
		}

		if currentPath == conf.UserPath {
			userPathNotDone = false
		}
	}

	if userPathNotDone {
		if version, err := retrieveVersionFromDir(versionFiles, conf.UserPath, conf); err != nil || version != "" {
			return version, err
		}
	}

	return flatparser.RetrieveVersion(rootVersionFilePath, conf)
}

func retrieveVersionFromDir(versionFiles []types.VersionFile, dirPath string, conf *config.Config) (string, error) {
	for _, versionFile := range versionFiles {
		if version, err := versionFile.Parser(filepath.Join(dirPath, versionFile.Name), conf); err != nil || version != "" {
			return version, err
		}
	}

	return "", nil
}
