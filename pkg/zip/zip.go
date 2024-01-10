/*
 *
 * Copyright 2024 gotofuenv authors.
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
	"path"
)

func UnzipToDir(zipBodyReader io.Reader, dirPath string) error {
	zipBody, err := io.ReadAll(zipBodyReader)
	if err != nil {
		return err
	}

	byteReader := bytes.NewReader(zipBody)
	zipReader, err := zip.NewReader(byteReader, int64(len(zipBody)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if err = copyZipFileToDir(file, dirPath); err != nil {
			return err
		}
	}
	return nil
}

// a separate function allows deferred Close to execute earlier
func copyZipFileToDir(zipFile *zip.File, dirPath string) error {
	destPath := path.Join(dirPath, zipFile.Name)
	if destPath[len(destPath)-1] == '/' {
		// trailing slash indicates a directory
		return os.MkdirAll(destPath, 0755)
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
	return os.WriteFile(destPath, data, 0644)
}
