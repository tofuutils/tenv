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
	"github.com/tofuutils/tenv/pkg/github"
	main2 "github.com/tofuutils/tenv/tofuenv"
	"github.com/spf13/cobra"
)

// tofuInstallCmd represents the tofuList command
var tofuInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a specific version of OpenTofu",
	Long:  "Install a specific version of OpenTofu",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := github.GetClient()
		if err != nil {
		}

		err = main2.InstallSpecificVersion(
			client, "opentofu", "opentofu", "1.6.0",
		)
		if err != nil {
		}
	},
}

func init() {
	tofuCmd.AddCommand(tofuInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tofuInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tofuInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
