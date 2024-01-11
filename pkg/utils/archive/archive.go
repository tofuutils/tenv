/*
 *
 * Copyright 2024 opentofuutils authors.
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
package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractTarGz(source, destination string) error {
	// Open the source file
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Track the top-level directory name
	var topLevelDir string

	// Iterate through the tar archive and extract files
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			// End of tar archive
			break
		}

		if err != nil {
			return err
		}

		// If this is the first entry, capture the top-level directory name
		if topLevelDir == "" && header.Name != "pax_global_header" {
			topLevelDir = filepath.Base(header.Name)
		}

		// Construct the destination path for the current file without the top-level directory
		target := filepath.Join(destination, strings.TrimPrefix(header.Name, topLevelDir+"/"))

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directories
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}

		case tar.TypeReg, tar.TypeRegA:
			// Create regular files
			file, err := os.Create(target)
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy file contents
			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}

			// Set file permissions
			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}

	return nil
}
