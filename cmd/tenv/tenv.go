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

	"github.com/spf13/cobra"
	"github.com/tofuutils/gotofuenv/config"
	"github.com/tofuutils/gotofuenv/versionmanager"
	"github.com/tofuutils/gotofuenv/versionmanager/builder"
)

const (
	rootVersionHelp = "Display tenv current version."
	tfHelp          = "subcommands that help manage several version of Terraform (https://www.terraform.io)."
)

// can be overridden with ldflags
var version = "dev"

type subCmdParams struct {
	needToken      bool
	remoteEnvName  string
	pRemote        *string
	pPublicKeyPath *string
}

func main() {
	conf, err := config.InitConfigFromEnv()
	if err != nil {
		fmt.Println("Configuration error :", err)
		os.Exit(1)
	}

	if err = initRootCmd(&conf).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initRootCmd(conf *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "tenv",
		Long:    "tenv help manage several version of OpenTofu (https://opentofu.org).",
		Version: version,
	}

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&conf.RootPath, "root-path", "r", conf.RootPath, "local path to install versions of OpenTofu and Terraform")
	flags.BoolVarP(&conf.Verbose, "verbose", "v", conf.Verbose, "verbose output")

	rootCmd.AddCommand(newVersionCmd())
	tofuParams := subCmdParams{
		needToken: true, remoteEnvName: config.TofuRemoteUrlEnvName,
		pRemote: &conf.TofuRemoteUrl, pPublicKeyPath: &conf.TofuKeyPath,
	}
	initSubCmds(rootCmd, conf, builder.BuildTofuManager(conf), tofuParams)

	tfCmd := &cobra.Command{
		Use:   "tf",
		Short: tfHelp,
		Long:  tfHelp,
	}

	tfParams := subCmdParams{
		needToken: false, remoteEnvName: config.TfRemoteUrlEnvName,
		pRemote: &conf.TfRemoteUrl, pPublicKeyPath: &conf.TfKeyPath,
	}
	initSubCmds(tfCmd, conf, builder.BuildTfManager(conf), tfParams)

	rootCmd.AddCommand(tfCmd)

	return rootCmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: rootVersionHelp,
		Long:  rootVersionHelp,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("tenv from GoTofuEnv", version)
		},
	}
}

func initSubCmds(cmd *cobra.Command, conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) {
	cmd.AddCommand(newDetectCmd(conf, versionManager, params))
	cmd.AddCommand(newInstallCmd(conf, versionManager, params))
	cmd.AddCommand(newListCmd(conf, versionManager))
	cmd.AddCommand(newListRemoteCmd(conf, versionManager, params))
	cmd.AddCommand(newResetCmd(conf, versionManager))
	cmd.AddCommand(newUninstallCmd(conf, versionManager))
	cmd.AddCommand(newUseCmd(conf, versionManager, params))
}
