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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	bincheck "github.com/tofuutils/tenv/pkg/check/binary"
)

// ensure the directory exists with a MkdirAll call.
func UnzipToDir(dataZip []byte, dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	dataReader := bytes.NewReader(dataZip)
	zipReader, err := zip.NewReader(dataReader, int64(len(dataZip)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if err = copyZipFileToDir(file, dirPath); err != nil {
			return err
		}

		files, err := os.ReadDir(dirPath)
		if err != nil {
			return err
		}

		for _, file := range files {
			filePath := filepath.Join(dirPath, file.Name())

			isBinary, err := bincheck.Check(filePath)
			if err != nil {
				return err
			}

			if !isBinary {
				if err = os.Remove(filePath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// a separate function allows deferred Close to execute earlier.
func copyZipFileToDir(zipFile *zip.File, dirPath string) error {
	destPath, err := sanitizeArchivePath(dirPath, zipFile.Name)
	if err != nil {
		return err
	}

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

	return os.WriteFile(destPath, data, zipFile.Mode())
}

// Sanitize archive file pathing from "G305" (file traversal).
func sanitizeArchivePath(dirPath string, fileName string) (string, error) {
	destPath := filepath.Join(dirPath, fileName)
	if strings.HasPrefix(destPath, filepath.Clean(dirPath)) {
		return destPath, nil
	}

	return "", fmt.Errorf("content filepath is tainted: %s", fileName)
}
