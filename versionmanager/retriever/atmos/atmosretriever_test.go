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

package atmosretriever

import (
	"runtime"
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/winbin"
)

func TestBuildAssetNames(t *testing.T) {
	t.Parallel()

	fileName, shaFileName := buildAssetNames("1.2.3", "amd64")

	expectedSha := "atmos_1.2.3_SHA256SUMS"
	expectedFile := "atmos_1.2.3_" + runtime.GOOS + "_amd64" + winbin.GetBinaryName("")

	if shaFileName != expectedSha {
		t.Errorf("buildAssetNames() shaFileName = %q, want %q", shaFileName, expectedSha)
	}

	if fileName != expectedFile {
		t.Errorf("buildAssetNames() fileName = %q, want %q", fileName, expectedFile)
	}
}

func TestBuildAssetNames_DifferentArch(t *testing.T) {
	t.Parallel()

	fileName, shaFileName := buildAssetNames("0.100.0", "arm64")

	if shaFileName != "atmos_0.100.0_SHA256SUMS" {
		t.Errorf("buildAssetNames() shaFileName = %q, want %q", shaFileName, "atmos_0.100.0_SHA256SUMS")
	}

	expectedFile := "atmos_0.100.0_" + runtime.GOOS + "_arm64" + winbin.GetBinaryName("")
	if fileName != expectedFile {
		t.Errorf("buildAssetNames() fileName = %q, want %q", fileName, expectedFile)
	}
}
