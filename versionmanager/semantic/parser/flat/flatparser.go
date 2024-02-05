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
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/fatih/color"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
)

func RetrieveVersion(filePath string, conf *config.Config) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		conf.AppLogger.Log(loghelper.LevelWarnOrInfo(errors.Is(err, fs.ErrNotExist)), "Failed to read file", loghelper.Error, err)

		return "", nil
	}

	resolvedVersion := string(bytes.TrimSpace(data))
	if resolvedVersion == "" {
		return "", nil
	}
	if conf.DisplayNormal {
		fmt.Println("Resolved version from", filePath, ":", color.GreenString(resolvedVersion)) //nolint
	}

	return resolvedVersion, nil
}
