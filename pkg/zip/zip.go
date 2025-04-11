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

	"github.com/tofuutils/tenv/v4/pkg/fileperm"
)

// ensure the directory exists with a MkdirAll call.
func UnzipToDir(dataZip []byte, dirPath string, filter func(string) bool) error {
	err := os.MkdirAll(dirPath, fileperm.RWE)
	if err != nil {
		return err
	}

	dataReader := bytes.NewReader(dataZip)
	zipReader, err := zip.NewReader(dataReader, int64(len(dataZip)))
	if err != nil {
		return err
	}

	// First pass: create all directories
	for _, file := range zipReader.File {
		destPath, err := SanitizeArchivePath(dirPath, file.Name)
		if err != nil {
			return err
		}

		if destPath[len(destPath)-1] == '/' {
			// trailing slash indicates a directory
			if err := os.MkdirAll(destPath, fileperm.RWE); err != nil {
				return err
			}
		} else {
			// Create parent directory for files
			if err := os.MkdirAll(filepath.Dir(destPath), fileperm.RWE); err != nil {
				return err
			}
		}
	}

	// Second pass: extract files
	for _, file := range zipReader.File {
		if err = copyZipFileToDir(file, dirPath, filter); err != nil {
			return err
		}
	}

	return nil
}

// a separate function allows deferred Close to execute earlier.
func copyZipFileToDir(zipFile *zip.File, dirPath string, filter func(string) bool) error {
	destPath, err := SanitizeArchivePath(dirPath, zipFile.Name)
	if err != nil {
		return err
	}

	if destPath[len(destPath)-1] == '/' {
		// Directory already created in first pass
		return nil
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

// SanitizeArchivePath sanitizes archive file pathing from "G305" (file traversal).
func SanitizeArchivePath(dirPath string, fileName string) (string, error) {
	// Handle empty filename
	if fileName == "" {
		return dirPath, nil
	}

	// Check for absolute paths
	if filepath.IsAbs(fileName) {
		return "", fmt.Errorf("content filepath is tainted: %s", fileName)
	}

	// Clean the paths to handle any path traversal attempts
	cleanDirPath := filepath.Clean(dirPath)
	cleanFileName := filepath.Clean(fileName)

	// Join the paths
	destPath := filepath.Join(cleanDirPath, cleanFileName)

	// Check if the resulting path is still within the target directory
	if !strings.HasPrefix(destPath, cleanDirPath) {
		return "", fmt.Errorf("content filepath is tainted: %s", fileName)
	}

	return destPath, nil
}
