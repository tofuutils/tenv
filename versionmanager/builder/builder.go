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
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/versionmanager"
	terraformretriever "github.com/tofuutils/tenv/versionmanager/retriever/terraform"
	terragruntretriever "github.com/tofuutils/tenv/versionmanager/retriever/terragrunt"
	tofuretriever "github.com/tofuutils/tenv/versionmanager/retriever/tofu"
	"github.com/tofuutils/tenv/versionmanager/semantic"
	flatparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/flat"
	tgswitchparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/tgswitch"
)

func BuildTfManager(conf *config.Config) versionmanager.VersionManager {
	tfRetriever := terraformretriever.NewTerraformRetriever(conf)
	versionFiles := []semantic.VersionFile{{Name: ".terraform-version", Parser: flatparser.RetrieveVersionFromFile}, {Name: ".tfswitchrc", Parser: flatparser.RetrieveVersionFromFile}}

	return versionmanager.MakeVersionManager(conf, "Terraform", semantic.TfPredicateReaders, tfRetriever, config.TfVersionEnvName, versionFiles)
}

func BuildTgManager(conf *config.Config) versionmanager.VersionManager {
	tgRetriever := terragruntretriever.NewTerragruntRetriever(conf)
	versionFiles := []semantic.VersionFile{{Name: ".terragrunt-version", Parser: flatparser.RetrieveVersionFromFile}, {Name: ".tgswitchrc", Parser: flatparser.RetrieveVersionFromFile}, {Name: ".tgswitch.toml", Parser: tgswitchparser.RetrieveTerraguntVersionFromFile}}

	return versionmanager.MakeVersionManager(conf, "Terragrunt", semantic.TgPredicateReaders, tgRetriever, config.TgVersionEnvName, versionFiles)
}

func BuildTofuManager(conf *config.Config) versionmanager.VersionManager {
	tofuRetriever := tofuretriever.NewTofuRetriever(conf)
	versionFiles := []semantic.VersionFile{{Name: ".opentofu-version", Parser: flatparser.RetrieveVersionFromFile}}

	return versionmanager.MakeVersionManager(conf, "OpenTofu", semantic.TfPredicateReaders, tofuRetriever, config.TofuVersionEnvName, versionFiles)
}
