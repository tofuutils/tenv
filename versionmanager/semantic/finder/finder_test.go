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

	versionfinder "github.com/tofuutils/tenv/v4/versionmanager/semantic/finder"
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

func TestFindVersionTerragruntAlpha(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Find("alpha2025022701"); version != "alpha2025022701" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestFindVersionTerragruntAlphaDash(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Find("gruntwork-io/terragrunt/releases/tag/alpha-2025040801/"); version != "alpha-2025040801" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestFindVersionConstraint(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Find("~> v1.11.0"); version != "1.11.0" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestIsValidVersion(t *testing.T) {
	t.Parallel()

	if !versionfinder.IsValid("1.11.0") {
		t.Error("Unexpected result : should be valid")
	}
}

func TestIsValidVersionV(t *testing.T) {
	t.Parallel()

	if !versionfinder.IsValid("v2.13.0") {
		t.Error("Unexpected result : should be valid")
	}
}

func TestIsValidVersionNoFixes(t *testing.T) {
	t.Parallel()

	if !versionfinder.IsValid("3.12") {
		t.Error("Unexpected result : should be valid")
	}
}

func TestIsValidVersionVNoMinor(t *testing.T) {
	t.Parallel()

	if !versionfinder.IsValid("v4") {
		t.Error("Unexpected result : should be valid")
	}
}

func TestIsValidVersionConstraint(t *testing.T) {
	t.Parallel()

	if versionfinder.IsValid("~> v1.11.0") {
		t.Error("Unexpected result : should not be valid")
	}
}

func TestIsValidIP(t *testing.T) {
	t.Parallel()

	if versionfinder.IsValid("1.2.3.4") {
		t.Error("Unexpected result : should not be valid")
	}
}

func TestIsValidVersionTerragruntAlpha(t *testing.T) {
	t.Parallel()

	if !versionfinder.IsValid("alpha2025022701") {
		t.Error("Unexpected result : should be valid")
	}
}

func TestIsValidTerragruntAlphaDash(t *testing.T) {
	t.Parallel()

	if !versionfinder.IsValid("alpha-2025040801") {
		t.Error("Unexpected result : should be valid")
	}
}

func TestCleanVersion(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Clean("1.11.0"); version != "1.11.0" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestCleanVersionV(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Clean("v2.13.0"); version != "2.13.0" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestCleanVersionNoFixes(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Clean("3.12"); version != "3.12.0" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestCleanVersionVNoMinor(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Clean("v4"); version != "4.0.0" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestCleanVersionTerragruntAlpha(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Clean("alpha2025022701"); version != "alpha2025022701" {
		t.Error("Unexpected result, get :", version)
	}
}

func TestCleanTerragruntAlphaDash(t *testing.T) {
	t.Parallel()

	if version := versionfinder.Clean("alpha-2025040801"); version != "alpha-2025040801" {
		t.Error("Unexpected result, get :", version)
	}
}
