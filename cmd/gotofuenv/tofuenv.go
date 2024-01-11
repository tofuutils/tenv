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
	"github.com/dvaumoron/gotofuenv/tofuversion"
	"github.com/spf13/cobra"
)

// can be overridden with ldflags
var version = "dev"

func main() {
	conf, err := config.InitConfigFromEnv()
	if err != nil {
		fmt.Println("Configuration error :", err)
		os.Exit(1)
	}

	if err = initCmds(&conf).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCmds(conf *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "gotofuenv",
		Long:    "gotofuenv help manage several version of OpenTofu (https://opentofu.org).",
		Version: version,
	}

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&conf.RemoteUrl, "remote-url", "u", conf.RemoteUrl, "remote url to install from")
	flags.StringVarP(&conf.RootPath, "root-path", "r", conf.RootPath, "local path to install OpenTofu versions")
	flags.StringVarP(&conf.Token, "github-token", "t", "", "GitHub token (increases GitHub REST API rate limits)")
	flags.BoolVarP(&conf.Verbose, "verbose", "v", conf.Verbose, "verbose output")

	rootCmd.AddCommand(newInstallCmd(conf))
	rootCmd.AddCommand(newListCmd(conf))
	rootCmd.AddCommand(newListRemoteCmd(conf))
	rootCmd.AddCommand(newResetCmd(conf))
	rootCmd.AddCommand(newUninstallCmd(conf))
	rootCmd.AddCommand(newUseCmd(conf))
	return rootCmd
}

func newInstallCmd(conf *config.Config) *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install [version]",
		Short: "Install a specific version of OpenTofu.",
		Long: `Install a specific version of OpenTofu (into TOFUENV_ROOT directory from TOFUENV_REMOTE url).

Without parameter the version to use is resolved automatically via TOFUENV_TOFU_VERSION or .opentofu-version files
(searched in working directory, user home directory and TOFUENV_ROOT directory).
Use "latest-stable" when none are found.

If a parameter is passed, available options:
- an exact Semver 2.0.0 version string to install
- a version constraint string (checked against available at TOFUENV_REMOTE url)
- latest (checked against available at TOFUENV_REMOTE url)
- latest-allowed is a syntax to scan your OpenTofu files to detect which version is maximally allowed.
- min-required is a syntax to scan your OpenTofu files to detect which version is minimally required.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			requestedVersion := ""
			if len(args) == 0 {
				requestedVersion = conf.ResolveVersion(config.LatestStableKey)
			} else {
				requestedVersion = args[0]
			}
			return tofuversion.Install(requestedVersion, conf)
		},
	}
	return installCmd
}

func newListCmd(conf *config.Config) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed OpenTofu versions.",
		Long:  "List installed OpenTofu versions (located in TOFUENV_ROOT directory), sorted in ascending version order.",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			versions, err := tofuversion.ListLocal(conf)
			if err != nil {
				return err
			}

			filePath := conf.RootVersionFilePath()
			data, err := os.ReadFile(filePath)
			if err != nil && conf.Verbose {
				fmt.Println("Can not read used version :", err)
			}
			usedVersion := strings.TrimSpace(string(data))

			for _, version := range versions {
				if usedVersion == version {
					fmt.Println("*", version, "(set by", filePath, ")")
				} else {
					fmt.Println(" ", version)
				}
			}
			return nil
		},
	}
	return listCmd
}

func newListRemoteCmd(conf *config.Config) *cobra.Command {
	listRemoteCmd := &cobra.Command{
		Use:   "list-remote",
		Short: "List installable OpenTofu versions.",
		Long:  "List installable OpenTofu versions (from TOFUENV_REMOTE url), sorted in ascending version order.",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			versions, err := tofuversion.ListRemote(conf)
			if err != nil {
				return err
			}

			localSet := tofuversion.LocalSet(conf)
			for _, version := range versions {
				if _, installed := localSet[version]; installed {
					fmt.Println(version, "(installed)")
				} else {
					fmt.Println(version)
				}
			}
			return err
		},
	}
	return listRemoteCmd
}

func newResetCmd(conf *config.Config) *cobra.Command {
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset used version of OpenTofu.",
		Long:  "Reset used version of OpenTofu (remove .opentofu-version file from TOFUENV_ROOT).",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return tofuversion.Reset(conf)
		},
	}
	return resetCmd
}

func newUninstallCmd(conf *config.Config) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall version",
		Short: "Uninstall a specific version of OpenTofu.",
		Long:  "Uninstall a specific version of OpenTofu (remove it from TOFUENV_ROOT directory).",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return tofuversion.Uninstall(args[0], conf)
		},
	}
	return uninstallCmd
}

func newUseCmd(conf *config.Config) *cobra.Command {
	useCmd := &cobra.Command{
		Use:   "use version",
		Short: "Switch the default OpenTofu version to use.",
		Long: `Switch the default OpenTofu version to use (set in .opentofu-version file in TOFUENV_ROOT)

Available parameter options:
- an exact Semver 2.0.0 version string to use
- a version constraint string (checked against available in TOFUENV_ROOT directory)
- latest (checked against available in TOFUENV_ROOT directory)
- latest-allowed is a syntax to scan your OpenTofu files to detect which version is maximally allowed.
- min-required is a syntax to scan your OpenTofu files to detect which version is minimally required.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return tofuversion.Use(args[0], conf)
		},
	}

	flags := useCmd.Flags()
	flags.BoolVarP(&conf.NoInstall, "no-install", "n", conf.NoInstall, "disable installation of missing version")
	flags.BoolVarP(&conf.WorkingDir, "working-dir", "w", false, "create .opentofu-version file in working directory")

	return useCmd
}
