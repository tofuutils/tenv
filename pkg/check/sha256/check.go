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

package sha256check

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	ErrCheck = errors.New("invalid sha256 checksum")
	ErrNoSum = errors.New("file sha256 checksum not found for current platform")
)

func Check(data []byte, dataSums []byte, fileName string) error {
	dataSum, err := extract(dataSums, fileName)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(data)
	if !bytes.Equal(dataSum, hashed[:]) {
		return ErrCheck
	}

	return nil
}

func extract(dataSums []byte, fileName string) ([]byte, error) {
	dataSumsStr := string(dataSums)
	for _, dataSumStr := range strings.Split(dataSumsStr, "\n") {
		dataSumStr, ok := strings.CutSuffix(dataSumStr, fileName)
		if ok {
			dataSumStr = strings.TrimSpace(dataSumStr)

			return hex.DecodeString(dataSumStr)
		}
	}

	return nil, ErrNoSum
}
