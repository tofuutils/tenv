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

package builder

import (
	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/versionmanager"
	atmosretriever "github.com/tofuutils/tenv/v4/versionmanager/retriever/atmos"
	terraformretriever "github.com/tofuutils/tenv/v4/versionmanager/retriever/terraform"
	terragruntretriever "github.com/tofuutils/tenv/v4/versionmanager/retriever/terragrunt"
	tofuretriever "github.com/tofuutils/tenv/v4/versionmanager/retriever/tofu"
	asdfparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/asdf"
	flatparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/flat"
	iacparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/iac"
	terragruntparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/terragrunt"
	tomlparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/toml"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

var Builders = map[string]Func{ //nolint
	cmdconst.TofuName:       BuildTofuManager,
	cmdconst.TerraformName:  BuildTfManager,
	cmdconst.TerragruntName: BuildTgManager,
	cmdconst.AtmosName:      BuildAtmosManager,
}

type Func = func(*config.Config, *hclparse.Parser) versionmanager.VersionManager

func BuildAtmosManager(conf *config.Config, _ *hclparse.Parser) versionmanager.VersionManager {
	atmosRetriever := atmosretriever.Make(conf)
	versionFiles := []types.VersionFile{
		{Name: ".atmos-version", Parser: flatparser.RetrieveVersion},
		{Name: asdfparser.ToolFileName, Parser: asdfparser.RetrieveAtmosVersion},
	}

	return versionmanager.Make(conf, config.AtmosPrefix, "Atmos", nil, atmosRetriever, versionFiles)
}

func BuildTfManager(conf *config.Config, hclParser *hclparse.Parser) versionmanager.VersionManager {
	tfRetriever := terraformretriever.Make(conf)
	gruntParser := terragruntparser.Make(hclParser)
	versionFiles := []types.VersionFile{
		{Name: ".terraform-version", Parser: flatparser.RetrieveVersion},
		{Name: ".tfswitchrc", Parser: flatparser.RetrieveVersion},
		{Name: asdfparser.ToolFileName, Parser: asdfparser.RetrieveTfVersion},
		{Name: terragruntparser.HCLName, Parser: gruntParser.RetrieveTerraformVersionConstraintFromHCL},
		{Name: terragruntparser.JSONName, Parser: gruntParser.RetrieveTerraformVersionConstraintFromJSON},
	}

	iacExts := []iacparser.ExtDescription{
		{Value: ".tf", Parser: hclParser.ParseHCLFile},
		{Value: ".tf.json", Parser: hclParser.ParseJSONFile},
	}

	return versionmanager.Make(conf, config.TfenvPrefix, "Terraform", iacExts, tfRetriever, versionFiles)
}

func BuildTgManager(conf *config.Config, hclParser *hclparse.Parser) versionmanager.VersionManager {
	tgRetriever := terragruntretriever.Make(conf)
	gruntParser := terragruntparser.Make(hclParser)
	versionFiles := []types.VersionFile{
		{Name: ".terragrunt-version", Parser: flatparser.RetrieveVersion},
		{Name: ".tgswitchrc", Parser: flatparser.RetrieveVersion},
		{Name: ".tgswitch.toml", Parser: tomlparser.RetrieveVersion},
		{Name: asdfparser.ToolFileName, Parser: asdfparser.RetrieveTgVersion},
		{Name: terragruntparser.HCLName, Parser: gruntParser.RetrieveTerragruntVersionConstraintFromHCL},
		{Name: terragruntparser.JSONName, Parser: gruntParser.RetrieveTerragruntVersionConstraintFromJSON},
	}

	return versionmanager.Make(conf, config.TgPrefix, "Terragrunt", nil, tgRetriever, versionFiles)
}

func BuildTofuManager(conf *config.Config, hclParser *hclparse.Parser) versionmanager.VersionManager {
	tofuRetriever := tofuretriever.Make(conf)
	gruntParser := terragruntparser.Make(hclParser)
	versionFiles := []types.VersionFile{
		{Name: ".opentofu-version", Parser: flatparser.RetrieveVersion},
		{Name: asdfparser.ToolFileName, Parser: asdfparser.RetrieveTofuVersion},
		{Name: terragruntparser.HCLName, Parser: gruntParser.RetrieveTerraformVersionConstraintFromHCL},
		{Name: terragruntparser.JSONName, Parser: gruntParser.RetrieveTerraformVersionConstraintFromJSON},
	}

	iacExts := []iacparser.ExtDescription{
		{Value: ".tofu", Parser: hclParser.ParseHCLFile},
		{Value: ".tofu.json", Parser: hclParser.ParseJSONFile},
		{Value: ".tf", Parser: hclParser.ParseHCLFile},
		{Value: ".tf.json", Parser: hclParser.ParseJSONFile},
	}

	return versionmanager.Make(conf, config.TofuenvPrefix, "OpenTofu", iacExts, tofuRetriever, versionFiles)
}
