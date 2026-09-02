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

package terraformretriever

import (
	"runtime"
	"testing"
)

func TestBuildAssetNames(t *testing.T) {
	t.Parallel()

	fileName, sumsFileName, sigFileName := buildAssetNames("1.6.0", "amd64")

	expectedSums := "terraform_1.6.0_SHA256SUMS"
	expectedFile := "terraform_1.6.0_" + runtime.GOOS + "_amd64.zip"
	expectedSig := expectedSums + ".sig"

	if fileName != expectedFile {
		t.Errorf("buildAssetNames() fileName = %q, want %q", fileName, expectedFile)
	}

	if sumsFileName != expectedSums {
		t.Errorf("buildAssetNames() sumsFileName = %q, want %q", sumsFileName, expectedSums)
	}

	if sigFileName != expectedSig {
		t.Errorf("buildAssetNames() sigFileName = %q, want %q", sigFileName, expectedSig)
	}
}

func TestBuildAssetNames_DifferentArch(t *testing.T) {
	t.Parallel()

	fileName, sumsFileName, sigFileName := buildAssetNames("1.5.7", "arm64")

	expectedSums := "terraform_1.5.7_SHA256SUMS"
	expectedFile := "terraform_1.5.7_" + runtime.GOOS + "_arm64.zip"

	if fileName != expectedFile {
		t.Errorf("buildAssetNames() fileName = %q, want %q", fileName, expectedFile)
	}

	if sumsFileName != expectedSums {
		t.Errorf("buildAssetNames() sumsFileName = %q, want %q", sumsFileName, expectedSums)
	}

	if sigFileName != expectedSums+".sig" {
		t.Errorf("buildAssetNames() sigFileName = %q, want %q", sigFileName, expectedSums+".sig")
	}
}
