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

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/cobra"

	"github.com/tofuutils/tenv/v2/config"
	"github.com/tofuutils/tenv/v2/config/cmdconst"
	"github.com/tofuutils/tenv/v2/pkg/loghelper"
	"github.com/tofuutils/tenv/v2/versionmanager"
	"github.com/tofuutils/tenv/v2/versionmanager/builder"
	"github.com/tofuutils/tenv/v2/versionmanager/proxy"
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
	needToken      bool
	remoteEnvName  string
	pRemote        *string
	pPublicKeyPath *string
}

func main() {
	conf, err := config.InitConfigFromEnv()
	if err != nil {
		loghelper.StdDisplay(loghelper.Concat("Configuration error : ", err.Error()))
		os.Exit(1)
	}

	builders := map[string]builder.BuilderFunc{
		cmdconst.TofuName:       builder.BuildTofuManager,
		cmdconst.TerraformName:  builder.BuildTfManager,
		cmdconst.TerragruntName: builder.BuildTgManager,
		cmdconst.AtmosName:      builder.BuildAtmosManager,
	}

	hclParser := hclparse.NewParser()
	manageNoArgsCmd(&conf, builders, hclParser)     // call os.Exit when necessary
	manageHiddenCallCmd(&conf, builders, hclParser) // proxy call use os.Exit when called

	if err = initRootCmd(&conf, builders, hclParser).Execute(); err != nil {
		loghelper.StdDisplay(err.Error())
		os.Exit(1)
	}
}

func initRootCmd(conf *config.Config, builders map[string]builder.BuilderFunc, hclParser *hclparse.Parser) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     cmdconst.TenvName,
		Long:    "tenv help manage several versions of OpenTofu (https://opentofu.org), Terraform (https://www.terraform.io), Terragrunt (https://terragrunt.gruntwork.io), and Atmos (https://atmos.tools/).",
		Version: version,
	}

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&conf.ForceQuiet, "quiet", "q", conf.ForceQuiet, "no unnecessary output (and no log)")
	flags.StringVarP(&conf.RootPath, "root-path", "r", conf.RootPath, "local path to install versions of OpenTofu, Terraform, Terragrunt, and Atmos")
	flags.BoolVarP(&conf.DisplayVerbose, "verbose", "v", false, "verbose output (and set log level to Trace)")

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newUpdatePathCmd(conf.GithubActions))

	tofuCmd := &cobra.Command{
		Use:     cmdconst.TofuName,
		Aliases: []string{"opentofu"},
		Short:   tofuHelp,
		Long:    tofuHelp,
	}

	tofuParams := subCmdParams{
		needToken: true, remoteEnvName: config.TofuRemoteURLEnvName,
		pRemote: &conf.Tofu.RemoteURL, pPublicKeyPath: &conf.TofuKeyPath,
	}
	initSubCmds(tofuCmd, conf, builders[cmdconst.TofuName](conf, hclParser), tofuParams)

	rootCmd.AddCommand(tofuCmd)

	tfCmd := &cobra.Command{
		Use:     "tf",
		Aliases: []string{cmdconst.TerraformName},
		Short:   tfHelp,
		Long:    tfHelp,
	}

	tfParams := subCmdParams{
		needToken: false, remoteEnvName: config.TfRemoteURLEnvName,
		pRemote: &conf.Tf.RemoteURL, pPublicKeyPath: &conf.TfKeyPath,
	}
	initSubCmds(tfCmd, conf, builders[cmdconst.TerraformName](conf, hclParser), tfParams)

	rootCmd.AddCommand(tfCmd)

	tgCmd := &cobra.Command{
		Use:     "tg",
		Aliases: []string{cmdconst.TerragruntName},
		Short:   tgHelp,
		Long:    tgHelp,
	}

	tgParams := subCmdParams{
		needToken: true, remoteEnvName: config.TgRemoteURLEnvName, pRemote: &conf.Tg.RemoteURL,
	}
	initSubCmds(tgCmd, conf, builders[cmdconst.TerragruntName](conf, hclParser), tgParams)

	rootCmd.AddCommand(tgCmd)

	atmosCmd := &cobra.Command{
		Use:     "at",
		Aliases: []string{cmdconst.AtmosName},
		Short:   atmosHelp,
		Long:    atmosHelp,
	}

	atmosParams := subCmdParams{
		needToken: true, remoteEnvName: config.AtmosRemoteURLEnvName, pRemote: &conf.Atmos.RemoteURL,
	}
	initSubCmds(atmosCmd, conf, builders[cmdconst.AtmosName](conf, hclParser), atmosParams)

	rootCmd.AddCommand(atmosCmd)

	return rootCmd
}

func manageNoArgsCmd(conf *config.Config, builders map[string]builder.BuilderFunc, hclParser *hclparse.Parser) {
	if len(os.Args) > 1 {
		return
	}

	if err := toolUI(conf, builders, hclParser); err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	os.Exit(0)
}

func manageHiddenCallCmd(conf *config.Config, builders map[string]builder.BuilderFunc, hclParser *hclparse.Parser) {
	if len(os.Args) < 3 || os.Args[1] != cmdconst.CallSubCmd {
		return
	}

	calledNamed, cmdArgs := os.Args[2], os.Args[3:]
	if builder, ok := builders[calledNamed]; ok {
		proxy.Exec(conf, builder, hclParser, calledNamed, cmdArgs)
	} else if calledNamed == cmdconst.AgnosticName {
		proxy.ExecAgnostic(conf, builders, hclParser, cmdArgs)
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   versionName,
		Short: rootVersionHelp,
		Long:  rootVersionHelp,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			loghelper.StdDisplay(loghelper.Concat(cmdconst.TenvName, " ", versionName, " ", version))
		},
	}
}

func newUpdatePathCmd(gha bool) *cobra.Command {
	return &cobra.Command{
		Use:   "update-path",
		Short: updatePathHelp,
		Long:  updatePathHelp,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			execPath, err := os.Executable()
			if err != nil {
				loghelper.StdDisplay(err.Error())

				return
			}

			execDirPath := filepath.Dir(execPath)
			if gha {
				pathfilePath := os.Getenv("GITHUB_PATH")
				if pathfilePath != "" {
					pathfile, err := os.OpenFile(pathfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
					if err != nil {
						return
					}
					defer pathfile.Close()

					_, err = pathfile.Write(append([]byte(execDirPath), '\n'))
					if err != nil {
						return
					}
				}
			}

			var pathBuilder strings.Builder
			pathBuilder.WriteString(execDirPath)
			pathBuilder.WriteRune(os.PathListSeparator)
			pathBuilder.WriteString(os.Getenv(pathEnvName))
			loghelper.StdDisplay(pathBuilder.String())
		},
	}
}

func initSubCmds(cmd *cobra.Command, conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) {
	cmd.AddCommand(newConstraintCmd(conf, versionManager))
	cmd.AddCommand(newDetectCmd(conf, versionManager, params))
	cmd.AddCommand(newInstallCmd(conf, versionManager, params))
	cmd.AddCommand(newListCmd(conf, versionManager))
	cmd.AddCommand(newListRemoteCmd(conf, versionManager, params))
	cmd.AddCommand(newResetCmd(conf, versionManager))
	cmd.AddCommand(newUninstallCmd(conf, versionManager))
	cmd.AddCommand(newUseCmd(conf, versionManager, params))
}
