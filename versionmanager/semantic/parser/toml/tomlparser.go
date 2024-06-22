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

	"github.com/tofuutils/tenv/v2/config"
	"github.com/tofuutils/tenv/v2/pkg/loghelper"
	"github.com/tofuutils/tenv/v2/versionmanager/semantic/types"
)

const versionName = "version"

func RetrieveVersion(filePath string, conf *config.Config) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		conf.Displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Failed to read tgswitch file", loghelper.Error, err)

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

	return types.DisplayDetectionInfo(conf.Displayer, resolvedVersion, filePath), nil
}
