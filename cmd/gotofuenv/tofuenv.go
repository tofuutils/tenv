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

	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/versionmanager"
	"github.com/dvaumoron/gotofuenv/versionmanager/builder"
	"github.com/spf13/cobra"
)

const (
	rootVersionHelp = "Display gotofuenv current version."
	tfHelp          = "subcommands that help manage several version of Terraform (https://www.terraform.io)."
)

// can be overridden with ldflags
var version = "dev"

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
		Use:     "gotofuenv",
		Long:    "gotofuenv help manage several version of OpenTofu (https://opentofu.org).",
		Version: version,
	}

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&conf.RootPath, "root-path", "r", conf.RootPath, "local path to install versions of OpenTofu and Terraform")
	flags.StringVarP(&conf.GithubToken, "github-token", "t", "", "GitHub token (increases GitHub REST API rate limits)")
	flags.BoolVarP(&conf.Verbose, "verbose", "v", conf.Verbose, "verbose output")

	rootCmd.AddCommand(newVersionCmd())
	initSubCmds(rootCmd, conf, builder.BuildTofuManager(conf), config.TofuRemoteUrlEnvName, &conf.TofuRemoteUrl)

	tfCmd := &cobra.Command{
		Use:   "tf",
		Short: tfHelp,
		Long:  tfHelp,
	}

	initSubCmds(tfCmd, conf, builder.BuildTfManager(conf), config.TfRemoteUrlEnvName, &conf.TfRemoteUrl)

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
			fmt.Println("gotofuenv", version)
		},
	}
}

func initSubCmds(cmd *cobra.Command, conf *config.Config, versionManager versionmanager.VersionManager, remoteEnvName string, pRemote *string) {
	cmd.AddCommand(newDetectCmd(versionManager))
	cmd.AddCommand(newInstallCmd(conf, versionManager, remoteEnvName, pRemote))
	cmd.AddCommand(newListCmd(conf, versionManager))
	cmd.AddCommand(newListRemoteCmd(conf, versionManager, remoteEnvName, pRemote))
	cmd.AddCommand(newResetCmd(conf, versionManager))
	cmd.AddCommand(newUninstallCmd(conf, versionManager))
	cmd.AddCommand(newUseCmd(conf, versionManager, remoteEnvName, pRemote))
}
