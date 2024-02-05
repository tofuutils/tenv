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
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/versionmanager"
	"github.com/tofuutils/tenv/versionmanager/builder"
)

const (
	rootVersionHelp = "Display tenv current version."
	tfHelp          = "subcommand to manage several versions of Terraform (https://www.terraform.io)."
	tgHelp          = "subcommand to manage several versions of Terragrunt (https://terragrunt.gruntwork.io/)."
	tofuHelp        = "subcommand to manage several versions of OpenTofu (https://opentofu.org)."
)

// can be overridden with ldflags.
var version = "dev"

type subCmdParams struct {
	deprecated     bool
	needToken      bool
	remoteEnvName  string
	pRemote        *string
	pPublicKeyPath *string
}

func main() {
	conf, err := config.InitConfigFromEnv()
	if err != nil {
		fmt.Println("Configuration error :", err) //nolint
		os.Exit(1)
	}

	if err = initRootCmd(&conf).Execute(); err != nil {
		fmt.Println(err) //nolint
		os.Exit(1)
	}
}

func initRootCmd(conf *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     config.TenvName,
		Long:    "tenv help manage several versions of OpenTofu (https://opentofu.org), Terraform (https://www.terraform.io) and Terragrunt (https://terragrunt.gruntwork.io/).",
		Version: version,
	}

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&conf.ForceQuiet, "quiet", "q", conf.ForceQuiet, "no output (and no log)")
	flags.StringVarP(&conf.RootPath, "root-path", "r", conf.RootPath, "local path to install versions of OpenTofu, Terraform and Terragrunt")
	flags.BoolVarP(&conf.DisplayVerbose, "verbose", "v", false, "verbose output (and set log level to Trace)")

	rootCmd.AddCommand(newVersionCmd())
	tofuParams := subCmdParams{
		deprecated: true, // direct use should display a deprecation message
		needToken:  true, remoteEnvName: config.TofuRemoteURLEnvName,
		pRemote: &conf.Tofu.RemoteURL, pPublicKeyPath: &conf.TofuKeyPath,
	}
	tofuManager := builder.BuildTofuManager(conf)
	initSubCmds(rootCmd, conf, tofuManager, tofuParams)

	// Add this in your main function, after the tfCmd and before the tgCmd
	tofuCmd := &cobra.Command{
		Use:   config.TofuName,
		Short: tofuHelp,
		Long:  tofuHelp,
	}
	tofuParams.deprecated = false // usage with tofu subcommand are ok
	initSubCmds(tofuCmd, conf, tofuManager, tofuParams)

	rootCmd.AddCommand(tofuCmd)

	tfCmd := &cobra.Command{
		Use:   "tf",
		Short: tfHelp,
		Long:  tfHelp,
	}

	tfParams := subCmdParams{
		needToken: false, remoteEnvName: config.TfRemoteURLEnvName,
		pRemote: &conf.Tf.RemoteURL, pPublicKeyPath: &conf.TfKeyPath,
	}
	initSubCmds(tfCmd, conf, builder.BuildTfManager(conf), tfParams)

	rootCmd.AddCommand(tfCmd)

	tgCmd := &cobra.Command{
		Use:   "tg",
		Short: tgHelp,
		Long:  tgHelp,
	}

	tgParams := subCmdParams{
		needToken: true, remoteEnvName: config.TgRemoteURLEnvName, pRemote: &conf.Tg.RemoteURL,
	}
	initSubCmds(tgCmd, conf, builder.BuildTgManager(conf), tgParams)

	rootCmd.AddCommand(tgCmd)

	return rootCmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: rootVersionHelp,
		Long:  rootVersionHelp,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(config.TenvName, version) //nolint
		},
	}
}

func initSubCmds(cmd *cobra.Command, conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) {
	cmd.AddCommand(newDetectCmd(conf, versionManager, params))
	cmd.AddCommand(newInstallCmd(conf, versionManager, params))
	cmd.AddCommand(newListCmd(conf, versionManager, params))
	cmd.AddCommand(newListRemoteCmd(conf, versionManager, params))
	cmd.AddCommand(newResetCmd(conf, versionManager, params))
	cmd.AddCommand(newUninstallCmd(conf, versionManager, params))
	cmd.AddCommand(newUseCmd(conf, versionManager, params))
}
