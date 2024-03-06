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

package main

import (
	"fmt"
	"os"

	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/versionmanager/builder"
	"github.com/tofuutils/tenv/versionmanager/proxy"
	terragruntparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/terragrunt"
)

func main() {
	conf, err := config.InitConfigFromEnv()
	if err != nil {
		fmt.Println("Configuration error :", err) //nolint
		os.Exit(1)
	}

	conf.InitDisplayer(true)
	gruntParser := terragruntparser.Make()
	tofuManager := builder.BuildTofuManager(&conf, gruntParser)
	detectedVersion, err := tofuManager.ResolveWithVersionFiles()
	if err != nil {
		fmt.Println("Failed to resolve a version allowing to call tofu :", err) //nolint
		os.Exit(1)
	}

	execName := ""
	installPath := ""
	if detectedVersion == "" {
		terraformManager := builder.BuildTfManager(&conf, gruntParser)
		detectedVersion, err = terraformManager.ResolveWithVersionFiles()
		if err != nil {
			fmt.Println("Failed to resolve a version allowing to call terraform :", err) //nolint
			os.Exit(1)
		}

		if detectedVersion == "" {
			fmt.Println("No version files found corresponding to opentofu or terraform") //nolint
			os.Exit(1)
		}

		execName = config.TerraformName
		installPath = terraformManager.InstallPath()
		detectedVersion, err = terraformManager.Evaluate(detectedVersion, true)
	} else {
		execName = config.TofuName
		installPath = tofuManager.InstallPath()
		detectedVersion, err = tofuManager.Evaluate(detectedVersion, true)
	}
	if err != nil {
		fmt.Println("Failed to evaluate the requested version in a specific version allowing to call", execName, ":", err) //nolint
		os.Exit(1)
	}

	proxy.RunCmd(installPath, detectedVersion, execName)
}
