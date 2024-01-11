/*
 *
 * Copyright 2024 opentofuutils authors.
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
	"github.com/opentofuutils/tenv/pkg/tf"
	"github.com/spf13/cobra"
)

// tfInstallCmd represents the tfList command
var tfInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a specific version of Terraform",
	Long:  "Install a specific version of Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		err := tf.InstallSpecificVersion("1.6.6")
		if err != nil {
		}
	},
}

func init() {
	tfCmd.AddCommand(tfInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tfInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	tfInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
