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

package asdfparser

import (
	"bufio"

	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/pkg/loghelper"
	"github.com/tofuutils/tenv/v3/versionmanager/semantic/types"
)

func NoMsg(_ loghelper.Displayer, value string, _ string) string {
	return value
}

func Retrieve(filePath, toolName string, conf *config.Config, displayMsg func(loghelper.Displayer, string, string) string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		conf.Displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Failed to read file", loghelper.Error, err)

		return "", nil
	}

	resolvedVersion := ""

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[0] == toolName {
			resolvedVersion = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		conf.Displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Failed to read file", loghelper.Error, err)

		return "", nil
	}

	if resolvedVersion == "" {
		return "", nil
	}

	return displayMsg(conf.Displayer, resolvedVersion, filePath), nil
}

func RetrieveTfVersion(filePath string, conf *config.Config) (string, error) {
	return Retrieve(filePath, "terraform", conf, types.DisplayDetectionInfo)
}

func RetrieveTofuVersion(filePath string, conf *config.Config) (string, error) {
	return Retrieve(filePath, "tofu", conf, types.DisplayDetectionInfo)
}

func RetrieveAtmosVersion(filePath string, conf *config.Config) (string, error) {
	return Retrieve(filePath, "atmos", conf, types.DisplayDetectionInfo)
}

func RetrieveTgVersion(filePath string, conf *config.Config) (string, error) {
	return Retrieve(filePath, "terragrunt", conf, types.DisplayDetectionInfo)
}
