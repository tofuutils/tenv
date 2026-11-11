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

package iacparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	iacparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/iac"
)

func TestGatherRequiredVersionHiddenFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.tf"), []byte("terraform {\n  required_version = \">= 1.5.0\"\n}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "._main.tf"), []byte("\x00\x05\x16\x07\x00\x02\x00\x00Mac OS X"), 0o600); err != nil {
		t.Fatal(err)
	}

	parser := hclparse.NewParser()
	exts := []iacparser.ExtDescription{{Value: ".tf", Parser: parser.ParseHCLFile}}
	conf := config.Config{Displayer: loghelper.InertDisplayer, WorkPath: dir}
	versions, err := iacparser.GatherRequiredVersion(&conf, exts)
	if err != nil {
		t.Fatal("Unexpected error : ", err)
	}
	if len(versions) != 1 || versions[0] != ">= 1.5.0" {
		t.Error("Unexpected versions : ", versions)
	}
}
