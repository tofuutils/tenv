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
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

const ToolFileName = ".tool-versions"

func RetrieveTofuVersion(filePath string, conf *config.Config) (string, error) {
	return retrieveVersionFromToolFile(filePath, cmdconst.OpentofuName, conf)
}

func RetrieveTfVersion(filePath string, conf *config.Config) (string, error) {
	return retrieveVersionFromToolFile(filePath, cmdconst.TerraformName, conf)
}

func RetrieveTgVersion(filePath string, conf *config.Config) (string, error) {
	return retrieveVersionFromToolFile(filePath, cmdconst.TerragruntName, conf)
}

func RetrieveAtmosVersion(filePath string, conf *config.Config) (string, error) {
	return retrieveVersionFromToolFile(filePath, cmdconst.AtmosName, conf)
}

func retrieveVersionFromToolFile(filePath, toolName string, conf *config.Config) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		conf.Displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Failed to open tool file", loghelper.Error, err)

		return "", nil
	}
	defer file.Close()

	return parseVersionFromToolFileReader(filePath, file, toolName, conf.Displayer), nil
}

func parseVersionFromToolFileReader(filePath string, reader io.Reader, toolName string, displayer loghelper.Displayer) string {
	resolvedVersion := ""
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		trimmedLine := strings.TrimSpace(scanner.Text())

		if trimmedLine == "" || trimmedLine[0] == '#' {
			continue
		}

		parts := strings.Fields(trimmedLine)
		if len(parts) >= 2 && parts[0] == toolName {
			resolvedVersion, _, _ = strings.Cut(parts[1], "#") // handle comment not separated by space
		}
	}

	if err := scanner.Err(); err != nil {
		displayer.Log(hclog.Warn, "Failed to parse tool file", loghelper.Error, err)

		return ""
	}

	if resolvedVersion == "" {
		return ""
	}

	return types.DisplayDetectionInfo(displayer, resolvedVersion, filePath)
}
