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

package pathfilter

import "testing"

func TestNameEqual(t *testing.T) {
	tests := []struct {
		name       string
		targetName string
		path       string
		want       bool
	}{
		{"unix path match", "foo", "/home/user/foo", true},
		{"unix path no match", "bar", "/home/user/foo", false},
		{"windows path match", "foo", `C:\Users\user\foo`, true},
		{"windows path no match", "bar", `C:\Users\user\foo`, false},
		{"relative path", "foo", "foo", true},
		{"relative path no match", "bar", "foo", false},
		{"nested unix", "foo", "/a/b/c/foo", true},
		{"nested windows", "foo", `C:\a\b\c\foo`, true},
		{"mixed separators last slash", "c\\foo", `/a\b/c\foo`, true},
		{"empty target no match", "", "foo", false},
		{"root path", "foo", "/foo", true},
		{"just filename", "foo", "foo", true},
		{"case sensitive", "Foo", "foo", false},
		{"trailing slash yields empty", "", "/home/user/foo/", true},
		{"trailing backslash yields empty", "", `C:\Users\user\foo\`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := NameEqual(tt.targetName)
			if got := fn(tt.path); got != tt.want {
				t.Errorf("NameEqual(%q)(%q) = %v, want %v", tt.targetName, tt.path, got, tt.want)
			}
		})
	}
}
