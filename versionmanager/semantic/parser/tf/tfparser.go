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

package tfparser

import (
	"io/fs"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

const requiredVersionName = "required_version"

type extDescription struct {
	value    string
	len      int
	parseHCL bool
}

var exts = []extDescription{{value: ".tf", parseHCL: true}, {value: ".tf.json", parseHCL: false}} //nolint

var terraformPartialSchema = &hcl.BodySchema{ //nolint
	Blocks: []hcl.BlockHeaderSchema{{Type: "terraform"}},
}

var versionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: requiredVersionName}},
}

func init() {
	for i, desc := range exts {
		desc.len = len(desc.value)
		exts[i] = desc // override with updated copy
	}
}

func GatherRequiredVersion(conf *config.Config) ([]string, error) {
	conf.Display("Scan project to find .tf files")

	var requireds []string
	var foundFiles []string
	parser := hclparse.NewParser()
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		pathLen := len(path)
		var parsedFile *hcl.File
		var diags hcl.Diagnostics
		for _, extDesc := range exts {
			if start := pathLen - extDesc.len; start < 0 || path[start:] != extDesc.value {
				continue
			}

			foundFiles = append(foundFiles, path)

			if extDesc.parseHCL {
				parsedFile, diags = parser.ParseHCLFile(path)
			} else {
				parsedFile, diags = parser.ParseJSONFile(path)
			}
			if diags.HasErrors() {
				return diags
			}
			if parsedFile == nil {
				return nil
			}

			extracted := extractRequiredVersion(parsedFile.Body, conf)
			requireds = append(requireds, extracted...)

			return nil
		}

		return nil
	})

	if conf.AppLogger.IsDebug() {
		if len(foundFiles) == 0 {
			conf.AppLogger.Debug("No .tf file found")
		} else {
			conf.AppLogger.Debug("Read", "filePaths", foundFiles)
		}
	}

	return requireds, err
}

func extractRequiredVersion(body hcl.Body, conf *config.Config) []string {
	rootContent, _, diags := body.PartialContent(terraformPartialSchema)
	if diags.HasErrors() {
		conf.AppLogger.Warn("Failed to parse tf file", loghelper.Error, diags)

		return nil
	}

	requireds := make([]string, 0, 1)
	for _, block := range rootContent.Blocks {
		content, _, diags := block.Body.PartialContent(versionPartialSchema)
		if diags.HasErrors() {
			conf.AppLogger.Warn("Failed to parse tf block", loghelper.Error, diags)

			return nil
		}

		attr, exists := content.Attributes[requiredVersionName]
		if !exists {
			continue
		}

		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			conf.AppLogger.Warn("Failures parsing tf attribute", loghelper.Error, diags)

			return nil
		}

		val, err := convert.Convert(val, cty.String)
		if err != nil {
			conf.AppLogger.Warn("Failed to convert tf attribute", loghelper.Error, err)

			return nil
		}

		if val.IsNull() {
			conf.AppLogger.Debug("Empty tf attribute")

			continue
		}

		if !val.IsWhollyKnown() {
			conf.AppLogger.Warn("Unknown tf attribute")

			continue
		}
		requireds = append(requireds, val.AsString())
	}

	return requireds
}
