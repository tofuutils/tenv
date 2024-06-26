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

package proxy

import (
	"fmt"
	"os"

	"github.com/tofuutils/tenv/v2/config"
	"github.com/tofuutils/tenv/v2/config/cmdconst"
	"github.com/tofuutils/tenv/v2/versionmanager/builder"
	terragruntparser "github.com/tofuutils/tenv/v2/versionmanager/semantic/parser/terragrunt"
)

func ExecAgnostic(conf *config.Config, builders map[string]builder.BuilderFunc, gruntParser terragruntparser.TerragruntParser, cmdArgs []string) {
	conf.InitDisplayer(true)
	manager := builders[cmdconst.TofuName](conf, gruntParser)
	detectedVersion, err := manager.ResolveWithVersionFiles()
	if err != nil {
		fmt.Println("Failed to resolve a version allowing to call tofu :", err) //nolint
		os.Exit(1)
	}

	execName := cmdconst.TofuName
	if detectedVersion == "" {
		execName = cmdconst.TerraformName
		manager = builders[cmdconst.TerraformName](conf, gruntParser)
		detectedVersion, err = manager.ResolveWithVersionFiles()
		if err != nil {
			fmt.Println("Failed to resolve a version allowing to call terraform :", err) //nolint
			os.Exit(1)
		}

		if detectedVersion == "" {
			fmt.Println("No version files found corresponding to opentofu or terraform") //nolint
			os.Exit(1)
		}
	}

	installPath, err := manager.InstallPath()
	if err != nil {
		fmt.Println("Failed to create installation directory for", execName, ":", err) //nolint
		os.Exit(1)
	}

	detectedVersion, err = manager.Evaluate(detectedVersion, true)
	if err != nil {
		fmt.Println("Failed to evaluate the requested version in a specific version allowing to call", execName, ":", err) //nolint
		os.Exit(1)
	}

	RunCmd(installPath, detectedVersion, execName, cmdArgs, conf.GithubActions)
}
