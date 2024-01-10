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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/tofuversion"
)

func main() {
	conf, err := config.InitConfigFromEnv()
	if err != nil {
		fmt.Println("Configuration error :", err)
		os.Exit(1)
	}

	// detect version (can install depending on GOTOFUENV_AUTO_INSTALL)
	configVersion := conf.ResolveVersion(config.LatestAllowedKey)
	detectedVersion, err := tofuversion.Detect(configVersion, &conf)
	if err != nil {
		fmt.Println("Failed to detect an OpenTofu version :", err)
		os.Exit(1)
	}

	cmdArgs := os.Args[1:]
	// proxy to selected version
	cmd := exec.Command(path.Join(conf.InstallPath(), detectedVersion, "tofu"), cmdArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Println("Failure during tofu call :", err)
	}
}
