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

package iacparser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

const requiredVersionName = "required_version"

type ExtDescription struct {
	Value  string
	Parser func(string) (*hcl.File, hcl.Diagnostics)
}

var terraformPartialSchema = &hcl.BodySchema{ //nolint
	Blocks: []hcl.BlockHeaderSchema{{Type: cmdconst.TerraformName}},
}

var versionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: requiredVersionName}},
}

func GatherRequiredVersion(conf *config.Config, exts []ExtDescription) ([]string, error) {
	if len(exts) == 0 {
		return nil, nil
	}

	conf.Displayer.Display("Scan project to find IAC files")

	var foundFiles []string //nolint
	if conf.Displayer.IsDebug() {
		defer func() {
			if len(foundFiles) == 0 {
				conf.Displayer.Log(hclog.Debug, "No IAC files found")
			} else {
				conf.Displayer.Log(hclog.Debug, "Read", "filePaths", foundFiles)
			}
		}()
	}

	entries, err := os.ReadDir(conf.WorkPath)
	if err != nil {
		return nil, err
	}

	similar := map[string]int{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		for index, extDesc := range exts {
			if cleanedName, found := strings.CutSuffix(name, extDesc.Value); found {
				extFlag := 1 << index
				similar[cleanedName] |= extFlag

				break
			}
		}
	}

	var requireds []string
	var parsedFile *hcl.File
	var diags hcl.Diagnostics
	foundFiles = make([]string, 0, len(similar))
	for cleanedName, fileExts := range similar {
		ext := filterExts(fileExts, exts)
		name := cleanedName + ext.Value
		foundFiles = append(foundFiles, name)

		parsedFile, diags = ext.Parser(filepath.Join(conf.WorkPath, name))
		if diags.HasErrors() {
			return nil, diags
		}
		if parsedFile == nil {
			continue
		}

		extracted := extractRequiredVersion(parsedFile.Body, conf)
		requireds = append(requireds, extracted...)
	}

	return requireds, nil
}

func extractRequiredVersion(body hcl.Body, conf *config.Config) []string {
	rootContent, _, diags := body.PartialContent(terraformPartialSchema)
	if diags.HasErrors() {
		conf.Displayer.Log(hclog.Warn, "Failed to parse hcl file", loghelper.Error, diags)

		return nil
	}

	requireds := make([]string, 0, 1)
	for _, block := range rootContent.Blocks {
		content, _, diags := block.Body.PartialContent(versionPartialSchema)
		if diags.HasErrors() {
			conf.Displayer.Log(hclog.Warn, "Failed to parse hcl block", loghelper.Error, diags)

			return nil
		}

		attr, exists := content.Attributes[requiredVersionName]
		if !exists {
			continue
		}

		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			conf.Displayer.Log(hclog.Warn, "Failed to parse hcl attribute", loghelper.Error, diags)

			return nil
		}

		val, err := convert.Convert(val, cty.String)
		if err != nil {
			conf.Displayer.Log(hclog.Warn, "Failed to convert hcl attribute", loghelper.Error, err)

			return nil
		}

		if val.IsNull() {
			conf.Displayer.Log(hclog.Debug, "Empty hcl attribute")

			continue
		}

		if !val.IsWhollyKnown() {
			conf.Displayer.Log(hclog.Warn, "Unknown hcl attribute")

			continue
		}
		requireds = append(requireds, val.AsString())
	}

	return requireds
}

func filterExts(fileExts int, exts []ExtDescription) ExtDescription {
	for index, ext := range exts { // has a meaningful order
		extFlag := 1 << index
		if fileExts&extFlag != 0 {
			return ext
		}
	}

	return ExtDescription{} // unreachable (fileExts should have at least one value from exts)
}
