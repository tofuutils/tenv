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
	"os/exec"
	"path"

	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/versionmanager"
)

func ExecProxy(builderFunc func(*config.Config) versionmanager.VersionManager, execName string) {
	conf, err := config.InitConfigFromEnv()
	if err != nil {
		fmt.Println("Configuration error :", err)
		os.Exit(1)
	}

	versionManager := builderFunc(&conf)
	detectedVersion, err := versionManager.Detect()
	if err != nil {
		fmt.Println("Failed to detect a version allowing to call", execName, ":", err)
		os.Exit(1)
	}

	cmdArgs := os.Args[1:]
	// proxy to selected version
	cmd := exec.Command(path.Join(versionManager.InstallPath(), detectedVersion, execName), cmdArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Println("Failure during", execName, "call :", err)
	}
}
