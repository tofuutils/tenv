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
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/versionmanager"
	"github.com/tofuutils/tenv/versionmanager/semantic"
)

const deprecationMsg = "Direct usage of this subcommand on tenv is deprecated, you should use tofu subcommand instead.\n\n"

func newDetectCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Display ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" current version.")
	detectHelp := descBuilder.String()

	descBuilder.Reset()
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Display ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" current version.")

	detectCmd := &cobra.Command{
		Use:   "detect",
		Short: detectHelp,
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			conf.LogLevelUpdate()
			addDeprecationMsg(conf, params)

			detectedVersion, err := versionManager.Detect(false)
			if err != nil {
				return err
			}
			fmt.Println(versionManager.FolderName, detectedVersion, "will be run from this directory.") //nolint

			return nil
		},
	}

	flags := detectCmd.Flags()
	addInstallationFlags(flags, conf, params)
	addOptionalInstallationFlags(flags, conf, params)
	addRemoteFlags(flags, conf, params)

	return detectCmd
}

func newInstallCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Install a specific version of ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteByte('.')
	shortMsg := descBuilder.String()

	descBuilder.Reset()
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Install a specific version of ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" (into TENV_ROOT directory from ")
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url).\n\nWithout parameter the version to use is resolved automatically via ")
	descBuilder.WriteString(versionManager.VersionEnvName)
	descBuilder.WriteString(` or version files
(searched in working directory, its parents, user home directory or TENV_ROOT directory).
Use "latest" when none are found.

If a parameter is passed, available options:
- an exact Semver 2.0.0 version string to install
- a version constraint string (checked against version available at `)
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url)\n- latest, latest-stable or latest-pre (checked against version available at ")
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url)\n- latest-allowed or min-required to scan your ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" files to detect which version is maximally allowed or minimally required.")

	installCmd := &cobra.Command{
		Use:   "install [version]",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			conf.LogLevelUpdate()
			addDeprecationMsg(conf, params)

			var requestedVersion string
			if len(args) == 0 {
				var err error
				requestedVersion, err = versionManager.Resolve(semantic.LatestKey)
				if err != nil {
					return err
				}
			} else {
				requestedVersion = args[0]
			}

			return versionManager.Install(requestedVersion)
		},
	}

	flags := installCmd.Flags()
	addInstallationFlags(flags, conf, params)
	addRemoteFlags(flags, conf, params)

	return installCmd
}

func newListCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("List installed ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" versions.")
	shortMsg := descBuilder.String()

	descBuilder.Reset()
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("List installed ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" versions (located in TENV_ROOT directory), sorted in ascending version order.")

	reverseOrder := false

	listCmd := &cobra.Command{
		Use:   "list",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			conf.LogLevelUpdate()
			addDeprecationMsg(conf, params)

			versions, err := versionManager.ListLocal(reverseOrder)
			if err != nil {
				return err
			}

			filePath := versionManager.RootVersionFilePath()
			data, err := os.ReadFile(filePath)
			if err != nil && conf.DisplayVerbose {
				fmt.Println("Can not read used version :", err) //nolint
			}
			usedVersion := string(bytes.TrimSpace(data))

			for _, version := range versions {
				if usedVersion == version {
					fmt.Println("*", version, "(set by", filePath+")") //nolint
				} else {
					fmt.Println(" ", version) //nolint
				}
			}
			if conf.DisplayVerbose {
				fmt.Println("found", len(versions), versionManager.FolderName, "version(s) managed by tenv.") //nolint
			}

			return nil
		},
	}

	addDescendingFlag(listCmd.Flags(), &reverseOrder)

	return listCmd
}

func newListRemoteCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("List installable ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" versions.")
	shortMsg := descBuilder.String()

	descBuilder.Reset()
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("List installable ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" versions (from ")
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url), sorted in ascending version order.")

	filterStable := false
	reverseOrder := false

	listRemoteCmd := &cobra.Command{
		Use:   "list-remote",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			conf.LogLevelUpdate()
			addDeprecationMsg(conf, params)

			versions, err := versionManager.ListRemote(reverseOrder)
			if err != nil {
				return err
			}

			countSkipped := 0
			localSet := versionManager.LocalSet()
			for _, version := range versions {
				if filterStable && !semantic.StableVersion(version) {
					countSkipped++

					continue
				}

				if _, installed := localSet[version]; installed {
					fmt.Println(version, "(installed)") //nolint
				} else {
					fmt.Println(version) //nolint
				}
			}
			if conf.DisplayVerbose {
				fmt.Println("found", len(versions), versionManager.FolderName, "version(s) (on", params.remoteEnvName+").") //nolint
				if filterStable {
					fmt.Println(countSkipped, "result(s) hidden (version not stable).") //nolint
				}
			}

			return err
		},
	}

	flags := listRemoteCmd.Flags()
	addDescendingFlag(flags, &reverseOrder)
	addRemoteFlags(flags, conf, params)
	flags.BoolVarP(&filterStable, "stable", "s", false, "display only stable version")

	return listRemoteCmd
}

func newResetCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Reset used version of ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteByte('.')
	shortMsg := descBuilder.String()

	descBuilder.Reset()
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Reset used version of ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" (remove TENV_ROOT/")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString("/version file).")

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			conf.LogLevelUpdate()
			addDeprecationMsg(conf, params)

			return versionManager.Reset()
		},
	}

	return resetCmd
}

func newUninstallCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Uninstall a specific version of ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteByte('.')
	shortMsg := descBuilder.String()

	descBuilder.Reset()
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Uninstall a specific version of ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" (remove it from TENV_ROOT directory).")

	uninstallCmd := &cobra.Command{
		Use:   "uninstall version",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			conf.LogLevelUpdate()
			addDeprecationMsg(conf, params)

			return versionManager.Uninstall(args[0])
		},
	}

	return uninstallCmd
}

func newUseCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Switch the default ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" version to use.")
	shortMsg := descBuilder.String()

	descBuilder.Reset()
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Switch the default ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" version to use (set in TENV_ROOT/")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(`/version file)

Available parameter options:
- an exact Semver 2.0.0 version string to use
- a version constraint string (checked against version available in TENV_ROOT directory)
- latest, latest-stable or latest-pre (checked against version available in TENV_ROOT directory)
- latest-allowed or min-required to scan your `)
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" files to detect which version is maximally allowed or minimally required.")

	workingDir := false

	useCmd := &cobra.Command{
		Use:   "use version",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			conf.LogLevelUpdate()
			addDeprecationMsg(conf, params)

			return versionManager.Use(args[0], workingDir)
		},
	}

	descBuilder.Reset()
	descBuilder.WriteString("create ")
	descBuilder.WriteString(versionManager.VersionFiles[0].Name)
	descBuilder.WriteString(" file in working directory")

	flags := useCmd.Flags()
	addInstallationFlags(flags, conf, params)
	addOptionalInstallationFlags(flags, conf, params)
	addRemoteFlags(flags, conf, params)
	flags.BoolVarP(&workingDir, "working-dir", "w", false, descBuilder.String())

	return useCmd
}

func addDeprecationHelpMsg(descBuilder *strings.Builder, params subCmdParams) {
	if params.deprecated {
		descBuilder.WriteString(deprecationMsg)
	}
}

func addDeprecationMsg(conf *config.Config, params subCmdParams) {
	if params.deprecated {
		conf.Display(deprecationMsg)
	}
}

func addDescendingFlag(flags *pflag.FlagSet, pReverseOrder *bool) {
	flags.BoolVarP(pReverseOrder, "descending", "d", false, "display list in descending version order")
}

func addInstallationFlags(flags *pflag.FlagSet, conf *config.Config, params subCmdParams) {
	flags.StringVarP(&conf.Arch, "arch", "a", conf.Arch, "specify arch for binaries downloading")
	if params.pPublicKeyPath != nil {
		flags.StringVarP(params.pPublicKeyPath, "key-file", "k", "", "local path to PGP public key file (replace check against remote one)")
	}
}

func addOptionalInstallationFlags(flags *pflag.FlagSet, conf *config.Config, params subCmdParams) {
	var descBuilder strings.Builder
	descBuilder.WriteString("force search on versions available at ")
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url")

	flags.BoolVarP(&conf.ForceRemote, "force-remote", "f", conf.ForceRemote, descBuilder.String())
	flags.BoolVarP(&conf.NoInstall, "no-install", "n", conf.NoInstall, "disable installation of missing version")
}

func addRemoteFlags(flags *pflag.FlagSet, conf *config.Config, params subCmdParams) {
	flags.StringVarP(&conf.RemoteConfPath, "remote-conf", "c", conf.RemoteConfPath, "path to remote configuration file (advanced settings)")
	if params.needToken {
		flags.StringVarP(&conf.GithubToken, "github-token", "t", conf.GithubToken, "GitHub token (increases GitHub REST API rate limits)")
	}
	flags.StringVarP(params.pRemote, "remote-url", "u", "", "remote url to install from")
}
