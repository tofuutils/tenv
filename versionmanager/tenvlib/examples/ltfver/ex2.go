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
	"context"
	"fmt"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/config/cmdconst"
	"github.com/tofuutils/tenv/v3/versionmanager/semantic"
	"github.com/tofuutils/tenv/v3/versionmanager/tenvlib"
)

func main() {
	conf, err := config.DefaultConfig() // does not read environment variables
	if err != nil {
		fmt.Println("init failed :", err)

		return
	}

	conf.SkipInstall = false // tenvlib.AutoInstall option equivalent

	tenv, err := tenvlib.Make(tenvlib.WithConfig(&conf), tenvlib.DisableDisplay)
	if err != nil {
		fmt.Println("should not occur when calling WithConfig :", err)

		return
	}

	ctx := context.Background()
	version, err := tenv.Evaluate(ctx, cmdconst.TerraformName, semantic.LatestKey)
	if err != nil {
		fmt.Println("eval failed :", err)

		return
	}

	conf.ForceRemote = true

	remoteVersion, err := tenv.Evaluate(ctx, cmdconst.TerraformName, semantic.LatestKey)
	if err != nil {
		fmt.Println("eval remote failed :", err)

		return
	}

	if version != remoteVersion {
		err = tenv.Uninstall(ctx, cmdconst.TerraformName, version)
		if err != nil {
			fmt.Println("uninstall failed :", err)
		}
	}

	fmt.Println("Last Terraform version :", version, "(local),", remoteVersion, "(remote)")
}
