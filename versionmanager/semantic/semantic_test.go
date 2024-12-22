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

package semantic_test

import (
	"slices"
	"testing"

	"github.com/tofuutils/tenv/v4/versionmanager/semantic"
)

func TestCmpVersion(t *testing.T) {
	t.Parallel()

	versions := []string{"1.6.0-beta5", "1.5.2", "1.6.0-alpha5", "1.6.0", "1.5.1", "1.5.0", "1.6.0-rc1"}
	slices.SortFunc(versions, semantic.CmpVersion)
	if !slices.Equal(versions, []string{"1.5.0", "1.5.1", "1.5.2", "1.6.0-alpha5", "1.6.0-beta5", "1.6.0-rc1", "1.6.0"}) {
		t.Error("Unmatching results, get :", versions)
	}
}

func TestStableVersion(t *testing.T) {
	t.Parallel()

	var filtered []string
	for _, version := range []string{"1.5.0", "1.5.1", "1.5.2", "1.6.0-alpha5", "1.6.0-beta5", "1.6.0-rc1", "1.6.0"} {
		if semantic.StableVersion(version) {
			filtered = append(filtered, version)
		}
	}

	if !slices.Equal(filtered, []string{"1.5.0", "1.5.1", "1.5.2", "1.6.0"}) {
		t.Error("Unmatching results, get :", filtered)
	}
}
