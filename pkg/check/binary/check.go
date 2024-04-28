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

package bincheck

import (
	"os"
)

func Check(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	const maxBytes = 8000 // Read up to 8000 bytes for analysis
	bytes := make([]byte, maxBytes)
	numberOfBytes, err := file.Read(bytes)
	if err != nil {
		return false, err
	}

	// Check for the presence of a null byte within the read bytes
	for i := 0; i < numberOfBytes; i++ {
		if bytes[i] == 0 {
			return true, nil // Null byte found, file is binary
		}
	}

	// If no null byte found, check for a UTF-8 encoding signature (BOM)
	if numberOfBytes >= 3 && bytes[0] == 0xEF && bytes[1] == 0xBB && bytes[2] == 0xBF {
		return false, nil // UTF-8 encoded text file
	}

	// If no null byte or UTF-8 BOM found, assume it's a text file
	return false, nil
}
