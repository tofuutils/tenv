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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/versionmanager"
	"github.com/tofuutils/tenv/versionmanager/builder"
	terragruntparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/terragrunt"
)

const (
	versionName     = "version"
	rootVersionHelp = "Display tenv current version."
	updatePathHelp  = "Display PATH updated with tenv directory location first."

	helpPrefix = "Subcommand to manage several versions of "
	atmosHelp  = helpPrefix + "Atmos (https://atmos.tools)."
	tfHelp     = helpPrefix + "Terraform (https://www.terraform.io)."
	tgHelp     = helpPrefix + "Terragrunt (https://terragrunt.gruntwork.io)."
	tofuHelp   = helpPrefix + "OpenTofu (https://opentofu.org)."

	pathEnvName = "PATH"
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
		Long:    "tenv help manage several versions of OpenTofu (https://opentofu.org), Terraform (https://www.terraform.io) and Terragrunt (https://terragrunt.gruntwork.io).",
		Version: version,
	}

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&conf.ForceQuiet, "quiet", "q", conf.ForceQuiet, "no unnecessary output (and no log)")
	flags.StringVarP(&conf.RootPath, "root-path", "r", conf.RootPath, "local path to install versions of OpenTofu, Terraform and Terragrunt")
	flags.BoolVarP(&conf.DisplayVerbose, "verbose", "v", false, "verbose output (and set log level to Trace)")

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newUpdatePathCmd())

	tofuParams := subCmdParams{
		deprecated: true, // direct use should display a deprecation message
		needToken:  true, remoteEnvName: config.TofuRemoteURLEnvName,
		pRemote: &conf.Tofu.RemoteURL, pPublicKeyPath: &conf.TofuKeyPath,
	}
	gruntParser := terragruntparser.Make()
	tofuManager := builder.BuildTofuManager(conf, gruntParser)
	initSubCmds(rootCmd, conf, tofuManager, tofuParams) // add tofu management at root level

	tofuCmd := &cobra.Command{
		Use:     config.TofuName,
		Aliases: []string{"opentofu"},
		Short:   tofuHelp,
		Long:    tofuHelp,
	}
	tofuParams.deprecated = false // usage with tofu subcommand are ok
	initSubCmds(tofuCmd, conf, tofuManager, tofuParams)

	rootCmd.AddCommand(tofuCmd) // add tofu management as subcommand

	tfCmd := &cobra.Command{
		Use:     "tf",
		Aliases: []string{config.TerraformName},
		Short:   tfHelp,
		Long:    tfHelp,
	}

	tfParams := subCmdParams{
		needToken: false, remoteEnvName: config.TfRemoteURLEnvName,
		pRemote: &conf.Tf.RemoteURL, pPublicKeyPath: &conf.TfKeyPath,
	}
	initSubCmds(tfCmd, conf, builder.BuildTfManager(conf, gruntParser), tfParams)

	rootCmd.AddCommand(tfCmd)

	tgCmd := &cobra.Command{
		Use:     "tg",
		Aliases: []string{config.TerragruntName},
		Short:   tgHelp,
		Long:    tgHelp,
	}

	tgParams := subCmdParams{
		needToken: true, remoteEnvName: config.TgRemoteURLEnvName, pRemote: &conf.Tg.RemoteURL,
	}
	initSubCmds(tgCmd, conf, builder.BuildTgManager(conf, gruntParser), tgParams)

	rootCmd.AddCommand(tgCmd)

	atmosCmd := &cobra.Command{
		Use:   config.AtmosName,
		Short: atmosHelp,
		Long:  atmosHelp,
	}

	atmosParams := subCmdParams{
		needToken: true, remoteEnvName: config.AtmosRemoteURLEnvName, pRemote: &conf.Atmos.RemoteURL,
	}
	initSubCmds(atmosCmd, conf, builder.BuildAtmosManager(conf, gruntParser), atmosParams)

	rootCmd.AddCommand(atmosCmd)

	return rootCmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   versionName,
		Short: rootVersionHelp,
		Long:  rootVersionHelp,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(config.TenvName, versionName, version) //nolint
		},
	}
}

func newUpdatePathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-path",
		Short: updatePathHelp,
		Long:  updatePathHelp,
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			execPath, err := os.Executable()
			if err != nil {
				return nil
			}

			gha, err := config.GetenvBool(false, config.GithubActionsEnvName)
			if err != nil {
				return err
			}

			execDirPath := filepath.Dir(execPath)
			if gha {
				pathfilePath := os.Getenv("GITHUB_PATH")
				if pathfilePath != "" {
					pathfile, err := os.OpenFile(pathfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						return err
					}
					defer pathfile.Close()

					_, err = pathfile.Write(append([]byte(execDirPath), '\n'))
					if err != nil {
						return err
					}
				}
			}

			var pathBuilder strings.Builder
			pathBuilder.WriteString(execDirPath)
			pathBuilder.WriteRune(os.PathListSeparator)
			pathBuilder.WriteString(os.Getenv(pathEnvName))
			fmt.Println(pathBuilder.String()) //nolint

			return nil
		},
	}
}

func initSubCmds(cmd *cobra.Command, conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) {
	cmd.AddCommand(newConstraintCmd(conf, versionManager, params))
	cmd.AddCommand(newDetectCmd(conf, versionManager, params))
	cmd.AddCommand(newInstallCmd(conf, versionManager, params))
	cmd.AddCommand(newListCmd(conf, versionManager, params))
	cmd.AddCommand(newListRemoteCmd(conf, versionManager, params))
	cmd.AddCommand(newResetCmd(conf, versionManager, params))
	cmd.AddCommand(newUninstallCmd(conf, versionManager, params))
	cmd.AddCommand(newUseCmd(conf, versionManager, params))
}
