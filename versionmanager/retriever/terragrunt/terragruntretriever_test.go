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

package terragruntretriever

import (
	"runtime"
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/winbin"
)

const expectedSumsFileName = "SHA256SUMS"

func TestBuildAssetNames(t *testing.T) {
	t.Parallel()

	fileName, sumsFileName := buildAssetNames("amd64")

	expectedFile := "terragrunt_" + runtime.GOOS + "_amd64" + winbin.GetBinaryName("")
	if fileName != expectedFile {
		t.Errorf("buildAssetNames() fileName = %q, want %q", fileName, expectedFile)
	}

	if sumsFileName != expectedSumsFileName {
		t.Errorf("buildAssetNames() sumsFileName = %q, want %q", sumsFileName, expectedSumsFileName)
	}
}

func TestBuildAssetNames_DifferentArch(t *testing.T) {
	t.Parallel()

	fileName, sumsFileName := buildAssetNames("arm64")

	expectedFile := "terragrunt_" + runtime.GOOS + "_arm64" + winbin.GetBinaryName("")
	if fileName != expectedFile {
		t.Errorf("buildAssetNames() fileName = %q, want %q", fileName, expectedFile)
	}

	if sumsFileName != expectedSumsFileName {
		t.Errorf("buildAssetNames() sumsFileName = %q, want %q", sumsFileName, expectedSumsFileName)
	}
}
