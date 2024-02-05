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
	"errors"
	"io/fs"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

const (
	hclName  = "terragrunt.hcl"
	jsonName = "terragrunt.hcl.json"

	terraformVersionConstraintName  = "terraform_version_constraint"
	terragruntVersionConstraintName = "terraform_version_constraint"
)

var terraformVersionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: terraformVersionConstraintName}},
}

var terragruntVersionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: terragruntVersionConstraintName}},
}

func RetrieveTerraformVersionConstraint(conf *config.Config) (string, error) {
	return retrieveVersionConstraint(terraformVersionPartialSchema, terraformVersionConstraintName, conf)
}

func RetrieveTerraguntVersionConstraint(conf *config.Config) (string, error) {
	return retrieveVersionConstraint(terragruntVersionPartialSchema, terragruntVersionConstraintName, conf)
}

func retrieveVersionConstraint(versionPartialShema *hcl.BodySchema, versionConstraintName string, conf *config.Config) (string, error) {
	parser := hclparse.NewParser()
	constraint, err := retrieveVersionConstraintFromFile(hclName, parser.ParseHCL, versionPartialShema, versionConstraintName, conf)
	if err != nil || constraint != "" {
		return constraint, err
	}

	return retrieveVersionConstraintFromFile(jsonName, parser.ParseJSON, versionPartialShema, versionConstraintName, conf)
}

func retrieveVersionConstraintFromFile(fileName string, fileParser func([]byte, string) (*hcl.File, hcl.Diagnostics), versionPartialShema *hcl.BodySchema, versionConstraintName string, conf *config.Config) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		conf.AppLogger.Log(loghelper.LevelWarnOrInfo(errors.Is(err, fs.ErrNotExist)), "Failed to read terragrunt file", loghelper.Error, err)

		return "", nil
	}

	parsedFile, diags := fileParser(data, fileName)
	if diags.HasErrors() {
		return "", diags
	}
	conf.AppLogger.Debug("Read", "fileName", fileName)
	if parsedFile == nil {
		return "", nil
	}

	content, _, diags := parsedFile.Body.PartialContent(versionPartialShema)
	if diags.HasErrors() {
		conf.AppLogger.Warn("Failed to parse terragrunt file", loghelper.Error, diags)

		return "", nil
	}

	attr, exists := content.Attributes[versionConstraintName]
	if !exists {
		return "", nil
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		conf.AppLogger.Warn("Failures parsing terragrunt attribute", loghelper.Error, diags)

		return "", nil
	}

	val, err = convert.Convert(val, cty.String)
	if err != nil {
		conf.AppLogger.Warn("Failed to convert terragrunt attribute", loghelper.Error, err)

		return "", nil
	}

	if val.IsNull() {
		conf.AppLogger.Info("Empty terragrunt attribute")

		return "", nil
	}

	if !val.IsWhollyKnown() {
		conf.AppLogger.Warn("Unknown terragrunt attribute")

		return "", nil
	}

	return val.AsString(), nil
}
