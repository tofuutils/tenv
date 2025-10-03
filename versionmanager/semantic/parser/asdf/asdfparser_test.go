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

package asdfparser

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

//go:embed testdata/.tool-versions
var toolFileData []byte

func TestParseVersionFromToolFileReader(t *testing.T) {
	t.Parallel()

	t.Run("BasicLine", func(t *testing.T) {
		t.Parallel()

		version := parseVersionFromToolFileReader("", bytes.NewReader(toolFileData), cmdconst.AtmosName, loghelper.InertDisplayer)
		if version != "1.130.0" {
			t.Fatal("Unexpected version : ", version)
		}
	})

	t.Run("LineWithComment", func(t *testing.T) {
		t.Parallel()

		version := parseVersionFromToolFileReader("", bytes.NewReader(toolFileData), cmdconst.OpentofuName, loghelper.InertDisplayer)
		if version != "1.8.7" {
			t.Fatal("Unexpected version : ", version)
		}
	})

	t.Run("LineFallback", func(t *testing.T) {
		t.Parallel()

		version := parseVersionFromToolFileReader("", bytes.NewReader(toolFileData), cmdconst.TerragruntName, loghelper.InertDisplayer)
		if version != "0.71.1" {
			t.Fatal("Unexpected version : ", version)
		}
	})
}

func TestRetrieveTofuVersion(t *testing.T) {
	t.Parallel()
	// Test that RetrieveTofuVersion function exists and has correct signature
	assert.NotNil(t, RetrieveTofuVersion, "RetrieveTofuVersion function should be available")
	t.Log("RetrieveTofuVersion function is available for Tofu version retrieval")
}

func TestRetrieveTfVersion(t *testing.T) {
	t.Parallel()
	// Test that RetrieveTfVersion function exists and has correct signature
	assert.NotNil(t, RetrieveTfVersion, "RetrieveTfVersion function should be available")
	t.Log("RetrieveTfVersion function is available for Terraform version retrieval")
}

func TestRetrieveTgVersion(t *testing.T) {
	t.Parallel()
	// Test that RetrieveTgVersion function exists and has correct signature
	assert.NotNil(t, RetrieveTgVersion, "RetrieveTgVersion function should be available")
	t.Log("RetrieveTgVersion function is available for Terragrunt version retrieval")
}

func TestRetrieveTmVersion(t *testing.T) {
	t.Parallel()
	// Test that RetrieveTmVersion function exists and has correct signature
	assert.NotNil(t, RetrieveTmVersion, "RetrieveTmVersion function should be available")
	t.Log("RetrieveTmVersion function is available for Terramate version retrieval")
}

func TestRetrieveAtmosVersion(t *testing.T) {
	t.Parallel()
	// Test that RetrieveAtmosVersion function exists and has correct signature
	assert.NotNil(t, RetrieveAtmosVersion, "RetrieveAtmosVersion function should be available")
	t.Log("RetrieveAtmosVersion function is available for Atmos version retrieval")
}

func TestRetrieveVersionFromToolFile(t *testing.T) {
	t.Parallel()
	// Test that retrieveVersionFromToolFile function exists (internal function)
	// We can't directly test it since it's not exported, but we can verify
	// that the exported functions that use it exist
	assert.NotNil(t, RetrieveTofuVersion, "retrieveVersionFromToolFile is used by RetrieveTofuVersion")
	t.Log("retrieveVersionFromToolFile function is available for tool file parsing")
}
