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

package flatparser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestRetrieveVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{name: "plain", content: "1.5.7\n", want: "1.5.7"},
		{name: "surrounding spaces", content: "  1.5.7  \n\n", want: "1.5.7"},
		{name: "windows line ending", content: "1.5.7\r\n", want: "1.5.7"},
		{name: "comment line before", content: "# pinned for the network module\n1.5.7\n", want: "1.5.7"},
		{name: "trailing comment", content: "1.5.7 # pinned\n", want: "1.5.7"},
		{name: "blank and comment lines", content: "\n# a\n\n  # b\nlatest-allowed\n", want: "latest-allowed"},
		{name: "constraint", content: ">= 1.5, < 1.7 # range\n", want: ">= 1.5, < 1.7"},
		{name: "only comments", content: "# nothing here\n", want: ""},
		{name: "empty", content: "", want: ""},
	}

	conf := config.Config{Displayer: loghelper.InertDisplayer} //nolint:exhaustruct
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			filePath := filepath.Join(t.TempDir(), ".terraform-version")
			if err := os.WriteFile(filePath, []byte(testCase.content), 0o600); err != nil {
				t.Fatal(err)
			}

			got, err := Retrieve(filePath, &conf, NoMsg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != testCase.want {
				t.Errorf("Retrieve() = %q, want %q", got, testCase.want)
			}
		})
	}
}
