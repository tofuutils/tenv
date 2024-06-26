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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tofuutils/tenv/v2/config"
	"github.com/tofuutils/tenv/v2/versionmanager/builder"
	cmdproxy "github.com/tofuutils/tenv/v2/versionmanager/proxy/cmd"
	terragruntparser "github.com/tofuutils/tenv/v2/versionmanager/semantic/parser/terragrunt"
)

var errDelimiter = errors.New("key and value should not contains delimiter")

func Exec(conf *config.Config, builderFunc builder.BuilderFunc, gruntParser terragruntparser.TerragruntParser, execName string, cmdArgs []string) {
	conf.InitDisplayer(true)
	versionManager := builderFunc(conf, gruntParser)
	detectedVersion, err := versionManager.Detect(true)
	if err != nil {
		fmt.Println("Failed to detect a version allowing to call", execName, ":", err) //nolint
		os.Exit(1)
	}

	installPath, err := versionManager.InstallPath()
	if err != nil {
		fmt.Println("Failed to create installation directory for", execName, ":", err) //nolint
		os.Exit(1)
	}

	RunCmd(installPath, detectedVersion, execName, cmdArgs)
}

func RunCmd(installPath string, detectedVersion string, execName string, cmdArgs []string) {
	cmdproxy.Run(filepath.Join(installPath, detectedVersion, execName), cmdArgs)
}
