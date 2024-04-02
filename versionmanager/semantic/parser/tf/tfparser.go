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
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
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
	parseHCL bool
}

var exts = []extDescription{{value: ".tf", parseHCL: true}, {value: ".tf.json", parseHCL: false}} //nolint

var terraformPartialSchema = &hcl.BodySchema{ //nolint
	Blocks: []hcl.BlockHeaderSchema{{Type: config.TerraformName}},
}

var versionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: requiredVersionName}},
}

func GatherRequiredVersion(conf *config.Config) ([]string, error) {
	conf.Displayer.Display("Scan project to find .tf files")

	var foundFiles []string
	if conf.Displayer.IsDebug() {
		defer func() {
			if len(foundFiles) == 0 {
				conf.Displayer.Log(hclog.Debug, "No .tf file found")
			} else {
				conf.Displayer.Log(hclog.Debug, "Read", "filePaths", foundFiles)
			}
		}()
	}

	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var requireds []string
	var parsedFile *hcl.File
	var diags hcl.Diagnostics
	parser := hclparse.NewParser()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		for _, extDesc := range exts {
			if !strings.HasSuffix(name, extDesc.value) {
				continue
			}

			foundFiles = append(foundFiles, name)

			if extDesc.parseHCL {
				parsedFile, diags = parser.ParseHCLFile(name)
			} else {
				parsedFile, diags = parser.ParseJSONFile(name)
			}
			if diags.HasErrors() {
				return nil, diags
			}
			if parsedFile == nil {
				continue
			}

			extracted := extractRequiredVersion(parsedFile.Body, conf)
			requireds = append(requireds, extracted...)
		}
	}

	return requireds, nil
}

func extractRequiredVersion(body hcl.Body, conf *config.Config) []string {
	rootContent, _, diags := body.PartialContent(terraformPartialSchema)
	if diags.HasErrors() {
		conf.Displayer.Log(hclog.Warn, "Failed to parse tf file", loghelper.Error, diags)

		return nil
	}

	requireds := make([]string, 0, 1)
	for _, block := range rootContent.Blocks {
		content, _, diags := block.Body.PartialContent(versionPartialSchema)
		if diags.HasErrors() {
			conf.Displayer.Log(hclog.Warn, "Failed to parse tf block", loghelper.Error, diags)

			return nil
		}

		attr, exists := content.Attributes[requiredVersionName]
		if !exists {
			continue
		}

		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			conf.Displayer.Log(hclog.Warn, "Failed to parse tf attribute", loghelper.Error, diags)

			return nil
		}

		val, err := convert.Convert(val, cty.String)
		if err != nil {
			conf.Displayer.Log(hclog.Warn, "Failed to convert tf attribute", loghelper.Error, err)

			return nil
		}

		if val.IsNull() {
			conf.Displayer.Log(hclog.Debug, "Empty tf attribute")

			continue
		}

		if !val.IsWhollyKnown() {
			conf.Displayer.Log(hclog.Warn, "Unknown tf attribute")

			continue
		}
		requireds = append(requireds, val.AsString())
	}

	return requireds
}
