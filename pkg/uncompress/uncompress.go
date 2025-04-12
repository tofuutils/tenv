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

package uncompress

import (
	"errors"
	"os"
	"strings"

	"github.com/tofuutils/tenv/v4/pkg/fileperm"
	"github.com/tofuutils/tenv/v4/pkg/uncompress/targz"
	"github.com/tofuutils/tenv/v4/pkg/uncompress/zip"
)

var errArchive = errors.New("unknown archive kind")

// ensure the directory exists with a MkdirAll call.
func ToDir(data []byte, filePath string, dirPath string, filter func(string) bool) error {
	err := os.MkdirAll(dirPath, fileperm.RWE)
	if err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(filePath, ".tar.gz"):
		return targz.UntarToDir(data, dirPath, filter)
	case strings.HasSuffix(filePath, ".zip"):
		return zip.UnzipToDir(data, dirPath, filter)
	}

	return errArchive
}
