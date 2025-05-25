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

package limitedcopy

import (
	"errors"
	"io"
	"os"
)

const copyAllowedSize = 200 << 20 // 200MB, should be enough for our use cases.

var errFileTooBig = errors.New("file too big, max allowed size is 200MB")

func Copy(destPath string, reader io.Reader, perm os.FileMode) error {
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	defer destFile.Close()

	switch n, err := io.CopyN(destFile, reader, copyAllowedSize); {
	case err != nil:
		return FilterEOF(err)
	case n == copyAllowedSize:
		return errFileTooBig
	default:
		return nil
	}
}

func FilterEOF(err error) error {
	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}
