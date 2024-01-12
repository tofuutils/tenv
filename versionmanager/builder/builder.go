/*
 *
 * Copyright 2024 gotofuenv authors.
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
	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/versionmanager"
	terraformretriever "github.com/dvaumoron/gotofuenv/versionmanager/terraform"
	tofuretriever "github.com/dvaumoron/gotofuenv/versionmanager/tofu"
)

func BuildTfManager(conf *config.Config) versionmanager.VersionManager {
	tfRetriever := terraformretriever.MakeTerraformRetriever(conf)
	return versionmanager.MakeVersionManager(conf, config.TfFolderName, tfRetriever, config.TfVersionEnvName, config.TfVersionFileName)
}

func BuildTofuManager(conf *config.Config) versionmanager.VersionManager {
	tofuRetriever := tofuretriever.MakeTofuRetriever(conf)
	return versionmanager.MakeVersionManager(conf, config.TofuFolderName, tofuRetriever, config.TofuVersionEnvName, config.TofuVersionFileName)
}
