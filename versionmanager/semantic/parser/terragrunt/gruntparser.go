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

package terragruntparser

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

const versionConstraintName = "terraform_version_constraint"

type nameDescription struct {
	value    string
	parseHCL bool
}

var names = []nameDescription{{value: "terragrunt.hcl", parseHCL: true}, {value: "terragrunt.hcl.json", parseHCL: false}} //nolint

var terraGruntPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: versionConstraintName}},
}

func RetrieveVersionConstraint(verbose bool) (string, error) {
	parser := hclparse.NewParser()

	var parsedFile *hcl.File
	var diags hcl.Diagnostics
	for _, nameDesc := range names {
		data, err := os.ReadFile(nameDesc.value)
		if err != nil {
			if verbose {
				fmt.Println("Failed to read terragrunt file :", err) //nolint
			}

			continue
		}

		if nameDesc.parseHCL {
			parsedFile, diags = parser.ParseHCL(data, nameDesc.value)
		} else {
			parsedFile, diags = parser.ParseJSON(data, nameDesc.value)
		}
		if diags.HasErrors() {
			return "", diags
		}
		if parsedFile == nil {
			return "", nil
		}

		return extractVersionConstraint(parsedFile.Body, verbose), nil
	}

	return "", nil
}

func extractVersionConstraint(body hcl.Body, verbose bool) string {
	content, _, diags := body.PartialContent(terraGruntPartialSchema)
	if diags.HasErrors() {
		if verbose {
			fmt.Println("Failed to parse terragrunt file :", diags) //nolint
		}

		return ""
	}

	attr, exists := content.Attributes[versionConstraintName]
	if !exists {
		return ""
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		if verbose {
			fmt.Println("Failures parsing terragrunt attribute :", diags) //nolint
		}

		return ""
	}

	val, err := convert.Convert(val, cty.String)
	if err != nil {
		if verbose {
			fmt.Println("Failed to convert terragrunt attribute :", err) //nolint
		}

		return ""
	}

	if val.IsNull() {
		if verbose {
			fmt.Println("Empty terragrunt attribute") //nolint
		}

		return ""
	}

	if !val.IsWhollyKnown() {
		if verbose {
			fmt.Println("Unknown terragrunt attribute") //nolint
		}

		return ""
	}

	return val.AsString()
}
