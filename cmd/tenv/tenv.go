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
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/cobra"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/config/envname"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager"
	"github.com/tofuutils/tenv/v4/versionmanager/builder"
	"github.com/tofuutils/tenv/v4/versionmanager/proxy"
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

	rwPerm = 0o600
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

	hclParser := hclparse.NewParser()
	manageNoArgsCmd(&conf, hclParser)     // call os.Exit when necessary
	manageHiddenCallCmd(&conf, hclParser) // proxy call use os.Exit when called

	if err = initRootCmd(&conf, hclParser).Execute(); err != nil {
		os.Exit(1)
	}
}

func initRootCmd(conf *config.Config, hclParser *hclparse.Parser) *cobra.Command {
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
		needToken: true, remoteEnvName: envname.TofuRemoteURL,
		pRemote: &conf.Tofu.RemoteURL, pPublicKeyPath: &conf.TofuKeyPath,
	}
	initSubCmds(tofuCmd, builder.BuildTofuManager(conf, hclParser), tofuParams)

	rootCmd.AddCommand(tofuCmd)

	tfCmd := &cobra.Command{
		Use:     "tf",
		Aliases: []string{cmdconst.TerraformName},
		Short:   tfHelp,
		Long:    tfHelp,
	}

	tfParams := subCmdParams{
		needToken: false, remoteEnvName: envname.TfRemoteURL,
		pRemote: &conf.Tf.RemoteURL, pPublicKeyPath: &conf.TfKeyPath,
	}
	initSubCmds(tfCmd, builder.BuildTfManager(conf, hclParser), tfParams)

	rootCmd.AddCommand(tfCmd)

	tgCmd := &cobra.Command{
		Use:     "tg",
		Aliases: []string{cmdconst.TerragruntName},
		Short:   tgHelp,
		Long:    tgHelp,
	}

	tgParams := subCmdParams{
		needToken: true, remoteEnvName: envname.TgRemoteURL, pRemote: &conf.Tg.RemoteURL,
	}
	initSubCmds(tgCmd, builder.BuildTgManager(conf, hclParser), tgParams)

	rootCmd.AddCommand(tgCmd)

	atmosCmd := &cobra.Command{
		Use:     "at",
		Aliases: []string{cmdconst.AtmosName},
		Short:   atmosHelp,
		Long:    atmosHelp,
	}

	atmosParams := subCmdParams{
		needToken: true, remoteEnvName: envname.AtmosRemoteURL, pRemote: &conf.Atmos.RemoteURL,
	}
	initSubCmds(atmosCmd, builder.BuildAtmosManager(conf, hclParser), atmosParams)

	rootCmd.AddCommand(atmosCmd)

	return rootCmd
}

func manageNoArgsCmd(conf *config.Config, hclParser *hclparse.Parser) {
	if len(os.Args) > 1 {
		return
	}

	ctx := context.Background()
	if err := toolUI(ctx, conf, hclParser); err != nil {
		loghelper.StdDisplay(err.Error())

		os.Exit(1)
	}

	os.Exit(0)
}

func manageHiddenCallCmd(conf *config.Config, hclParser *hclparse.Parser) {
	if len(os.Args) < 3 || os.Args[1] != cmdconst.CallSubCmd {
		return
	}

	calledNamed, cmdArgs := os.Args[2], os.Args[3:]
	if builderFunc, ok := builder.Builders[calledNamed]; ok {
		proxy.Exec(conf, builderFunc, hclParser, calledNamed, cmdArgs)
	} else if calledNamed == cmdconst.AgnosticName {
		proxy.ExecAgnostic(conf, hclParser, cmdArgs)
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:          versionName,
		Short:        rootVersionHelp,
		Long:         rootVersionHelp,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		Run: func(_ *cobra.Command, _ []string) {
			loghelper.StdDisplay(loghelper.Concat(cmdconst.TenvName, " ", versionName, " ", version))
		},
	}
}

func newUpdatePathCmd(gha bool) *cobra.Command {
	return &cobra.Command{
		Use:          "update-path",
		Short:        updatePathHelp,
		Long:         updatePathHelp,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			execPath, err := os.Executable()
			if err != nil {
				return err
			}

			execDirPath := filepath.Dir(execPath)
			if gha {
				pathfilePath := os.Getenv("GITHUB_PATH")
				if pathfilePath != "" {
					pathfile, err := os.OpenFile(pathfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, rwPerm)
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
			loghelper.StdDisplay(pathBuilder.String())

			return nil
		},
	}
}

func initSubCmds(cmd *cobra.Command, versionManager versionmanager.VersionManager, params subCmdParams) {
	cmd.AddCommand(newConstraintCmd(versionManager))
	cmd.AddCommand(newDetectCmd(versionManager, params))
	cmd.AddCommand(newInstallCmd(versionManager, params))
	cmd.AddCommand(newListCmd(versionManager))
	cmd.AddCommand(newListRemoteCmd(versionManager, params))
	cmd.AddCommand(newResetCmd(versionManager))
	cmd.AddCommand(newUninstallCmd(versionManager))
	cmd.AddCommand(newUseCmd(versionManager, params))
}
