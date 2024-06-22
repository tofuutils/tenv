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

import "strings"

// Handle '/' and '\' separated path.
func NameEqual(targetName string) func(string) bool {
	return func(path string) bool {
		separatorIndex := strings.LastIndexByte(path, '/')
		if separatorIndex == -1 {
			separatorIndex = strings.LastIndexByte(path, '\\')
		}

		return path[separatorIndex+1:] == targetName
	}
}
