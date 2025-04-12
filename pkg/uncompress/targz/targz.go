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

package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/tofuutils/tenv/v4/pkg/fileperm"
	"github.com/tofuutils/tenv/v4/pkg/uncompress/sanitize"
)

func UntarToDir(dataTarGz []byte, dirPath string, filter func(string) bool) error {
	uncompressedStream, err := gzip.NewReader(bytes.NewReader(dataTarGz))
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		headerName := header.Name
		if !filter(headerName) {
			continue
		}

		destPath, err := sanitize.ArchivePath(dirPath, headerName)
		if err != nil {
			return err
		}

		switch typeflag := header.Typeflag; typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, fileperm.RWE); err != nil {
				return err
			}
		case tar.TypeReg:
			destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fileperm.RWE)
			if err != nil {
				return err
			}
			defer destFile.Close()

			if _, err := io.Copy(destFile, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown type during tar extraction : %c in %s", typeflag, headerName)
		}
	}
}
