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
	"fmt"
	"github.com/spf13/cobra"
)

// tfUninstallCmd represents the tfList command
var tfUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall a specific version of Terraform",
	Long:  "Uninstall a specific version of Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tfList called")
	},
}

func init() {
	tfCmd.AddCommand(tfUninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tfUninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	tfUninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
