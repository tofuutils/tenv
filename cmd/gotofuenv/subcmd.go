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
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/versionmanager"
	"github.com/dvaumoron/gotofuenv/versionmanager/semantic"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newDetectCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Display ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" current version.")
	detectHelp := descBuilder.String()

	detectCmd := &cobra.Command{
		Use:   "detect",
		Short: detectHelp,
		Long:  detectHelp,
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			detectedVersion, err := versionManager.Detect()
			if err != nil {
				return err
			}
			fmt.Println(versionManager.FolderName, detectedVersion, "will be run from this directory.")
			return nil
		},
	}

	flags := detectCmd.Flags()
	addForceRemoteAndNoInstallFlags(flags, conf, params)
	addKeyFlag(flags, params)
	addRemoteUrlFlag(flags, conf, params)

	return detectCmd
}

func newInstallCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Install a specific version of ")
	descBuilder.WriteString(versionManager.FolderName)
	shortMsg := descBuilder.String() + "."

	descBuilder.WriteString(" (into TOFUENV_ROOT directory from ")
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url).\n\nWithout parameter the version to use is resolved automatically via ")
	descBuilder.WriteString(versionManager.VersionEnvName)
	descBuilder.WriteString(" or ")
	descBuilder.WriteString(versionManager.VersionFileName)
	descBuilder.WriteString(` files
(searched in working directory, user home directory and TOFUENV_ROOT directory).
Use "latest-stable" when none are found.

If a parameter is passed, available options:
- an exact Semver 2.0.0 version string to install
- a version constraint string (checked against version available at `)
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url)\n- latest or latest-stable (checked against version available at ")
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url)\n- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.")

	installCmd := &cobra.Command{
		Use:   "install [version]",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			requestedVersion := ""
			if len(args) == 0 {
				requestedVersion = versionManager.Resolve(semantic.LatestStableKey)
			} else {
				requestedVersion = args[0]
			}
			return versionManager.Install(requestedVersion)
		},
	}

	flags := installCmd.Flags()
	addKeyFlag(flags, params)
	addRemoteUrlFlag(flags, conf, params)

	return installCmd
}

func newListCmd(conf *config.Config, versionManager versionmanager.VersionManager) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("List installed ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" versions")
	shortMsg := descBuilder.String() + "."

	descBuilder.WriteString(" (located in TOFUENV_ROOT directory), sorted in ascending version order.")

	reverseOrder := false

	listCmd := &cobra.Command{
		Use:   "list",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			versions, err := versionManager.ListLocal(reverseOrder)
			if err != nil {
				return err
			}

			filePath := versionManager.RootVersionFilePath()
			data, err := os.ReadFile(filePath)
			if err != nil && conf.Verbose {
				fmt.Println("Can not read used version :", err)
			}
			usedVersion := string(bytes.TrimSpace(data))

			for _, version := range versions {
				if usedVersion == version {
					fmt.Println("*", version, "(set by", filePath+")")
				} else {
					fmt.Println(" ", version)
				}
			}
			if conf.Verbose {
				fmt.Println("found", len(versions), versionManager.FolderName, "version(s) managed by gotofuenv.")
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
	descBuilder.WriteString(" versions")
	shortMsg := descBuilder.String() + "."

	descBuilder.WriteString(" (from ")
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
					fmt.Println(version, "(installed)")
				} else {
					fmt.Println(version)
				}
			}
			if conf.Verbose {
				fmt.Println("found", len(versions), versionManager.FolderName, "version(s) (on", *params.pRemote+").")
				if filterStable {
					fmt.Println(countSkipped, "result(s) hidden (version not stable).")
				}
			}
			return err
		},
	}

	flags := listRemoteCmd.Flags()
	addDescendingFlag(flags, &reverseOrder)
	addRemoteUrlFlag(flags, conf, params)
	flags.BoolVarP(&filterStable, "stable", "s", false, "display only stable version")

	return listRemoteCmd
}

func newResetCmd(conf *config.Config, versionManager versionmanager.VersionManager) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Reset used version of ")
	descBuilder.WriteString(versionManager.FolderName)
	shortMsg := descBuilder.String() + "."

	descBuilder.WriteString(" (remove ")
	descBuilder.WriteString(versionManager.VersionFileName)
	descBuilder.WriteString(" file from TOFUENV_ROOT).")

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return versionManager.Reset()
		},
	}
	return resetCmd
}

func newUninstallCmd(conf *config.Config, versionManager versionmanager.VersionManager) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Uninstall a specific version of ")
	descBuilder.WriteString(versionManager.FolderName)
	shortMsg := descBuilder.String() + "."

	descBuilder.WriteString("(remove it from TOFUENV_ROOT directory).")

	uninstallCmd := &cobra.Command{
		Use:   "uninstall version",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return versionManager.Uninstall(args[0])
		},
	}
	return uninstallCmd
}

func newUseCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	descBuilder.WriteString("Switch the default ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" version to use")
	shortMsg := descBuilder.String() + "."

	descBuilder.WriteString(" (set in ")
	descBuilder.WriteString(versionManager.VersionFileName)
	descBuilder.WriteString(` file in TOFUENV_ROOT)

Available parameter options:
- an exact Semver 2.0.0 version string to use
- a version constraint string (checked against version available in TOFUENV_ROOT directory)
- latest or latest-stable (checked against version available in TOFUENV_ROOT directory)
- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.`)

	workingDir := false

	useCmd := &cobra.Command{
		Use:   "use version",
		Short: shortMsg,
		Long:  descBuilder.String(),
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return versionManager.Use(args[0], workingDir)
		},
	}

	descBuilder.Reset()
	descBuilder.WriteString("create ")
	descBuilder.WriteString(versionManager.VersionFileName)
	descBuilder.WriteString(" file in working directory")

	flags := useCmd.Flags()
	addForceRemoteAndNoInstallFlags(flags, conf, params)
	addKeyFlag(flags, params)
	addRemoteUrlFlag(flags, conf, params)
	flags.BoolVarP(&workingDir, "working-dir", "w", false, descBuilder.String())

	return useCmd
}

func addDescendingFlag(flags *pflag.FlagSet, pReverseOrder *bool) {
	flags.BoolVarP(pReverseOrder, "descending", "d", false, "display list in descending version order")
}

func addForceRemoteAndNoInstallFlags(flags *pflag.FlagSet, conf *config.Config, params subCmdParams) {
	var descBuilder strings.Builder
	descBuilder.WriteString("force search on versions available at ")
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url")

	flags.BoolVarP(&conf.ForceRemote, "force-remote", "f", false, descBuilder.String())
	flags.BoolVarP(&conf.NoInstall, "no-install", "n", conf.NoInstall, "disable installation of missing version")
}

func addKeyFlag(flags *pflag.FlagSet, params subCmdParams) {
	flags.StringVarP(params.pPublicKeyPath, "key-file", "k", "", "local path to PGP public key file (replace check against remote one)")
}

func addRemoteUrlFlag(flags *pflag.FlagSet, conf *config.Config, params subCmdParams) {
	if params.needToken {
		flags.StringVarP(&conf.GithubToken, "github-token", "t", "", "GitHub token (increases GitHub REST API rate limits)")
	}
	flags.StringVarP(params.pRemote, "remote-url", "u", *params.pRemote, "remote url to install from")
}
