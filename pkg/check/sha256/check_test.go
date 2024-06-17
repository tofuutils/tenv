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

package sha256check_test

import (
	_ "embed"
	"testing"

	sha256check "github.com/tofuutils/tenv/v2/pkg/check/sha256"
)

//go:embed testdata/hello.txt
var data []byte

//go:embed testdata/hello_SHA256SUMS
var dataSums []byte

func TestSha256CheckCorrect(t *testing.T) {
	t.Parallel()

	if err := sha256check.Check(data, dataSums, "hello.txt"); err != nil {
		t.Error("Unexpected error : ", err)
	}
}

func TestSha256CheckError(t *testing.T) {
	t.Parallel()

	if err := sha256check.Check(data, dataSums, "hello2.txt"); err == nil {
		t.Error("Should fail on non corresponding file and fileName")
	} else if err != sha256check.ErrCheck {
		t.Error("Incorrect error reported, get :", err)
	}
}

func TestSha256Extract(t *testing.T) {
	t.Parallel()

	if err := sha256check.Check(data, dataSums, "any_name.txt"); err == nil {
		t.Error("Should fail on non exiting fileName")
	} else if err != sha256check.ErrNoSum {
		t.Error("Incorrect error reported, get :", err)
	}
}
