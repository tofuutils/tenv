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

package archname

import "testing"

func TestConvert(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"amd64", "x86_64"},
		{"386", "i386"},
		{"arm64", "arm64"},
		{"arm", "arm"},
		{"", ""},
		{"x86_64", "x86_64"},
		{"i386", "i386"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Convert(tt.input); got != tt.want {
				t.Errorf("Convert(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
