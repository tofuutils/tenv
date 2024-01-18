/*
 *
 * Copyright 2024 gotofuenv authors.
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

package reversecmp_test

import (
	"cmp"
	"testing"

	"github.com/dvaumoron/gotofuenv/pkg/reversecmp"
)

func TestReverserFalse(t *testing.T) {
	reversed := reversecmp.Reverser(cmp.Compare[int], false)
	if reversed(0, 5) != -1 {
		t.Error("Not ordered")
	}
	if reversed(1, 1) != 0 {
		t.Error("WTF")
	}
	if reversed(10, 5) != 1 {
		t.Error("Not ordered again")
	}
}

func TestReverserTrue(t *testing.T) {
	reversed := reversecmp.Reverser(cmp.Compare[int], true)
	if reversed(0, 5) != 1 {
		t.Error("Not inversed")
	}
	if reversed(1, 1) != 0 {
		t.Error("WTF")
	}
	if reversed(10, 5) != -1 {
		t.Error("Not inversed again")
	}
}
