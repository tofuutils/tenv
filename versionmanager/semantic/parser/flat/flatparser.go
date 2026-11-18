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
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

func NoMsg(_ loghelper.Displayer, value string, _ string) string {
	return value
}

func Retrieve(filePath string, conf *config.Config, displayMsg func(loghelper.Displayer, string, string) string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		conf.Displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Failed to read file", loghelper.Error, err)

		return "", nil
	}

	resolvedVersion := parseContent(string(data))
	if resolvedVersion == "" {
		return "", nil
	}

	return displayMsg(conf.Displayer, resolvedVersion, filePath), nil
}

// parseContent returns the first line that is not empty once comments are removed,
// like tfenv and tgenv do with their version files.
func parseContent(content string) string {
	for line := range strings.Lines(content) {
		line, _, _ = strings.Cut(line, "#")
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}

	return ""
}

func RetrieveVersion(filePath string, conf *config.Config) (string, error) {
	return Retrieve(filePath, conf, types.DisplayDetectionInfo)
}
