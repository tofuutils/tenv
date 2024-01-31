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

package htmlquery

import (
	"bytes"
	_ "embed"
	"slices"
	"testing"
)

//go:embed testdata/artifactory.html
var artifactoryData []byte

func TestExtractAttr(t *testing.T) {
	t.Parallel()

	artifactoryReader := bytes.NewReader(artifactoryData)

	extracter := SelectionExtracter("style")
	extracted, err := extract(artifactoryReader, "address", extracter)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if len(extracted) != 1 && extracted[0] != "font-size:small;" {
		t.Error("Unmatching results, get :", extracted)
	}
}

func TestExtractText(t *testing.T) {
	t.Parallel()

	artifactoryReader := bytes.NewReader(artifactoryData)

	extracter := SelectionExtracter("#text")
	extracted, err := extract(artifactoryReader, "address", extracter)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if len(extracted) != 1 && extracted[0] != "Artifactory/7.68.20 Server" {
		t.Error("Unmatching results, get :", extracted)
	}
}

func TestExtractTexts(t *testing.T) {
	t.Parallel()

	artifactoryReader := bytes.NewReader(artifactoryData)

	extracted, err := extract(artifactoryReader, "a", selectionTextExtracter)
	if err != nil {
		t.Fatal("Unexpected extract error : ", err)
	}

	if !slices.Equal(extracted, []string{"../", "1.7.0/", "1.7.0-alpha20231025/", "1.7.0-rc2/", "1.7.1/", "index.json"}) {
		t.Error("Unmatching results, get :", extracted)
	}
}
