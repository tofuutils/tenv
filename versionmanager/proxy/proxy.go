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
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/pkg/cmdproxy"
	"github.com/tofuutils/tenv/v4/versionmanager/builder"
	"github.com/tofuutils/tenv/v4/versionmanager/lastuse"
)

const chdirFlagPrefix = "-chdir="

// Always call os.Exit.
func Exec(conf *config.Config, builderFunc builder.Func, hclParser *hclparse.Parser, execName string, cmdArgs []string) {
	conf.InitDisplayer(true)
	versionManager := builderFunc(conf, hclParser)

	updateWorkPath(conf, cmdArgs)

	ctx := context.Background()
	detectedVersion, err := versionManager.Detect(ctx, true)
	if err != nil {
		fmt.Println("Failed to detect a version allowing to call", execName, ":", err) //nolint
		os.Exit(cmdconst.EarlyErrorExitCode)
	}

	installPath, err := versionManager.InstallPath()
	if err != nil {
		fmt.Println("Failed to create installation directory for", execName, ":", err) //nolint
		os.Exit(cmdconst.EarlyErrorExitCode)
	}

	execPath := ExecPath(installPath, detectedVersion, execName, conf)

	cmd := exec.CommandContext(ctx, execPath, cmdArgs...)

	cmdproxy.Run(cmd, conf.GithubActions, conf.Getenv)
}

func ExecPath(installPath string, version string, execName string, conf *config.Config) string {
	versionPath := filepath.Join(installPath, version)
	lastuse.WriteNow(versionPath, conf)

	return filepath.Join(versionPath, execName)
}

func updateWorkPath(conf *config.Config, cmdArgs []string) {
	for _, arg := range cmdArgs {
		if chdirPath, ok := strings.CutPrefix(arg, chdirFlagPrefix); ok {
			conf.WorkPath = chdirPath

			return
		}
	}
}
