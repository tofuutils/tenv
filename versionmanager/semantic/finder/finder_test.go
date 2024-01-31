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

package versionfinder_test

import (
	"testing"

	versionfinder "github.com/tofuutils/tenv/versionmanager/semantic/finder"
)

func TestFindVersionAlphaSlash(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Find("1.7.0-alpha20231025/"); version != "1.7.0-alpha20231025" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestFindVersionEmpty(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Find("index.json"); version != "" {
		t.Error("Should not find a version, get :", version)
	}
}

func TestFindVersionPrefix(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Find("terraform_1.6.6"); version != "1.6.6" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestFindVersionPrefixSlash(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Find("terraform/1.7.0/"); version != "1.7.0" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestFindVersionPrefixV(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Find("terraform/v1.7.1/"); version != "1.7.1" {
		t.Error("Unexpected result, get :", version)
	}
}
