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

package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"os"

	"github.com/tofuutils/tenv/v4/pkg/fileperm"
	"github.com/tofuutils/tenv/v4/pkg/uncompress/sanitize"
)

func UnzipToDir(dataZip []byte, dirPath string, filter func(string) bool) error {
	dataReader := bytes.NewReader(dataZip)
	zipReader, err := zip.NewReader(dataReader, int64(len(dataZip)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if err = copyZipFileToDir(file, dirPath, filter); err != nil {
			return err
		}
	}

	return nil
}

// a separate function allows deferred Close to execute earlier.
func copyZipFileToDir(zipFile *zip.File, dirPath string, filter func(string) bool) error {
	destPath, err := sanitize.ArchivePath(dirPath, zipFile.Name)
	if err != nil {
		return err
	}

	if destPath[len(destPath)-1] == '/' {
		// trailing slash indicates a directory
		return os.MkdirAll(destPath, fileperm.RWE)
	}

	reader, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if !filter(destPath) {
		return nil
	}

	return os.WriteFile(destPath, data, zipFile.Mode())
}
