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
	"os"
	"strconv"
	"strings"

	"github.com/tofuutils/tenv/v2/config"
	"github.com/tofuutils/tenv/v2/pkg/loghelper"
	"github.com/tofuutils/tenv/v2/versionmanager"
	"github.com/tofuutils/tenv/v2/versionmanager/semantic"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const deprecationMsg = "Direct usage of this subcommand on tenv is deprecated, you should use tofu subcommand instead.\n\n"

func newConstraintCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Set a default constraint expression for ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" (set in TENV_ROOT/")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(`/constraint file).

Without expression reset the default constraint.

The default constraint is added while using latest-allowed, min-required or custom constraint.`)

	constraintCmd := &cobra.Command{
		Use:   "constraint [expression]",
		Short: loghelper.Concat("Set a default constraint expression for ", versionManager.FolderName, "."),
		Long:  descBuilder.String(),
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			conf.InitDisplayer(false)
			addDeprecationMsg(conf, params)

			if len(args) == 0 || args[0] == "" {
				if err := versionManager.ResetConstraint(); err != nil {
					loghelper.StdDisplay(err.Error())
				}

				return
			}

			if err := versionManager.SetConstraint(args[0]); err != nil {
				loghelper.StdDisplay(err.Error())
			}
		},
	}

	return constraintCmd
}

func newDetectCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Display ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" current version.")

	detectCmd := &cobra.Command{
		Use:   "detect",
		Short: loghelper.Concat("Display ", versionManager.FolderName, " current version."),
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			conf.InitDisplayer(false)
			addDeprecationMsg(conf, params)

			detectedVersion, err := versionManager.Detect(false)
			if err != nil {
				loghelper.StdDisplay(err.Error())
			}
			loghelper.StdDisplay(loghelper.Concat(versionManager.FolderName, " ", detectedVersion, " will be run from this directory."))
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
- a version constraint expression (checked against version available at `)
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url)\n- latest, latest-stable or latest-pre (checked against version available at ")
	descBuilder.WriteString(params.remoteEnvName)
	descBuilder.WriteString(" url)\n- latest:<re> or min:<re> to get first version matching with <re> as a regexp after a version sort\n- latest-allowed or min-required to scan your ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" files to detect which version is maximally allowed or minimally required")

	installCmd := &cobra.Command{
		Use:   "install [version]",
		Short: loghelper.Concat("Install a specific version of ", versionManager.FolderName, "."),
		Long:  descBuilder.String(),
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			conf.InitDisplayer(false)
			addDeprecationMsg(conf, params)

			if len(args) == 0 {
				version, err := versionManager.Resolve(semantic.LatestKey)
				if err != nil {
					loghelper.StdDisplay(err.Error())

					return
				}

				if err = versionManager.Install(version); err != nil {
					loghelper.StdDisplay(err.Error())
				}

				return
			}

			if err := versionManager.Install(args[0]); err != nil {
				loghelper.StdDisplay(err.Error())
			}
		},
	}

	flags := installCmd.Flags()
	addInstallationFlags(flags, conf, params)
	addRemoteFlags(flags, conf, params)

	return installCmd
}

func newListCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("List installed ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" versions (located in TENV_ROOT directory), sorted in ascending version order.")

	reverseOrder := false

	listCmd := &cobra.Command{
		Use:   "list",
		Short: loghelper.Concat("List installed ", versionManager.FolderName, " versions."),
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			conf.InitDisplayer(false)
			addDeprecationMsg(conf, params)

			versions, err := versionManager.ListLocal(reverseOrder)
			if err != nil {
				loghelper.StdDisplay(err.Error())

				return
			}

			filePath := versionManager.RootVersionFilePath()
			data, err := os.ReadFile(filePath)
			if err != nil && conf.DisplayVerbose {
				loghelper.StdDisplay("Can not read used version : " + err.Error())
			}
			usedVersion := string(bytes.TrimSpace(data))

			for _, version := range versions {
				if usedVersion == version {
					loghelper.StdDisplay(loghelper.Concat("* ", version, " (set by ", filePath, ")"))
				} else {
					loghelper.StdDisplay("  " + version)
				}
			}
			if conf.DisplayVerbose {
				loghelper.StdDisplay(loghelper.Concat("found ", strconv.Itoa(len(versions)), " ", versionManager.FolderName, " version(s) managed by tenv."))
			}
		},
	}

	addDescendingFlag(listCmd.Flags(), &reverseOrder)

	return listCmd
}

func newListRemoteCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
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
		Short: loghelper.Concat("List installable ", versionManager.FolderName, " versions."),
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			conf.InitDisplayer(false)
			addDeprecationMsg(conf, params)

			versions, err := versionManager.ListRemote(reverseOrder)
			if err != nil {
				loghelper.StdDisplay(err.Error())

				return
			}

			countSkipped := 0
			localSet := versionManager.LocalSet()
			for _, version := range versions {
				if filterStable && !semantic.StableVersion(version) {
					countSkipped++

					continue
				}

				if _, installed := localSet[version]; installed {
					loghelper.StdDisplay(version + " (installed)")
				} else {
					loghelper.StdDisplay(version)
				}
			}
			if conf.DisplayVerbose {
				loghelper.StdDisplay(loghelper.Concat("found ", strconv.Itoa(len(versions)), " ", versionManager.FolderName, " version(s) (on ", params.remoteEnvName, ")."))
				if filterStable {
					loghelper.StdDisplay(strconv.Itoa(countSkipped) + " result(s) hidden (version not stable).")
				}
			}

			return
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
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Reset used version of ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" (remove TENV_ROOT/")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString("/version file).")

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: loghelper.Concat("Reset used version of ", versionManager.FolderName, "."),
		Long:  descBuilder.String(),
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			conf.InitDisplayer(false)
			addDeprecationMsg(conf, params)

			if err := versionManager.ResetVersion(); err != nil {
				loghelper.StdDisplay(err.Error())
			}
		},
	}

	return resetCmd
}

func newUninstallCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Uninstall a specific version of ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" (remove it from TENV_ROOT directory).")

	uninstallCmd := &cobra.Command{
		Use:   "uninstall version",
		Short: loghelper.Concat("Uninstall a specific version of ", versionManager.FolderName, "."),
		Long:  descBuilder.String(),
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			conf.InitDisplayer(false)
			addDeprecationMsg(conf, params)

			if err := versionManager.Uninstall(args[0]); err != nil {
				loghelper.StdDisplay(err.Error())
			}
		},
	}

	return uninstallCmd
}

func newUseCmd(conf *config.Config, versionManager versionmanager.VersionManager, params subCmdParams) *cobra.Command {
	var descBuilder strings.Builder
	addDeprecationHelpMsg(&descBuilder, params)
	descBuilder.WriteString("Switch the default ")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" version to use (set in TENV_ROOT/")
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(`/version file)

Available parameter options:
- an exact Semver 2.0.0 version string to use
- a version constraint expression (checked against version available in TENV_ROOT directory)
- latest, latest-stable or latest-pre (checked against version available in TENV_ROOT directory)
- latest:<re> or min:<re> to get first version matching with <re> as a regexp after a version sort
- latest-allowed or min-required to scan your `)
	descBuilder.WriteString(versionManager.FolderName)
	descBuilder.WriteString(" files to detect which version is maximally allowed or minimally required")

	workingDir := false

	useCmd := &cobra.Command{
		Use:   "use version",
		Short: loghelper.Concat("Switch the default ", versionManager.FolderName, " version to use."),
		Long:  descBuilder.String(),
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			conf.InitDisplayer(false)
			addDeprecationMsg(conf, params)

			if err := versionManager.Use(args[0], workingDir); err != nil {
				loghelper.StdDisplay(err.Error())
			}
		},
	}

	flags := useCmd.Flags()
	addInstallationFlags(flags, conf, params)
	addOptionalInstallationFlags(flags, conf, params)
	addRemoteFlags(flags, conf, params)
	flags.BoolVarP(&workingDir, "working-dir", "w", false, loghelper.Concat("create ", versionManager.VersionFiles[0].Name, " file in working directory"))

	return useCmd
}

func addDeprecationHelpMsg(descBuilder *strings.Builder, params subCmdParams) {
	if params.deprecated {
		descBuilder.WriteString(deprecationMsg)
	}
}

func addDeprecationMsg(conf *config.Config, params subCmdParams) {
	if params.deprecated {
		conf.Displayer.Display(deprecationMsg)
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
	flags.BoolVarP(&conf.ForceRemote, "force-remote", "f", conf.ForceRemote, loghelper.Concat("force search on versions available at ", params.remoteEnvName, " url"))
	flags.BoolVarP(&conf.NoInstall, "no-install", "n", conf.NoInstall, "disable installation of missing version")
}

func addRemoteFlags(flags *pflag.FlagSet, conf *config.Config, params subCmdParams) {
	flags.StringVarP(&conf.RemoteConfPath, "remote-conf", "c", conf.RemoteConfPath, "path to remote configuration file (advanced settings)")
	if params.needToken {
		flags.StringVarP(&conf.GithubToken, "github-token", "t", conf.GithubToken, "GitHub token (increases GitHub REST API rate limits)")
	}
	flags.StringVarP(params.pRemote, "remote-url", "u", "", "remote url to install from")
}
