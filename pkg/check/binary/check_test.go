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

package bincheck_test

import (
	"testing"

	bincheck "github.com/tofuutils/tenv/v2/pkg/check/binary"
)

func TestTextFileCheck(t *testing.T) {
	t.Parallel()

	isBinary, err := bincheck.Check("testdata/test.txt")
	if err != nil {
		t.Fatalf("Error checking non-binary file: %v", err)
	}
	if isBinary {
		t.Errorf("Expected non-binary file, got binary")
	}
}

func TestBinaryFileCheck(t *testing.T) {
	t.Parallel()

	isBinary, err := bincheck.Check("testdata/test.bin")
	if err != nil {
		t.Fatalf("Error checking non-binary file: %v", err)
	}
	if !isBinary {
		t.Errorf("Expected binary file, got non-binary")
	}
}
