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

package lastuse

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/tofuutils/tenv/v3/pkg/loghelper"
)

const fileName = "last-use.txt"

func Read(dirPath string, displayer loghelper.Displayer) time.Time {
	data, err := os.ReadFile(filepath.Join(dirPath, fileName))
	if err != nil {
		displayer.Log(loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), "Unable to read date in file", loghelper.Error, err)

		return time.Time{}
	}

	parsed, err := time.Parse(time.DateOnly, string(data)) //nolint
	if err != nil {
		displayer.Log(hclog.Warn, "Unable to parse date in file", loghelper.Error, err)

		return time.Time{}
	}

	return parsed
}

func WriteNow(dirPath string, displayer loghelper.Displayer) {
	lastUsePath := filepath.Join(dirPath, fileName)
	nowData := time.Now().AppendFormat(nil, time.DateOnly) //nolint

	if err := os.WriteFile(lastUsePath, nowData, 0o644); err != nil {
		displayer.Log(hclog.Warn, "Unable to write date in file", loghelper.Error, err)
	}
}
