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
	"strings"

	"github.com/dvaumoron/gotofuenv/config"
	"github.com/dvaumoron/gotofuenv/pkg/iterate"
	"github.com/dvaumoron/gotofuenv/versionmanager"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newInstallCmd(conf *config.Config, versionManager versionmanager.VersionManager, pRemote *string) *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install [version]",
		Short: "Install a specific version of OpenTofu.",
		Long: `Install a specific version of OpenTofu (into TOFUENV_ROOT directory from TOFUENV_REMOTE url).

Without parameter the version to use is resolved automatically via TOFUENV_TOFU_VERSION or .opentofu-version files
(searched in working directory, user home directory and TOFUENV_ROOT directory).
Use "latest-stable" when none are found.

If a parameter is passed, available options:
- an exact Semver 2.0.0 version string to install
- a version constraint string (checked against version available at TOFUENV_REMOTE url)
- latest or latest-stable (checked against version available at TOFUENV_REMOTE url)
- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			requestedVersion := ""
			if len(args) == 0 {
				requestedVersion = versionManager.Resolve(config.LatestStableKey)
			} else {
				requestedVersion = args[0]
			}
			return versionManager.Install(requestedVersion)
		},
	}

	addRemoteUrlFlag(installCmd.Flags(), pRemote)

	return installCmd
}

func newListCmd(conf *config.Config, versionManager versionmanager.VersionManager) *cobra.Command {
	reverseOrder := false

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed OpenTofu versions.",
		Long:  "List installed OpenTofu versions (located in TOFUENV_ROOT directory), sorted in ascending version order.",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			versions, err := versionManager.ListLocal()
			if err != nil {
				return err
			}

			filePath := versionManager.RootVersionFilePath()
			data, err := os.ReadFile(filePath)
			if err != nil && conf.Verbose {
				fmt.Println("Can not read used version :", err)
			}
			usedVersion := strings.TrimSpace(string(data))

			versionReceiver, done := iterate.Iterate(versions, reverseOrder)
			defer done()

			for version := range versionReceiver {
				if usedVersion == version {
					fmt.Println("*", version, "(set by", filePath+")")
				} else {
					fmt.Println(" ", version)
				}
			}
			return nil
		},
	}

	addDescendingFlag(listCmd.Flags(), &reverseOrder)

	return listCmd
}

func newListRemoteCmd(conf *config.Config, versionManager versionmanager.VersionManager, pRemote *string) *cobra.Command {
	filterStable := false
	reverseOrder := false

	listRemoteCmd := &cobra.Command{
		Use:   "list-remote",
		Short: "List installable OpenTofu versions.",
		Long:  "List installable OpenTofu versions (from TOFUENV_REMOTE url), sorted in ascending version order.",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			versions, err := versionManager.ListRemote()
			if err != nil {
				return err
			}

			versionReceiver, done := iterate.Iterate(versions, reverseOrder)
			defer done()

			localSet := versionManager.LocalSet()
			for version := range versionReceiver {
				if filterStable && !versionmanager.StableVersion(version) {
					continue
				}

				if _, installed := localSet[version]; installed {
					fmt.Println(version, "(installed)")
				} else {
					fmt.Println(version)
				}
			}
			return err
		},
	}

	flags := listRemoteCmd.Flags()
	addDescendingFlag(flags, &reverseOrder)
	addRemoteUrlFlag(flags, pRemote)
	flags.BoolVarP(&filterStable, "stable", "s", false, "display only stable version")

	return listRemoteCmd
}

func newResetCmd(conf *config.Config, versionManager versionmanager.VersionManager) *cobra.Command {
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset used version of OpenTofu.",
		Long:  "Reset used version of OpenTofu (remove .opentofu-version file from TOFUENV_ROOT).",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return versionManager.Reset()
		},
	}
	return resetCmd
}

func newUninstallCmd(conf *config.Config, versionManager versionmanager.VersionManager) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall version",
		Short: "Uninstall a specific version of OpenTofu.",
		Long:  "Uninstall a specific version of OpenTofu (remove it from TOFUENV_ROOT directory).",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return versionManager.Uninstall(args[0])
		},
	}
	return uninstallCmd
}

func newUseCmd(conf *config.Config, versionManager versionmanager.VersionManager, pRemote *string) *cobra.Command {
	forceRemote := false
	workingDir := false

	useCmd := &cobra.Command{
		Use:   "use version",
		Short: "Switch the default OpenTofu version to use.",
		Long: `Switch the default OpenTofu version to use (set in .opentofu-version file in TOFUENV_ROOT)

Available parameter options:
- an exact Semver 2.0.0 version string to use
- a version constraint string (checked against version available in TOFUENV_ROOT directory)
- latest or latest-stable (checked against version available in TOFUENV_ROOT directory)
- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return versionManager.Use(args[0], forceRemote, workingDir)
		},
	}

	flags := useCmd.Flags()
	flags.BoolVarP(&forceRemote, "force-remote", "f", false, "force search version available at TOFUENV_REMOTE url")
	flags.BoolVarP(&conf.NoInstall, "no-install", "n", conf.NoInstall, "disable installation of missing version")
	addRemoteUrlFlag(flags, pRemote)
	flags.BoolVarP(&workingDir, "working-dir", "w", false, "create .opentofu-version file in working directory")

	return useCmd
}

func addDescendingFlag(flags *pflag.FlagSet, pReverseOrder *bool) {
	flags.BoolVarP(pReverseOrder, "descending", "d", false, "display list in descending version order")
}

func addRemoteUrlFlag(flags *pflag.FlagSet, pRemote *string) {
	flags.StringVarP(pRemote, "remote-url", "u", *pRemote, "remote url to install from")
}
