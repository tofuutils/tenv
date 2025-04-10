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

package pathfilter_test

import (
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/pathfilter"
)

func TestNameEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		targetName string
		path       string
		want       bool
	}{
		{
			name:       "Unix path match",
			targetName: "file.txt",
			path:       "/path/to/file.txt",
			want:       true,
		},
		{
			name:       "Unix path no match",
			targetName: "file.txt",
			path:       "/path/to/other.txt",
			want:       false,
		},
		{
			name:       "Windows path match",
			targetName: "file.txt",
			path:       `C:\path\to\file.txt`,
			want:       true,
		},
		{
			name:       "Windows path no match",
			targetName: "file.txt",
			path:       `C:\path\to\other.txt`,
			want:       false,
		},
		{
			name:       "No directory separator match",
			targetName: "file.txt",
			path:       "file.txt",
			want:       true,
		},
		{
			name:       "No directory separator no match",
			targetName: "file.txt",
			path:       "other.txt",
			want:       false,
		},
		{
			name:       "Mixed separators match",
			targetName: "file.txt",
			path:       `path/to\file.txt`,
			want:       true,
		},
		{
			name:       "Empty target name",
			targetName: "",
			path:       "/path/to/",
			want:       true,
		},
		{
			name:       "Empty path",
			targetName: "file.txt",
			path:       "",
			want:       false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filter := pathfilter.NameEqual(tt.targetName)
			got := filter(tt.path)
			if got != tt.want {
				t.Errorf("NameEqual(%q)(%q) = %v, want %v", tt.targetName, tt.path, got, tt.want)
			}
		})
	}
}
