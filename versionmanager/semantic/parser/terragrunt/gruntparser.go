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
	"github.com/tofuutils/tenv/config"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

const (
	hclName  = "terragrunt.hcl"
	jsonName = "terragrunt.hcl.json"

	terraformVersionConstraintName  = "terraform_version_constraint"
	terragruntVersionConstraintName = "terraform_version_constraint"

	msgTerraGruntErr = "Failed to read terragrunt file :"
)

var terraformVersionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: terraformVersionConstraintName}},
}

var terragruntVersionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: terragruntVersionConstraintName}},
}

func RetrieveTerraformVersionConstraint(conf *config.Config) (string, error) {
	return retrieveVersionConstraint(terraformVersionPartialSchema, terraformVersionConstraintName, conf.Verbose)
}

func RetrieveTerraguntVersionConstraint(conf *config.Config) (string, error) {
	return retrieveVersionConstraint(terragruntVersionPartialSchema, terragruntVersionConstraintName, conf.Verbose)
}

func retrieveVersionConstraint(versionPartialShema *hcl.BodySchema, versionConstraintName string, verbose bool) (string, error) {
	parser := hclparse.NewParser()

	data, err := os.ReadFile(hclName)
	if err == nil {
		parsedFile, diags := parser.ParseHCL(data, hclName)

		return extractVersionConstraint(parsedFile, diags, versionPartialShema, versionConstraintName, verbose)
	}

	if verbose {
		fmt.Println(msgTerraGruntErr, err) //nolint
	}

	data, err = os.ReadFile(jsonName)
	if err != nil {
		if verbose {
			fmt.Println(msgTerraGruntErr, err) //nolint
		}

		return "", nil
	}

	parsedFile, diags := parser.ParseJSON(data, jsonName)

	return extractVersionConstraint(parsedFile, diags, versionPartialShema, versionConstraintName, verbose)
}

func extractVersionConstraint(parsedFile *hcl.File, diags hcl.Diagnostics, versionPartialShema *hcl.BodySchema, versionConstraintName string, verbose bool) (string, error) {
	if diags.HasErrors() {
		return "", diags
	}
	if parsedFile == nil {
		return "", nil
	}

	content, _, diags := parsedFile.Body.PartialContent(versionPartialShema)
	if diags.HasErrors() {
		if verbose {
			fmt.Println("Failed to parse terragrunt file :", diags) //nolint
		}

		return "", nil
	}

	attr, exists := content.Attributes[versionConstraintName]
	if !exists {
		return "", nil
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		if verbose {
			fmt.Println("Failures parsing terragrunt attribute :", diags) //nolint
		}

		return "", nil
	}

	val, err := convert.Convert(val, cty.String)
	if err != nil {
		if verbose {
			fmt.Println("Failed to convert terragrunt attribute :", err) //nolint
		}

		return "", nil
	}

	if val.IsNull() {
		if verbose {
			fmt.Println("Empty terragrunt attribute") //nolint
		}

		return "", nil
	}

	if !val.IsWhollyKnown() {
		if verbose {
			fmt.Println("Unknown terragrunt attribute") //nolint
		}

		return "", nil
	}

	return val.AsString(), nil
}
