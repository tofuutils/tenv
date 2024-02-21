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
	"errors"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
	"github.com/tofuutils/tenv/versionmanager/semantic/parser/types"
)

const versionName = "version"

func RetrieveVersion(filePath string, conf *config.Config) (types.DetectionInfo, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		conf.AppLogger.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Failed to read tgswitch file", loghelper.Error, err)

		return types.DetectionInfo{}, nil
	}

	var parsed map[string]string
	if _, err = toml.Decode(string(data), &parsed); err != nil {
		return types.DetectionInfo{}, err
	}

	resolvedVersion := parsed[versionName]
	if resolvedVersion == "" {
		return types.DetectionInfo{}, nil
	}

	return types.MakeDetectionInfo(resolvedVersion, filePath), nil
}
