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
	"path/filepath"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

func RetrieveVersion(versionFiles []types.VersionFile, conf *config.Config) (string, error) {
	previousPath, err := filepath.Abs(conf.WorkPath)
	if err != nil {
		return "", err
	}

	if version, err := retrieveVersionFromDir(versionFiles, previousPath, conf); err != nil || version != "" {
		return version, err
	}

	userPathDone := false
	for currentPath := filepath.Dir(previousPath); currentPath != previousPath; previousPath, currentPath = currentPath, filepath.Dir(currentPath) {
		if version, err := retrieveVersionFromDir(versionFiles, currentPath, conf); err != nil || version != "" {
			return version, err
		}

		if currentPath == conf.UserPath {
			userPathDone = true
		}
	}

	if userPathDone {
		return "", nil
	}

	return retrieveVersionFromDir(versionFiles, conf.UserPath, conf)
}

func retrieveVersionFromDir(versionFiles []types.VersionFile, dirPath string, conf *config.Config) (string, error) {
	for _, versionFile := range versionFiles {
		if version, err := versionFile.Parser(filepath.Join(dirPath, versionFile.Name), conf); err != nil || version != "" {
			return version, err
		}
	}

	return "", nil
}
