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

	"github.com/hashicorp/go-hclog"
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
	terragruntVersionConstraintName = "terragrunt_version_constraint"
)

var terraformVersionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: terraformVersionConstraintName}},
}

var terragruntVersionPartialSchema = &hcl.BodySchema{ //nolint
	Attributes: []hcl.AttributeSchema{{Name: terragruntVersionConstraintName}},
}

func RetrieveTerraformVersionConstraint(conf *config.Config) (string, []loghelper.RecordedMessage, error) {
	return retrieveVersionConstraint(terraformVersionPartialSchema, terraformVersionConstraintName, conf)
}

func RetrieveTerraguntVersionConstraint(conf *config.Config) (string, []loghelper.RecordedMessage, error) {
	return retrieveVersionConstraint(terragruntVersionPartialSchema, terragruntVersionConstraintName, conf)
}

func retrieveVersionConstraint(versionPartialShema *hcl.BodySchema, versionConstraintName string, conf *config.Config) (string, []loghelper.RecordedMessage, error) {
	parser := hclparse.NewParser()
	constraint, recorded, err := retrieveVersionConstraintFromFile(hclName, parser.ParseHCL, versionPartialShema, versionConstraintName, conf)
	if err != nil || constraint != "" {
		return constraint, recorded, err
	}

	constraint, recorded2, err := retrieveVersionConstraintFromFile(jsonName, parser.ParseJSON, versionPartialShema, versionConstraintName, conf)
	recorded = append(recorded, recorded2...)

	return constraint, recorded, err
}

func retrieveVersionConstraintFromFile(fileName string, fileParser func([]byte, string) (*hcl.File, hcl.Diagnostics), versionPartialShema *hcl.BodySchema, versionConstraintName string, conf *config.Config) (string, []loghelper.RecordedMessage, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		recordeds := []loghelper.RecordedMessage{{Level: loghelper.LevelWarnOrDebug(errors.Is(err, fs.ErrNotExist)), Message: "Failed to read terragrunt file", Args: []any{loghelper.Error, err}}}

		return "", recordeds, nil
	}

	parsedFile, diags := fileParser(data, fileName)
	if diags.HasErrors() {
		return "", nil, diags
	}

	recordeds := []loghelper.RecordedMessage{{Level: hclog.Debug, Message: "Read", Args: []any{"fileName", fileName}}}
	if parsedFile == nil {
		return "", recordeds, nil
	}

	content, _, diags := parsedFile.Body.PartialContent(versionPartialShema)
	if diags.HasErrors() {
		recordeds = append(recordeds, loghelper.RecordedMessage{Level: hclog.Warn, Message: "Failed to parse terragrunt file", Args: []any{loghelper.Error, diags}})

		return "", recordeds, nil
	}

	attr, exists := content.Attributes[versionConstraintName]
	if !exists {
		return "", recordeds, nil
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		recordeds = append(recordeds, loghelper.RecordedMessage{Level: hclog.Warn, Message: "Failed to parse terragrunt attribute", Args: []any{loghelper.Error, diags}})

		return "", recordeds, nil
	}

	val, err = convert.Convert(val, cty.String)
	if err != nil {
		recordeds = append(recordeds, loghelper.RecordedMessage{Level: hclog.Warn, Message: "Failed to convert terragrunt attribute", Args: []any{loghelper.Error, err}})

		return "", recordeds, nil
	}

	if val.IsNull() {
		recordeds = append(recordeds, loghelper.RecordedMessage{Level: hclog.Debug, Message: "Empty terragrunt attribute"})

		return "", recordeds, nil
	}

	if !val.IsWhollyKnown() {
		recordeds = append(recordeds, loghelper.RecordedMessage{Level: hclog.Warn, Message: "Unknown terragrunt attribute"})

		return "", recordeds, nil
	}

	return val.AsString(), recordeds, nil
}
